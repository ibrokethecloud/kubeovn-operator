package controller

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/go-logr/logr"
	kubeovniov1 "github.com/harvester/kubeovn-operator/api/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// +kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch;create;update;patch;delete
type HealthCheckReconciler struct {
	client.Client
	RestConfig          *rest.Config
	Scheme              *runtime.Scheme
	EventRecorder       record.EventRecorder
	Namespace           string
	Log                 logr.Logger
	HealthCheckInterval int
}

func (r *HealthCheckReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubeovniov1.Configuration{}).
		Named("kubeovn-healthcheck-controller").Complete(r)
}

// HealthCheckReconciler runs a separate reconcile loop where configuration object is requeued every 300 seconds
// to force a re-run of the healthcheck
func (r *HealthCheckReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
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

	// deletion timestamp is set on object, no need for run further checks
	if !config.DeletionTimestamp.IsZero() {
		r.Log.WithValues("name", configObj.Name).Info("configuration being deleted, no further healthchecks needed")
		return ctrl.Result{}, nil
	}

	if config.Status.Status != kubeovniov1.ConfigurationStatusDeployed {
		r.Log.WithValues("name", configObj.Name).Info("waiting for resources to be deployed")
		return ctrl.Result{}, nil
	}

	if err := r.reconcileOVNDBHealth(ctx, config); err != nil {
		return ctrl.Result{}, fmt.Errorf("error during execution of reconcileOVNDBHealth: %v", err)
	}

	// healthcheck only updates conditions. since object is also reconciled by another controller we ignore the rest
	if !reflect.DeepEqual(config.Status.Conditions, configObj.Status.Conditions) {
		if err := r.Client.Status().Patch(ctx, config, client.MergeFrom(configObj)); err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{RequeueAfter: time.Duration(r.HealthCheckInterval) * time.Second}, nil
}

// reconcileOVNDBHealth reconciles health of north and south db
func (r *HealthCheckReconciler) reconcileOVNDBHealth(ctx context.Context, config *kubeovniov1.Configuration) error {
	if !r.checkNeeded(config) {
		return nil
	}
	var runNBCheck, runSBCheck bool

	// check if NorthBound leader exists
	nbPods, err := podList(ctx, kubeovniov1.NBLeaderLabel, r.Client)
	if err != nil {
		return fmt.Errorf("error fetching northbound leader: %v", err)
	}

	if len(nbPods.Items) == 0 {
		config.SetCondition(kubeovniov1.OVNNBLeaderFound, metav1.ConditionFalse, "", "no pods matching northbound leader label requirements found")
	} else {
		runNBCheck = true
		config.SetCondition(kubeovniov1.OVNNBLeaderFound, metav1.ConditionTrue, "", fmt.Sprintf("northbound leader found %s", nbPods.Items[0].GetName()))
	}

	// check if SouthBound leader exists
	sbPods, err := podList(ctx, kubeovniov1.SBLeaderLabel, r.Client)
	if err != nil {
		return fmt.Errorf("error fetching southbound leader: %v", err)
	}

	if len(sbPods.Items) == 0 {
		config.SetCondition(kubeovniov1.OVNSBLeaderFound, metav1.ConditionFalse, "", "no pods matching southbound leader label requirements found")
	} else {
		runSBCheck = true
		config.SetCondition(kubeovniov1.OVNSBLeaderFound, metav1.ConditionTrue, "", fmt.Sprintf("northbound leader found %s", sbPods.Items[0].GetName()))
	}

	// run health check on northbound db
	if runNBCheck {
		result, err := executeOVNCentralCommand(ctx, kubeovniov1.NBCheckScript, kubeovniov1.SBLeaderLabel, r.Client, r.RestConfig)
		if err != nil {
			r.Log.Error(err, "NBCheck failure")
			config.SetCondition(kubeovniov1.OVNNBDBHealth, metav1.ConditionFalse, "", string(result))
		} else {
			config.SetCondition(kubeovniov1.OVNNBDBHealth, metav1.ConditionTrue, "", string(result))
		}
	}

	if runSBCheck {
		result, err := executeOVNCentralCommand(ctx, kubeovniov1.SBCheckScript, kubeovniov1.SBLeaderLabel, r.Client, r.RestConfig)
		if err != nil {
			r.Log.Error(err, "NBCheck failure")
			config.SetCondition(kubeovniov1.OVNNBDBHealth, metav1.ConditionFalse, "", string(result))
		} else {
			config.SetCondition(kubeovniov1.OVNNBDBHealth, metav1.ConditionTrue, "", string(result))
		}
	}
	return nil
}

// checkNeeded calculates if Healthcheck interval has passed before triggering another health check
func (r *HealthCheckReconciler) checkNeeded(config *kubeovniov1.Configuration) bool {
	condition := config.LookupCondition(kubeovniov1.OVNNBLeaderFound)
	return condition.LastTransitionTime.Add(time.Duration(r.HealthCheckInterval) * time.Second).Before(metav1.Now().Time)
}
