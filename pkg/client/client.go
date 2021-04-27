/*
Copyright 2020 jlandowner.
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

package client

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

// NewClient returns kubernetes Clientset
func NewClient(kubeconfigPath *string, kubecontext *string) (*kubernetes.Clientset, error) {
	if *kubeconfigPath == "" {
		*kubeconfigPath = filepath.Join(homeDir(), ".kube", "config")
	}

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: *kubeconfigPath},
		&clientcmd.ConfigOverrides{CurrentContext: *kubecontext}).ClientConfig()

	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func GetDefaultNamespace(kubeconfigPath *string) (string, error) {
	if *kubeconfigPath == "" {
		*kubeconfigPath = filepath.Join(homeDir(), ".kube", "config")
	}
	config, err := clientcmd.LoadFromFile(*kubeconfigPath)
	if err != nil {
		return "", err
	}

	currentContext, ok := config.Contexts[config.CurrentContext]
	if !ok {
		return "", fmt.Errorf("Failed to get currentcontext %s in kubeconfig %v", config.CurrentContext, *kubeconfigPath)
	}

	namespace := currentContext.Namespace
	if namespace == "" {
		return "default", nil
	}
	return namespace, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
