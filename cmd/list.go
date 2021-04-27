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

	"github.com/jlandowner/psp-util/cmd/options"
	"github.com/jlandowner/psp-util/pkg/client"
	"github.com/jlandowner/psp-util/pkg/printers"
	"github.com/jlandowner/psp-util/pkg/relations"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&l.NoHeader, "no-headers", false, "output without header")
	listCmd.Flags().BoolVarP(&l.ClusterRole, "cluster-role", "c", false, "output only clusterroles associated with PSP")
	listCmd.Flags().BoolVarP(&l.Role, "role", "r", false, "output only roles associated with PSP")
}

var (
	l = &options.ListOptions{}

	listCmd = &cobra.Command{
		Use:               "list",
		Short:             "List PSPs and the related RBACs in cluster",
		PersistentPreRunE: l.PreRunE,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			k8sclient, err := client.NewClient(&kubeconfigPath, &kubecontext)
			if err != nil {
				return fmt.Errorf("Failed to load kubeconfig %v: %v", kubeconfigPath, err.Error())
			}

			psps, err := relations.GetRelationalPSPs(ctx, k8sclient)
			if err != nil {
				return err
			}

			printOpt := printers.ListPrinterOptions{true, true, true, true, true, true}
			if l.ClusterRole {
				printOpt.Role = false
				printOpt.PSPUtilManaged = false
			}
			if l.Role {
				printOpt.ClusterRole = false
				printOpt.ClusterRoleBinding = false
				printOpt.PSPUtilManaged = false
			}

			printer := printers.NewListPrinter(os.Stdout, printOpt)
			defer printer.Flush()

			if !l.NoHeader {
				printer.PrintHeader()
			}

			for _, psp := range psps {
				if len(psp.ClusterRoles) == 0 && len(psp.Roles) == 0 {
					printer.PrintLine(printers.ListPrinterLine{})
					continue
				}

				if printOpt.ClusterRole {
					// clusterrole can bind to either clusterrolebinding or rolebinding
					for _, cr := range psp.ClusterRoles {
						if len(cr.ClusterRoleBindings) == 0 && len(cr.RoleBindings) == 0 {
							printer.PrintLine(printers.ListPrinterLine{
								PSP:            psp.Name,
								ClusterRole:    cr.Name,
								PSPUtilManaged: strconv.FormatBool(cr.IsManaged())})
							continue
						}
						for _, crb := range cr.ClusterRoleBindings {
							printer.PrintLine(printers.ListPrinterLine{
								PSP:                psp.Name,
								ClusterRole:        cr.Name,
								ClusterRoleBinding: crb.Name,
								PSPUtilManaged:     strconv.FormatBool(cr.IsManaged())})
						}
						for _, rb := range cr.RoleBindings {
							rbname := fmt.Sprintf("%v/%v", rb.Namespace, rb.Name)
							printer.PrintLine(printers.ListPrinterLine{
								PSP:            psp.Name,
								ClusterRole:    cr.Name,
								RoleBinding:    rbname,
								PSPUtilManaged: strconv.FormatBool(false)})
						}
					}
				}

				if printOpt.Role {
					// role can only bind to rolebinding
					for _, r := range psp.Roles {
						rname := fmt.Sprintf("%v/%v", r.Namespace, r.Name)
						if len(r.RoleBindings) == 0 {
							printer.PrintLine(printers.ListPrinterLine{
								PSP:            psp.Name,
								Role:           rname,
								PSPUtilManaged: strconv.FormatBool(false)})
							continue
						}
						for _, rb := range r.RoleBindings {
							rbname := fmt.Sprintf("%v/%v", r.Namespace, rb.Name)
							printer.PrintLine(printers.ListPrinterLine{
								PSP:            psp.Name,
								Role:           rname,
								RoleBinding:    rbname,
								PSPUtilManaged: strconv.FormatBool(false)})
						}
					}
				}

			}
			return nil
		},
	}
)
