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
	"github.com/jlandowner/psp-util/pkg/policy"
	"github.com/jlandowner/psp-util/pkg/rbac"
	"github.com/jlandowner/psp-util/pkg/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(detachCmd)
	detachCmd.Flags().BoolVarP(&d.useGroup, "group", "g", false, "use Kind Group as Subject Kind")
	detachCmd.Flags().BoolVarP(&d.useUser, "user", "u", false, "use Kind User as Subject Kind")
	detachCmd.Flags().BoolVarP(&d.useSA, "sa", "s", false, "use Kind ServiceAccount as Subject Kind")
	detachCmd.Flags().StringVarP(&d.SubjectNamespace, "namespace", "n", "", "only used when kind is namedspaced resource(e.g. ServiceAccount)")
}

var (
	d = &AttachDetachOptions{}

	detachCmd = &cobra.Command{
		Use:               "detach PSP-NAME SUBJECT-NAME [ --group | --user | --sa ]",
		Short:             "Detach PSP from RBAC Subject",
		PersistentPreRunE: d.Complete,

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			k8sclient, err := client.NewClient(&kubeconfigPath)
			if err != nil {
				return fmt.Errorf("Failed to load kubeconfig %v: %v", kubeconfigPath, err.Error())
			}

			sub, err := d.GenerateSubject()
			if err != nil {
				return fmt.Errorf("Invalid options: %v", err.Error())
			}

			// Get PodSecurityPolicy
			psp, err := policy.GetPSP(ctx, k8sclient, d.PSPName)
			if err != nil {
				return fmt.Errorf("Failed to get PSP: %s", err.Error())
			}
			resourceName := utils.GenerateName(psp.Name)

			// Get ClusterRole
			cr, err := rbac.GetClusterRole(ctx, k8sclient, resourceName)
			if cr == nil || err != nil {
				fmt.Printf("Failed to get ClusterRole: %s\n", err.Error())
				return fmt.Errorf("psp '%s' has NOT been attached to %s. See `psp-util tree`", psp.Name, sub.String())
			}

			// Get ClusterRoleBinding
			crb, err := rbac.GetClusterRoleBinding(ctx, k8sclient, resourceName)
			if err != nil {
				fmt.Printf("Failed to get ClusterRoleBinding: %s\n", err.Error())
				return fmt.Errorf("psp '%s' has NOT been attached to %s. See `psp-util tree`", psp.Name, sub.String())
			}

			// Add Subject to ClusterRoleBinding
			hasGivenSubject := rbac.DetachSubjectToClusterRoleBinding(crb, *sub)
			if !hasGivenSubject {
				fmt.Printf("psp '%s' has NOT been attached to %s. See `psp-util tree`\n", psp.Name, sub.String())
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
