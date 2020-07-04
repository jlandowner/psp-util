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
	"os"

	"github.com/disiqueira/gotree"
	"github.com/jlandowner/psp-util/pkg/client"
	"github.com/jlandowner/psp-util/pkg/printers"
	"github.com/jlandowner/psp-util/pkg/relations"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(treeCmd)
}

var treeCmd = &cobra.Command{
	Use:   "tree",
	Short: "View relational tree between PSP and Subjects",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		k8sclient, err := client.NewClient(&kubeconfigPath)
		if err != nil {
			return fmt.Errorf("Failed to load kubeconfig %v: %v", kubeconfigPath, err.Error())
		}

		psps, err := relations.GetRelationalPSPs(ctx, k8sclient)
		if err != nil {
			return err
		}

		w := os.Stdout
		for _, psp := range psps {
			pspTree := gotree.New(fmt.Sprintf("📙PSP "+printers.GreenString, psp.Name))
			for _, cr := range psp.ClusterRoles {
				crTree := gotree.New(fmt.Sprintf("📕ClusterRole "+printers.GreenString, cr.Name))
				for _, crb := range cr.ClusterRoleBindings {
					crbTree := gotree.New(fmt.Sprintf("📘ClusterRoleBinding "+printers.GreenString, crb.Name))
					for _, sub := range crb.Subjects {
						crbTree.Add(fmt.Sprintf("📗Subject{Kind: "+printers.CianString+", Name: "+printers.RedString+", Namespace: "+printers.BlueString+"}", sub.Kind, sub.Name, sub.Namespace))
					}
					crTree.AddTree(crbTree)
				}
				pspTree.AddTree(crTree)
			}
			fmt.Fprintln(w, pspTree.Print())
		}
		return nil

	},
}
