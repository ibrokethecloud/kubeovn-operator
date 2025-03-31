package render

import (
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	sourcetemplate "github.com/harvester/kubeovn-operator/internal/templates"

	ovnoperatorv1 "github.com/harvester/kubeovn-operator/api/v1"
)

var config = &ovnoperatorv1.Configuration{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "sample",
		Namespace: "kube-system",
	},
	Spec: ovnoperatorv1.ConfigurationSpec{
		Global: ovnoperatorv1.GlobalSpec{
			Registry: ovnoperatorv1.RegistrySpec{
				ImagePullSecrets: []string{"registry-one-secret", "registry-two-secret"},
			},
		},
		Namespace: "kube-system",
	},
}

func Test_GenerateSAObjects(t *testing.T) {
	sa := corev1.ServiceAccount{}
	assert := require.New(t)
	returnedObjects, err := GenerateObjects(sourcetemplate.ServiceAccountList, config, &sa)
	assert.NoError(err)
	assert.Len(returnedObjects, 4)

}

func Test_generateMap(t *testing.T) {
	assert := require.New(t)
	_, err := generateMap(config)
	assert.NoError(err)
}
