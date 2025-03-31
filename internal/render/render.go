package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/Masterminds/sprig/v3"

	ovnoperatorv1 "github.com/harvester/kubeovn-operator/api/v1"
)

type values struct {
	Values ovnoperatorv1.ConfigurationSpec `json:"Values"`
}

func GenerateObjects(templates []string, config *ovnoperatorv1.Configuration, object client.Object) ([]client.Object, error) {
	var returnedObjects []client.Object

	valsObj, err := generateMap(config)
	if err != nil {
		return nil, fmt.Errorf("failed to generate values: %w", err)
	}
	for _, sourceTemplate := range templates {
		returnedObject, err := generateObject(sourceTemplate, valsObj, object)
		if err != nil {
			return nil, fmt.Errorf("error returned from generateTemplate: %v", err)
		}
		returnedObjects = append(returnedObjects, returnedObject)
	}
	return returnedObjects, nil
}

func generateObject(input string, valuesObj map[string]interface{}, object client.Object) (client.Object, error) {
	var output bytes.Buffer
	tmpl := template.Must(template.New("objects").Funcs(sprig.FuncMap()).Parse(input))
	err := tmpl.Execute(&output, valuesObj)
	if err != nil {
		return nil, fmt.Errorf("error rending template: %v", err)
	}

	err = yaml.Unmarshal(output.Bytes(), object)
	return object, err
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
