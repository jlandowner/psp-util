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
	"context"
	"fmt"

	"github.com/jlandowner/psp-util/pkg/client"
	"github.com/jlandowner/psp-util/pkg/rbac"
	"github.com/jlandowner/psp-util/pkg/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use:   "clean PSP-NAME",
	Short: "Clean managed ClusterRole and ClusterRoleBinding",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Args is invalid. Required: `PSP-NAME`")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		k8sclient, err := client.NewClient(&kubeconfigPath)
		if err != nil {
			return fmt.Errorf("Failed to load kubeconfig %v: %v", kubeconfigPath, err.Error())
		}
		pspName := args[0]

		name := utils.GenerateName(pspName)
		err = rbac.DeleteClusterRoleBindings(ctx, k8sclient, name)
		if err != nil {
			return err
		}

		err = rbac.DeleteClusterRole(ctx, k8sclient, name)
		if err != nil {
			return err
		}
		return nil
	},
}
