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
	rootCmd.AddCommand(detachCmd)
	detachCmd.Flags().StringVarP(&d.Group, "group", "g", "", "set Subject's Name and use Kind Group")
	detachCmd.Flags().StringVarP(&d.User, "user", "u", "", "set Subject's Name and use Kind User")
	detachCmd.Flags().StringVarP(&d.ServiceAccount, "sa", "s", "", "set Subject's Name and use Kind ServiceAccount")

	detachCmd.Flags().StringVar(&d.SubjectKind, "kind", "", "set Subject's Kind")
	detachCmd.Flags().StringVar(&d.SubjectName, "name", "", "set Subject's Name")
	detachCmd.Flags().StringVar(&d.SubjectAPIGroup, "api-group", "", "set Subject's APIGroup")

	detachCmd.Flags().StringVarP(&d.SubjectNamespace, "namespace", "n", "", "only used when kind is namedspaced resource(e.g. ServiceAccount)")
}

var (
	d = &options.AttachDetachOptions{}

	detachCmd = &cobra.Command{
		Use:               "detach PSP-NAME [ --group | --user | --sa ] SUBJECT-NAME",
		Short:             "Detach PSP from RBAC Subject",
		PersistentPreRunE: d.PreRunE,

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			k8sclient, err := client.NewClient(&kubeconfigPath)
			if err != nil {
				return fmt.Errorf("Failed to load kubeconfig %v: %v", kubeconfigPath, err.Error())
			}

			sub, err := d.GenerateSubject(&kubeconfigPath)
			if err != nil {
				return fmt.Errorf("Invalid options: %v", err.Error())
			}

			// Get PodSecurityPolicy
			psp, err := policy.GetPSP(ctx, k8sclient, d.PSPName)
			if apierrs.IsNotFound(err) {
				return fmt.Errorf("PSP %s is not found. See `psp-util tree`", a.PSPName)
			}
			if err != nil {
				return fmt.Errorf("Failed to get PSP: %s", err.Error())
			}
			resourceName := utils.GenerateName(psp.Name)

			// Get ClusterRole
			cr, err := rbac.GetClusterRole(ctx, k8sclient, resourceName)
			if apierrs.IsNotFound(err) {
				return fmt.Errorf("Managed ClusterRole is not found. Please remove subjects manually from the ClusterRoleBindings. See the resources by `psp-util tree`")
			}
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
