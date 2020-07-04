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

	"github.com/jlandowner/psp-util/cmd/options"
	"github.com/jlandowner/psp-util/pkg/client"
	"github.com/jlandowner/psp-util/pkg/policy"
	"github.com/jlandowner/psp-util/pkg/rbac"
	"github.com/jlandowner/psp-util/pkg/utils"
	"github.com/spf13/cobra"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
)

func init() {
	rootCmd.AddCommand(attachCmd)
	attachCmd.Flags().StringVarP(&a.Group, "group", "g", "", "set Subject's Name and use Kind Group")
	attachCmd.Flags().StringVarP(&a.User, "user", "u", "", "set Subject's Name and use Kind User")
	attachCmd.Flags().StringVarP(&a.ServiceAccount, "sa", "s", "", "set Subject's Name and use Kind ServiceAccount")

	attachCmd.Flags().StringVar(&a.SubjectKind, "kind", "", "set Subject's Kind")
	attachCmd.Flags().StringVar(&a.SubjectName, "name", "", "set Subject's Name")
	attachCmd.Flags().StringVar(&a.SubjectAPIGroup, "api-group", "", "set Subject's APIGroup")

	attachCmd.Flags().StringVarP(&a.SubjectNamespace, "namespace", "n", "", "set Subject's Namespace (only used when kind is ServiceAccount)")
}

var (
	a = &options.AttachDetachOptions{}

	attachCmd = &cobra.Command{
		Use:               "attach PSP-NAME [ --group | --user | --sa ] SUBJECT-NAME",
		Short:             "Attach PSP to RBAC Subject (Auto generate managed ClusterRole and ClusterRoleBinding)",
		PersistentPreRunE: a.PreRunE,

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			k8sclient, err := client.NewClient(&kubeconfigPath)
			if err != nil {
				return fmt.Errorf("Failed to load kubeconfig %v: %v", kubeconfigPath, err.Error())
			}

			sub, err := a.GenerateSubject(&kubeconfigPath)
			if err != nil {
				return fmt.Errorf("Invalid options: %v", err.Error())
			}

			// Get PodSecurityPolicy
			psp, err := policy.GetPSP(ctx, k8sclient, a.PSPName)
			if apierrs.IsNotFound(err) {
				return fmt.Errorf("PSP %s is not found. See `psp-util tree`", a.PSPName)
			}
			if err != nil {
				return fmt.Errorf("Failed to get PSP: %s", err.Error())
			}
			resourceName := utils.GenerateName(psp.Name)

			// Get or Create ClusterRole
			cr, err := rbac.GetClusterRole(ctx, k8sclient, resourceName)
			if apierrs.IsNotFound(err) {
				fmt.Printf("Managed ClusterRole is not found...")
				cr, err = rbac.CreatePSPRole(ctx, k8sclient, psp)
				if err != nil {
					return fmt.Errorf("Failed to create ClusterRole: %s", err.Error())
				}
				fmt.Printf("Created\n")
			}
			if cr == nil || err != nil {
				return fmt.Errorf("Failed to get ClusterRole: %s", err.Error())
			}

			// Get or Create ClusterRoleBinding
			crb, err := rbac.GetClusterRoleBinding(ctx, k8sclient, resourceName)
			if apierrs.IsNotFound(err) {
				fmt.Printf("Managed ClusterRoleBinding is not found...")
				crb, err = rbac.CreatePSPRoleBinding(ctx, k8sclient, psp)
				if err != nil {
					return fmt.Errorf("Failed to create ClusterRoleBinding: %s", err.Error())
				}
				fmt.Printf("Created\n")
			}
			if err != nil {
				return fmt.Errorf("Failed to get ClusterRoleBinding: %s", err.Error())
			}
			// Add Subject to ClusterRoleBinding
			if rbac.AttachSubjectToClusterRoleBinding(crb, *sub) {
				fmt.Printf("psp '%s' has already been attached to %s. See `psp-util tree`\n", psp.Name, sub.String())
				return nil
			}

			// Update ClusterRoleBinding to attach subjects
			_, err = rbac.UpdateClusterRoleBinding(ctx, k8sclient, crb)
			if err != nil {
				return fmt.Errorf("Failed to update ClusterRoleBinding: %s", err.Error())
			}
			return nil
		},
	}
)
