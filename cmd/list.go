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
	"strconv"

	"github.com/jlandowner/psp-util/pkg/client"
	"github.com/jlandowner/psp-util/pkg/printers"
	"github.com/jlandowner/psp-util/pkg/relations"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List PSPs and the related RBACs in cluster",
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

		w := printers.GetNewTabWriter(os.Stdout)
		defer w.Flush()

		printers.PrintLine(w, []string{"PSP-NAME", "CLUSTER-ROLE", "CLUSTER-ROLE-BINDING", "NS/ROLE", "NS/ROLE-BINDING", "PSP-UTIL-MANAGED"})

		for _, psp := range psps {
			if len(psp.ClusterRoles) == 0 && len(psp.Roles) == 0 {
				printers.PrintLine(w, []string{psp.Name, "", "", "", "", ""})
				continue
			}
			// clusterrole can bind to either clusterrolebinding or rolebinding
			for _, cr := range psp.ClusterRoles {
				if len(cr.ClusterRoleBindings) == 0 && len(cr.RoleBindings) == 0 {
					printers.PrintLine(w, []string{psp.Name, cr.Name, "", "", "", strconv.FormatBool(cr.IsManaged())})
					continue
				}
				for _, crb := range cr.ClusterRoleBindings {
					printers.PrintLine(w, []string{psp.Name, cr.Name, crb.Name, "", "", strconv.FormatBool(cr.IsManaged())})
				}
				for _, rb := range cr.RoleBindings {
					rbname := fmt.Sprintf("%v/%v", rb.Namespace, rb.Name)
					printers.PrintLine(w, []string{psp.Name, cr.Name, "", "", rbname, strconv.FormatBool(false)})
				}
			}
			// role can only bind to rolebinding
			for _, r := range psp.Roles {
				rname := fmt.Sprintf("%v/%v", r.Namespace, r.Name)
				if len(r.RoleBindings) == 0 {
					printers.PrintLine(w, []string{psp.Name, "", "", rname, "", strconv.FormatBool(false)})
					continue
				}
				for _, rb := range r.RoleBindings {
					rbname := fmt.Sprintf("%v/%v", r.Namespace, rb.Name)
					printers.PrintLine(w, []string{psp.Name, "", "", rname, rbname, strconv.FormatBool(false)})
				}
			}
		}
		return nil
	},
}
