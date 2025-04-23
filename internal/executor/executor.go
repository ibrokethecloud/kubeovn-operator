package executor

import (
	"context"
	"fmt"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/cmd/exec"
	"k8s.io/kubectl/pkg/scheme"
)

type RemoteCommandExecutor struct {
	client *kubernetes.Clientset
	pod    *corev1.Pod
	cfg    *rest.Config
	ctx    context.Context
}

// NewRemoteCommandExecutor is an implementation of Executor that runs commands in the driver pod
// which allows us to ship custom drivers as container images
func NewRemoteCommandExecutor(ctx context.Context, config *rest.Config, pod *corev1.Pod) (*RemoteCommandExecutor, error) {
	cfgCopy := *config
	cfgCopy.GroupVersion = &schema.GroupVersion{Group: "", Version: "v1"}
	cfgCopy.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: scheme.Codecs}
	cfgCopy.APIPath = "/api"
	client, err := kubernetes.NewForConfig(&cfgCopy)
	if err != nil {
		return nil, fmt.Errorf("error generating client for config in remote command executor: %v", err)
	}

	r := &RemoteCommandExecutor{
		client: client,
		cfg:    &cfgCopy,
		ctx:    ctx,
		pod:    pod,
	}
	return r, nil
}

func (r *RemoteCommandExecutor) Run(containerName string, cmd string) ([]byte, error) {
	iostreams, _, outBuffer, errBuffer := genericclioptions.NewTestIOStreams()
	streamOpts := exec.StreamOptions{
		Namespace:     r.pod.Namespace,
		PodName:       r.pod.Name,
		ContainerName: containerName,
		IOStreams:     iostreams,
		TTY:           false,
		Quiet:         false,
		Stdin:         false,
	}

	options := &exec.ExecOptions{
		StreamOptions: streamOpts,
		PodClient:     r.client.CoreV1(),
		Config:        r.cfg,
		Executor:      &exec.DefaultRemoteExecutor{},
	}

	options.Command = []string{"/bin/sh", cmd}
	err := options.Run()
	if err != nil {
		return errBuffer.Bytes(), fmt.Errorf("error during command execution: %v", err)
	}
	return outBuffer.Bytes(), nil
}
