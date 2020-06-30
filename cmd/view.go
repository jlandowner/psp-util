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
	"log"
	"os"

	"github.com/disiqueira/gotree"
	"github.com/jlandowner/psp-util/pkg/client"
	"github.com/jlandowner/psp-util/pkg/printers"
	"github.com/jlandowner/psp-util/pkg/structured"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(viewCmd)
}

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View relation between PSP and RBAC",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		k8sclient, err := client.NewClient(&kubeconfigPath)
		if err != nil {
			log.Fatalf("Failed to load kubeconfig %v (%v)", kubeconfigPath, err.Error())
		}

		psps, err := structured.GetStructuredPSPs(ctx, k8sclient)
		if err != nil {
			log.Fatalln(err)
		}

		w := os.Stdout
		for _, psp := range psps {
			pspTree := gotree.New(fmt.Sprintf("POD-SECURITY-POLICY {Name: "+printers.GreenString+"}", psp.Name))
			for _, cr := range psp.ClusterRoles {
				crTree := gotree.New(fmt.Sprintf("CLUSTER-ROLE {Name: "+printers.GreenString+"}", cr.Name))
				for _, crb := range cr.ClusterRoleBindings {
					crbTree := gotree.New(fmt.Sprintf("CLUSTER-ROLE-BINDING {Name: "+printers.GreenString+"}", crb.Name))
					for _, sub := range crb.Subjects {
						crbTree.Add(fmt.Sprintf("SUBJECT {Kind: "+printers.CianString+", Name: "+printers.RedString+", Namespace:"+printers.BlueString+"}", sub.Kind, sub.Name, sub.Namespace))
					}
					crTree.AddTree(crbTree)
				}
				pspTree.AddTree(crTree)
			}
			fmt.Fprintln(w, pspTree.Print())
		}

	},
}
