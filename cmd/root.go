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

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	kubeconfigPath string

	rootCmd = &cobra.Command{
		Use:   "psp-util",
		Short: "Utility to manage Pod Security Policy(PSP) and the related RBAC Resources",
		Long: `
A Kubectl plugin to manage Pod Security Policy(PSP) and the related RBAC Resources.
Attach/Detach PSP to/from RBACs(Group, User) or ServiceAccounts and
view the relations which PSP is effected to the Subjects in cluster.

Complete documentation is available at http://github.com/jlandowner/psp-util

Apache License Version 2.0 Copyright 2020 jlandowner
`,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&kubeconfigPath, "kubeconfig", "", "kube config file (default is $HOME/.kube/config)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
