package render

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"

	ovnoperatorv1 "github.com/harvester/kubeovn-operator/api/v1"
	sourcetemplate "github.com/harvester/kubeovn-operator/internal/templates"
)

func generateConfigObject() (*ovnoperatorv1.Configuration, error) {
	content, err := os.ReadFile("../../config/samples/kubeovn.io_v1_configuration.yaml")
	if err != nil {
		return nil, fmt.Errorf("error reading sample configuration file: %v", err)
	}
	config := &ovnoperatorv1.Configuration{}
	err = yaml.Unmarshal(content, config)
	return config, err
}

func Test_ObjectRendering(t *testing.T) {
	assert := require.New(t)
	c, err := generateConfigObject()
	assert.NoError(err, "expected no error while generating config object")
	for objectType, objectList := range sourcetemplate.OrderedObjectList {
		returnedObjects, err := GenerateObjects(objectList, c, objectType, nil)
		assert.NoError(err, "expected no error while generating object", objectType)
		for _, object := range returnedObjects {
			assert.NotEmpty(object.GetName())
		}
	}
}
