/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"

	kubeovniov1 "github.com/harvester/kubeovn-operator/api/v1"
	"github.com/harvester/kubeovn-operator/internal/render"
	"github.com/harvester/kubeovn-operator/internal/templates"
)

// ConfigurationReconciler reconciles a Configuration object
type ConfigurationReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	Namespace     string
	EventRecorder record.EventRecorder
	Log           logr.Logger
}

type reconcileFuncs func(context.Context, *kubeovniov1.Configuration) error

// orderedObjectList iterates templated object lists and applies them in order
var orderedObjectList = map[client.Object][]string{
	&apiextensionsv1.CustomResourceDefinition{}: templates.CRDList,
	&corev1.Secret{}:             templates.SecretList,
	&corev1.ServiceAccount{}:     templates.ServiceAccountList,
	&rbacv1.RoleBinding{}:        templates.RoleBindingList,
	&rbacv1.ClusterRole{}:        templates.ClusterRoleList,
	&rbacv1.ClusterRoleBinding{}: templates.ClusterRoleBindingList,
	&corev1.ConfigMap{}:          templates.ConfigMapList,
	&appsv1.Deployment{}:         templates.DeploymentList,
	&appsv1.DaemonSet{}:          templates.DaemonsetList,
	&corev1.Service{}:            templates.ServicesList,
}

// +kubebuilder:rbac:groups=kubeovn.io,resources=configurations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kubeovn.io,resources=configurations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=kubeovn.io,resources=configurations/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Configuration object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.4/pkg/reconcile
func (r *ConfigurationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	configObj := &kubeovniov1.Configuration{}
	err := r.Get(ctx, req.NamespacedName, configObj)
	if err != nil {
		if apierrors.IsNotFound(err) {
			r.Log.WithValues("name", configObj.Name).Info("configuration not found")
			return ctrl.Result{}, nil
		}
		r.Log.WithValues("name", configObj.Name).Error(err, "error fetching object")
		return ctrl.Result{}, err
	}

	config := configObj.DeepCopy()
	// if deletiontimestamp is set, then no further processing is needed as we let k8s gc the associated objects
	if config.DeletionTimestamp != nil {
		return reconcile.Result{}, nil
	}

	reconcileSteps := []reconcileFuncs{r.initializeConditions, r.findMasterNodes, r.checkObjects, r.applyObject}

	for _, v := range reconcileSteps {
		if err := v(ctx, config); err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigurationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	b := ctrl.NewControllerManagedBy(mgr).
		For(&kubeovniov1.Configuration{}).
		Named("configuration")
	return r.AddWatches(b).Complete(r)
}

// applyObject will check if Config object is not already deploying. If a change is needed then it triggers
// create/update of associated objects
func (r *ConfigurationReconciler) applyObject(ctx context.Context, config *kubeovniov1.Configuration) error {
	if len(config.Status.MatchingNodeAddresses) == 0 {
		r.Log.WithValues("name", config.Name).Info("waiting for matching master node requirement to be met")
		return nil
	}
	if config.Status.Status == kubeovniov1.ConfigurationStatusDeploying {
		r.Log.WithValues("name", config.Name).Info("skipping applying objects as objects are already deploying")
		return nil
	}

	for objectType, objectList := range orderedObjectList {
		objs, err := render.GenerateObjects(objectList, config, objectType)
		if err != nil {
			return fmt.Errorf("error during object generation for type %s: %v", objectType.GetObjectKind().GroupVersionKind(), err)
		}
		for _, obj := range objs {
			if err := controllerutil.SetControllerReference(config, obj, r.Scheme); err != nil {
				return fmt.Errorf("error setting controller reference on object %s/%s: %v", obj.GetNamespace(), obj.GetName(), err)
			}
			if err = r.reconcileObject(ctx, obj); err != nil {
				return fmt.Errorf("error reconcilling object %s/%s: %v", obj.GetNamespace(), obj.GetName(), err)
			}
		}
	}
	config.Status.Status = kubeovniov1.ConfigurationStatusDeployed
	return r.Status().Update(ctx, config)
}

// reconcileObject will create / update the managed objects
func (r *ConfigurationReconciler) reconcileObject(ctx context.Context, obj client.Object) error {
	existingObject := obj
	err := r.Client.Get(ctx, types.NamespacedName{Name: obj.GetName(), Namespace: obj.GetNamespace()}, existingObject)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return r.Create(ctx, obj)
		}
		return fmt.Errorf("error fetching object %s/%s: %v", obj.GetNamespace(), obj.GetName(), err)
	}

	return r.Patch(ctx, obj, client.MergeFrom(existingObject))
}

func (r *ConfigurationReconciler) filterObject(ctx context.Context, obj client.Object) []ctrl.Request {
	ownerRefs := obj.GetOwnerReferences()
	result := []ctrl.Request{}
	if len(ownerRefs) == 0 {
		return result
	}

	for _, v := range ownerRefs {
		if v.Kind == kubeovniov1.Kind && v.APIVersion == kubeovniov1.APIVersion {
			result = append(result, ctrl.Request{NamespacedName: types.NamespacedName{Name: v.Name, Namespace: r.Namespace}})
		}
	}
	return result
}

func (r *ConfigurationReconciler) AddWatches(b *builder.Builder) *builder.Builder {
	for key, _ := range orderedObjectList {
		b.Watches(key, handler.EnqueueRequestsFromMapFunc(r.filterObject))
	}
	return b
}

// checkObjects checks status of deployed objects
func (r *ConfigurationReconciler) checkObjects(ctx context.Context, config *kubeovniov1.Configuration) error {
	return nil
}

// findMasterNodes will find nodes matching the master label criteria in the configuration
func (r *ConfigurationReconciler) findMasterNodes(ctx context.Context, config *kubeovniov1.Configuration) error {
	nodeList := &corev1.NodeList{}
	err := r.List(ctx, nodeList)
	if err != nil {
		return fmt.Errorf("error listing nodes :%v", err)
	}

	var nodeAddresses []string
	for _, v := range nodeList.Items {
		if metav1.HasLabel(v.ObjectMeta, config.Spec.MasterNodesLabel) {
			for _, address := range v.Status.Addresses {
				if address.Type == corev1.NodeInternalIP {
					nodeAddresses = append(nodeAddresses, address.Address)
				}
			}
		}
	}

	// if no nodeAddresses are found then it is likely we had no matching nodes
	// we need to pause reconcile of the object until label matches
	if len(nodeAddresses) == 0 {
		r.EventRecorder.Event(config, corev1.EventTypeWarning,
			"ReconcilePaused", "no nodes matching master node labels found")
		return nil
	}

	config.Status.MatchingNodeAddresses = nodeAddresses
	return r.Status().Update(ctx, config)
}

// initializeConditions will initialise baseline conditions for the configuration object
func (r *ConfigurationReconciler) initializeConditions(ctx context.Context, config *kubeovniov1.Configuration) error {
	if len(config.Status.Conditions) != 2 {
		return nil
	}
	configObj := config.DeepCopy()
	config.SetCondition(kubeovniov1.ErroredObjectsCondition, metav1.ConditionUnknown, "", "")
	config.SetCondition(kubeovniov1.WaitingForMatchignNodesCondition, metav1.ConditionTrue, "", "")
	return r.Status().Patch(ctx, config, client.MergeFrom(configObj))
}
