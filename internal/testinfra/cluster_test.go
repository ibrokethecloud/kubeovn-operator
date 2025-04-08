package testinfra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func Test_CreateAndVerifyCluster(t *testing.T) {
	assert := require.New(t)
	ctx := context.TODO()
	err := CreateCluster(ctx, defaultClusterName)
	assert.NoError(err, "expected no error during cluster creation")
	cfgBytes, err := GetKubeconfig(ctx, defaultClusterName)
	assert.NoError(err, "expected no error during kubeconfig fetch")
	restConfig, err := clientcmd.RESTConfigFromKubeConfig(cfgBytes)
	assert.NoError(err, "expected no error during restconfig generation")
	client, err := kubernetes.NewForConfig(restConfig)
	nodes, err := client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	assert.NoError(err, "expected no error while fetching nodes")
	assert.Len(nodes.Items, 3, "expected to find 3 nodes in the test cluster")
}

func Test_DeleteCluster(t *testing.T) {
	assert := require.New(t)
	ctx := context.TODO()
	err := DeleteCluster(ctx, defaultClusterName)
	assert.NoError(err, "expected no error during cluster deletion")
}

func Test_GetKubeconfig(t *testing.T) {
	assert := require.New(t)
	ctx := context.TODO()
	cfgBytes, err := GetKubeconfig(ctx, defaultClusterName)
	assert.NoError(err, "expected no error during kubeconfig fetch")
	restConfig, err := clientcmd.RESTConfigFromKubeConfig(cfgBytes)
	assert.NoError(err, "expected no error during restconfig generation")
	client, err := kubernetes.NewForConfig(restConfig)
	nodes, err := client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	assert.NoError(err, "expected no error while fetching nodes")
	for _, node := range nodes.Items {
		t.Log(node.Name)
	}

}
