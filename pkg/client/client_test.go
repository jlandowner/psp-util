package client

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestGetDefaultNamespace(t *testing.T) {
	tests := []struct {
		title      string
		kubeconfig string
		expect     string
	}{
		{
			title:      "test1",
			kubeconfig: "",
			expect:     getCurrentNamespaceInDefaultKubeconfig(),
		},
		{
			title:      "test2",
			kubeconfig: "../../test/config",
			expect:     "test-a",
		},
	}

	for _, test := range tests {
		t.Log(test.title, test.kubeconfig, test.expect)
		kubeconfig := test.kubeconfig
		ns, err := GetDefaultNamespace(&kubeconfig)
		assert.Nil(t, err)
		assert.Equal(t, test.expect, ns)
	}
}

func getCurrentNamespaceInDefaultKubeconfig() string {
	config, err := readKubeconfig(homeDir() + "/.kube/config")
	if err != nil {
		panic(err)
	}
	numContext := -1
	for i, v := range config.Contexts {
		if config.CurrentContext == v.Name {
			numContext = i
			break
		}
	}
	if numContext < 0 {
		panic(config.Contexts)
	}
	return config.Contexts[numContext].Context.Namespace
}

type kubeconfig struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string
	Clusters   []interface{}
	Contexts   []struct {
		Context struct {
			Cluster   string
			Namespace string
			User      string
		}
		Name string
	}
	CurrentContext string `yaml:"current-context"`
	Preferences    interface{}
	Users          interface{}
}

func readKubeconfig(kubeconfigPath string) (*kubeconfig, error) {
	config := &kubeconfig{}
	buf, err := ioutil.ReadFile(kubeconfigPath)
	if err != nil {
		return nil, err
	}

	// yaml to struct
	err = yaml.Unmarshal(buf, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
