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

package v1

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	kubeovnv1 "github.com/harvester/kubeovn-operator/api/v1"
)

// nolint:unused
// log is for logging in this package.
var (
	configurationlog              = logf.Log.WithName("configuration-resource")
	ovnCentralDefaultResourceSpec = kubeovnv1.ResourceSpec{
		Requests: kubeovnv1.CPUMemSpec{
			CPU:    *resource.NewMilliQuantity(300, resource.BinarySI),
			Memory: *resource.NewMilliQuantity(200, resource.BinarySI),
		},
		Limits: kubeovnv1.CPUMemSpec{
			CPU:    *resource.NewQuantity(3, resource.BinarySI),
			Memory: *resource.NewQuantity(4, resource.BinarySI),
		},
	}
	ovsOVNDefaultResourceSpec = kubeovnv1.ResourceSpec{
		Requests: kubeovnv1.CPUMemSpec{
			CPU:    *resource.NewMilliQuantity(200, resource.BinarySI),
			Memory: *resource.NewMilliQuantity(200, resource.BinarySI),
		},
		Limits: kubeovnv1.CPUMemSpec{
			CPU:    *resource.NewQuantity(2, resource.BinarySI),
			Memory: *resource.NewQuantity(1, resource.BinarySI),
		},
	}
	kubeOVNControllerDefaultResourceSpec = kubeovnv1.ResourceSpec{
		Requests: kubeovnv1.CPUMemSpec{
			CPU:    *resource.NewMilliQuantity(200, resource.BinarySI),
			Memory: *resource.NewMilliQuantity(200, resource.BinarySI),
		},
		Limits: kubeovnv1.CPUMemSpec{
			CPU:    *resource.NewQuantity(1, resource.BinarySI),
			Memory: *resource.NewQuantity(1, resource.BinarySI),
		},
	}
	kubeOVNCNIDefaultResourceSpec = kubeovnv1.ResourceSpec{
		Requests: kubeovnv1.CPUMemSpec{
			CPU:    *resource.NewMilliQuantity(100, resource.BinarySI),
			Memory: *resource.NewMilliQuantity(100, resource.BinarySI),
		},
		Limits: kubeovnv1.CPUMemSpec{
			CPU:    *resource.NewQuantity(1, resource.BinarySI),
			Memory: *resource.NewQuantity(1, resource.BinarySI),
		},
	}
	kubeOVNPingerDefaultResourceSpec = kubeovnv1.ResourceSpec{
		Requests: kubeovnv1.CPUMemSpec{
			CPU:    *resource.NewMilliQuantity(100, resource.BinarySI),
			Memory: *resource.NewMilliQuantity(100, resource.BinarySI),
		},
		Limits: kubeovnv1.CPUMemSpec{
			CPU:    *resource.NewMilliQuantity(200, resource.BinarySI),
			Memory: *resource.NewMilliQuantity(400, resource.BinarySI),
		},
	}
	kubeOVNMonitorDefaultResourceSpec = kubeovnv1.ResourceSpec{
		Requests: kubeovnv1.CPUMemSpec{
			CPU:    *resource.NewMilliQuantity(200, resource.BinarySI),
			Memory: *resource.NewMilliQuantity(200, resource.BinarySI),
		},
		Limits: kubeovnv1.CPUMemSpec{
			CPU:    *resource.NewMilliQuantity(200, resource.BinarySI),
			Memory: *resource.NewMilliQuantity(200, resource.BinarySI),
		},
	}
)

// SetupConfigurationWebhookWithManager registers the webhook for Configuration in the manager.
func SetupConfigurationWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&kubeovnv1.Configuration{}).
		WithValidator(&ConfigurationCustomValidator{}).
		WithDefaulter(&ConfigurationCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-kubeovn-io-v1-configuration,mutating=true,failurePolicy=fail,sideEffects=None,groups=kubeovn.io,resources=configurations,verbs=create;update,versions=v1,name=mconfiguration-v1.kb.io,admissionReviewVersions=v1

// ConfigurationCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind Configuration when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type ConfigurationCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &ConfigurationCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind Configuration.
func (d *ConfigurationCustomDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	configuration, ok := obj.(*kubeovnv1.Configuration)

	if !ok {
		return fmt.Errorf("expected an Configuration object but got %T", obj)
	}
	configurationlog.Info("Defaulting for Configuration", "name", configuration.GetName())

	d.ApplyConfigurationDefaults(configuration)
	return nil
}

func (d *ConfigurationCustomDefaulter) ApplyConfigurationDefaults(config *kubeovnv1.Configuration) {
	config.Spec.OVNCentral = applyDefaults(config.Spec.OVNCentral, ovnCentralDefaultResourceSpec)
	config.Spec.OVSOVN = applyDefaults(config.Spec.OVSOVN, ovsOVNDefaultResourceSpec)
	config.Spec.KubeOVNController = applyDefaults(config.Spec.KubeOVNController, kubeOVNControllerDefaultResourceSpec)
	config.Spec.KubeOVNCNI = applyDefaults(config.Spec.KubeOVNCNI, kubeOVNCNIDefaultResourceSpec)
	config.Spec.KubeOVNPinger = applyDefaults(config.Spec.KubeOVNPinger, kubeOVNPingerDefaultResourceSpec)
	config.Spec.KubeOVNMonitor = applyDefaults(config.Spec.KubeOVNMonitor, kubeOVNMonitorDefaultResourceSpec)
}

// applyDefaults will apply baseline defaults for resource specs to configuration
func applyDefaults(resource, defaultValues kubeovnv1.ResourceSpec) kubeovnv1.ResourceSpec {
	if resource.Requests.CPU.IsZero() {
		resource.Requests.CPU = defaultValues.Requests.CPU
	}

	if resource.Requests.Memory.IsZero() {
		resource.Requests.Memory = defaultValues.Requests.Memory
	}

	if resource.Limits.CPU.IsZero() {
		resource.Limits.CPU = defaultValues.Limits.CPU
	}

	if resource.Limits.Memory.IsZero() {
		resource.Limits.Memory = defaultValues.Limits.Memory
	}
	return resource
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-kubeovn-io-v1-configuration,mutating=false,failurePolicy=fail,sideEffects=None,groups=kubeovn.io,resources=configurations,verbs=create;update,versions=v1,name=vconfiguration-v1.kb.io,admissionReviewVersions=v1

// ConfigurationCustomValidator struct is responsible for validating the Configuration resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type ConfigurationCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &ConfigurationCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type Configuration.
func (v *ConfigurationCustomValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	configuration, ok := obj.(*kubeovnv1.Configuration)
	if !ok {
		return nil, fmt.Errorf("expected a Configuration object but got %T", obj)
	}
	configurationlog.Info("Validation for Configuration upon creation", "name", configuration.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type Configuration.
func (v *ConfigurationCustomValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	configuration, ok := newObj.(*kubeovnv1.Configuration)
	if !ok {
		return nil, fmt.Errorf("expected a Configuration object for the newObj but got %T", newObj)
	}
	configurationlog.Info("Validation for Configuration upon update", "name", configuration.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type Configuration.
func (v *ConfigurationCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	configuration, ok := obj.(*kubeovnv1.Configuration)
	if !ok {
		return nil, fmt.Errorf("expected a Configuration object but got %T", obj)
	}
	configurationlog.Info("Validation for Configuration upon deletion", "name", configuration.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
