package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/go-logr/logr"
	kubeovniov1 "github.com/harvester/kubeovn-operator/api/v1"
)

const (
	KubeOVNNodeFinalizer = "nodes.kubeovn.io"
)

// +kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch;create;update;patch;delete
type NodeReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	EventRecorder record.EventRecorder
	Namespace     string
	Log           logr.Logger
}

func (r *NodeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return nil
}

// Reconcile finds nodes matching the configuration .spec.masterNodesLabel
// and adds a finalizer to all nodes. This finalizer is then used during node deletion
// events to trigger cleanup of ovn north and south bound databases
func (r *NodeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	nodeObj := &corev1.Node{}
	err := r.Get(ctx, req.NamespacedName, nodeObj)
	if err != nil {
		if apierrors.IsNotFound(err) {
			r.Log.WithValues("name", req.Name).Info("node not found")
			return ctrl.Result{}, nil
		}
	}

	node := nodeObj.DeepCopy()
	// fetch configuration object and identify matching nodes
	config, err := fetchKubeovnConfig(ctx, r.Client, r.Namespace)
	if err != nil {
		if apierrors.IsNotFound(err) {
			r.Log.Info("waiting for kubeovn configuration to be created, nothing to do: %v", err)
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if !metav1.HasLabel(node.ObjectMeta, config.Spec.MasterNodesLabel) {
		r.Log.WithValues("name", req.Name).Info("node does not match config masterNodesLabel, so ignoring")
		return ctrl.Result{}, nil
	}

	if node.DeletionTimestamp != nil {
		// reconcile ovn north and south databases and check if there is a condition matching
		return r.reconcileNodeDeletion(ctx, config, node)
	}

	// add finalizer if one does not exist and update node object
	if !controllerutil.ContainsFinalizer(node, kubeovniov1.KubeOVNNodeFinalizer) {
		controllerutil.AddFinalizer(node, kubeovniov1.KubeOVNNodeFinalizer)
		return ctrl.Result{}, r.Patch(ctx, node, client.MergeFrom(nodeObj))
	}

	return ctrl.Result{}, nil
}

// fetchKubeovnConfig fetches the default configuration
func fetchKubeovnConfig(ctx context.Context, client client.Client, namespace string) (*kubeovniov1.Configuration, error) {
	config := &kubeovniov1.Configuration{}
	err := client.Get(ctx, types.NamespacedName{Name: kubeovniov1.DefaultConfigurationName, Namespace: namespace}, config)
	return config, err
}

// reconcileNodeDeletion will delete the node entry from the ovn database and remove finalizer
// allowing node to be removed from apiserver
func (r *NodeReconciler) reconcileNodeDeletion(ctx context.Context, config *kubeovniov1.Configuration, node *corev1.Node) (ctrl.Result, error) {
	nodeObj := node.DeepCopy()
	// perform ovn northbound and southbound db cleanup operations
	// remove finalizer to allow node to be cleaned up
	if controllerutil.ContainsFinalizer(node, kubeovniov1.KubeOVNNodeFinalizer) {
		controllerutil.RemoveFinalizer(node, kubeovniov1.KubeOVNNodeFinalizer)
		return ctrl.Result{}, r.Patch(ctx, node, client.MergeFrom(nodeObj))
	}
	return ctrl.Result{}, nil
}
