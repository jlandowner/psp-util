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

	"github.com/spf13/cobra"
)

func init() {
	getCmd.AddCommand(getSaCmd)
	getSaCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "namespace (default is in kubeconfig)")
	getSaCmd.PersistentFlags().BoolVarP(&isAllNamespace, "all-namespaces", "A", false, "use all namespace (default is false)")
}

var (
	namespace      string
	isAllNamespace bool
	getSaCmd       = &cobra.Command{
		Use:   "sa",
		Short: "Get Service Accounts and PSP associated with it.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("namespace %v\n", namespace)
			fmt.Printf("isAllNamespace %v\n", isAllNamespace)
			fmt.Println("get sa")
		},
	}
)
