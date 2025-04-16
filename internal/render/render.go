package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	ovnoperatorv1 "github.com/harvester/kubeovn-operator/api/v1"
	"github.com/iancoleman/strcase"
	"helm.sh/helm/v4/pkg/engine"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

type values struct {
	Values ovnoperatorv1.ConfigurationSpec `json:"Values"`
}

func GenerateObjects(templates []string, config *ovnoperatorv1.Configuration, object client.Object, restConfig *rest.Config) ([]client.Object, error) {
	var returnedObjects []client.Object

	valsObj, err := generateMap(config)
	if err != nil {
		return nil, fmt.Errorf("failed to generate values: %w", err)
	}

	err = unstructured.SetNestedField(valsObj, config.GetNamespace(), "Values", "namespace")
	if err != nil {
		return nil, fmt.Errorf("failed to set namespace field in values map: %v", err)
	}

	// additional fields needed which are generated by the helper template in the helm chart
	kubeovn := generateIncludeValues(config)
	err = unstructured.SetNestedField(valsObj, kubeovn, "Values", "kubeovn")
	if err != nil {
		return nil, fmt.Errorf("failed to set kubeovn field in values map: %v", err)
	}

	//valsObj = AlternateMapValues(valsObj)
	for _, sourceTemplate := range templates {
		returnedObject, err := generateObject(sourceTemplate, valsObj, object, restConfig)
		if err != nil {
			return nil, fmt.Errorf("error returned from generateTemplate: %v", err)
		}
		if returnedObject != nil {
			returnedObjects = append(returnedObjects, returnedObject)
		}
	}
	return returnedObjects, nil
}

func generateObject(input string, valuesObj map[string]interface{}, object client.Object, restConfig *rest.Config) (client.Object, error) {
	newObj := initialiseNewObject(object)
	if newObj == nil {
		return nil, fmt.Errorf("could not initialise new object for type: %T", object)
	}
	var output bytes.Buffer
	f := sprig.TxtFuncMap()
	f["lookup"] = engine.NewLookupFunction(restConfig)
	f["include"] = include
	tmpl := template.Must(template.New("objects").Funcs(f).Parse(input))
	err := tmpl.Execute(&output, valuesObj)
	if err != nil {
		return nil, fmt.Errorf("error rending template: %v", err)
	}
	// object is skipped due to condition in the template
	if len(output.String()) == 0 {
		return nil, nil
	}
	err = yaml.Unmarshal(output.Bytes(), newObj)
	if err != nil {
		return nil, err
	}
	return newObj, nil
}

// needed to assert convert ensure templates need little to no change at all
func generateMap(config *ovnoperatorv1.Configuration) (map[string]interface{}, error) {
	val := values{Values: config.Spec}
	out, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}

	resultMap := make(map[string]interface{})
	err = json.Unmarshal(out, &resultMap)
	return resultMap, err
}

// generateIncludeValues renders the values generated the _helpers.tpl in the chart
func generateIncludeValues(config *ovnoperatorv1.Configuration) map[string]interface{} {
	nodeCount := strings.Join(config.Status.MatchingNodeAddresses, ",")
	upgradeStratergyMap := map[string]interface{}{
		"upgradeStratergy": "RollingUpdate",
	}

	// https://github.com/kubernetes/apimachinery/blob/v0.32.1/pkg/runtime/converter.go#L614
	// converter only supports int64 and float64, which is why the forced type conversion
	// for versionCompatability and runAsUser
	versionCompatibiltyMap := map[string]interface{}{
		"versionCompatability": float64(24.03),
	}

	runAsUser := int64(65534)
	if *config.Spec.Component.EnableOVNIPSec {
		runAsUser = int64(0)
	}

	return map[string]interface{}{
		"nodeCount": nodeCount,
		"ovs-ovn":   upgradeStratergyMap,
		"ovn":       versionCompatibiltyMap,
		"runAsUser": runAsUser,
	}

}

// mimic helm include but in this case values are already pre-calculate in the
// map when we render the values
func include(key string, data map[string]interface{}) string {
	// split key into fields

	fields := strings.Split(key, ".")
	val, _, _ := unstructured.NestedString(data, fields...)
	return val
}

// AlternateMapValues attempts to add screaming snake case keys for all existing keys to ensure that some of the helm
// templates which use the screaming snake case can be rendered correctly
func AlternateMapValues(data map[string]interface{}) map[string]interface{} {
	for key, value := range data {
		switch v := value.(type) {
		case map[string]interface{}:
			data[strcase.ToScreamingSnake(key)] = AlternateMapValues(v)
			data[strcase.ToKebab(key)] = AlternateMapValues(v)
		default:
			data[strcase.ToScreamingSnake(key)] = value
			data[strcase.ToKebab(key)] = value
		}
	}
	return data
}

func initialiseNewObject(object client.Object) client.Object {
	switch object.(type) {
	case *appsv1.Deployment:
		return &appsv1.Deployment{}
	case *apiextensionsv1.CustomResourceDefinition:
		return &apiextensionsv1.CustomResourceDefinition{}
	case *corev1.Secret:
		return &corev1.Secret{}
	case *corev1.ServiceAccount:
		return &corev1.ServiceAccount{}
	case *corev1.ConfigMap:
		return &corev1.ConfigMap{}
	case *rbacv1.RoleBinding:
		return &rbacv1.RoleBinding{}
	case *rbacv1.ClusterRole:
		return &rbacv1.ClusterRole{}
	case *rbacv1.ClusterRoleBinding:
		return &rbacv1.ClusterRoleBinding{}
	case *appsv1.DaemonSet:
		return &appsv1.DaemonSet{}
	case *corev1.Service:
		return &corev1.Service{}
	}
	return nil
}
