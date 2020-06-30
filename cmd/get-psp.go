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
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/jlandowner/psp-util/pkg/client"
	getpsp "github.com/jlandowner/psp-util/pkg/cmd/get-psp"
	"github.com/jlandowner/psp-util/pkg/printers"
)

func init() {
	getCmd.AddCommand(getPspCmd)
}

var getPspCmd = &cobra.Command{
	Use:   "psp",
	Short: "Get Pod Security Policies and RBAC associated with it.",
	Run: func(cmd *cobra.Command, args []string) {

		ctx := context.Background()
		k8sclient, err := client.NewClient(&kubeconfigPath)
		if err != nil {
			log.Fatalf("Failed to load kubeconfig %v (%v)", kubeconfigPath, err.Error())
		}

		psps, err := getpsp.GetPSP(ctx, k8sclient)
		if err != nil {
			log.Fatalln(err)
		}

		w := printers.GetNewTabWriter(os.Stdout)
		defer w.Flush()

		columnNames := []string{"PSP-NAME", "CLUSTER-ROLE", "CLUSTER-ROLE-BINDING"}
		printers.PrintLine(w, columnNames)

		for _, psp := range psps {
			if len(psp.ClusterRoles) == 0 {
				printers.PrintLine(w, []string{psp.Name})
				continue
			}
			for _, cr := range psp.ClusterRoles {
				if len(cr.ClusterRoleBindings) == 0 {
					printers.PrintLine(w, []string{psp.Name, cr.Name})
					continue
				}
				for _, crb := range cr.ClusterRoleBindings {
					printers.PrintLine(w, []string{psp.Name, cr.Name, crb.Name})
				}
			}
		}

	},
}
