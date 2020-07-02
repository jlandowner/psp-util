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
	rbacv1 "k8s.io/api/rbac/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
)

func init() {
	rootCmd.AddCommand(attachCmd)
	attachCmd.Flags().BoolVarP(&a.useGroup, "group", "g", false, "use Kind Group as Subject's Kind")
	attachCmd.Flags().BoolVarP(&a.useUser, "user", "u", false, "use Kind User as Subject's Kind")
	attachCmd.Flags().BoolVarP(&a.useSA, "sa", "s", false, "use Kind ServiceAccount as Subject's Kind")
	attachCmd.Flags().StringVarP(&a.SubjectNamespace, "namespace", "n", "", "set Subject's Namespace (only used when kind is ServiceAccount)")
}

var (
	subjectKindList = []string{"Group", "User", "ServiceAccount"}
	a               = &AttachDetachOptions{}

	attachCmd = &cobra.Command{
		Use:               "attach PSP-NAME SUBJECT-NAME [ --group | --user | --sa ]",
		Short:             "Attach PSP to RBAC Subject",
		PersistentPreRunE: a.Complete,

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			k8sclient, err := client.NewClient(&kubeconfigPath)
			if err != nil {
				return fmt.Errorf("Failed to load kubeconfig %v: %v", kubeconfigPath, err.Error())
			}

			sub, err := a.GenerateSubject()
			if err != nil {
				return fmt.Errorf("Invalid options: %v", err.Error())
			}

			// Get PodSecurityPolicy
			psp, err := policy.GetPSP(ctx, k8sclient, a.PSPName)
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

type AttachDetachOptions struct {
	PSPName          string
	SubjectKind      string
	SubjectName      string
	SubjectNamespace string

	useGroup bool
	useUser  bool
	useSA    bool
}

func (a AttachDetachOptions) Complete(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("Args is invalid. Required: `PSP-NAME SUBJECT-NAME`")
	}
	a.PSPName = args[0]
	a.SubjectName = args[1]

	if !a.useGroup && !a.useUser && !a.useSA {
		return fmt.Errorf("You must specify one of Subject Kind flag in %v", subjectKindList)
	}
	if a.useGroup {
		a.SubjectKind = "Group"
	}
	if a.useUser {
		if a.SubjectKind != "" {
			return fmt.Errorf("Multiple kind is not allowed. You must specify one of Subject Kind flag in %v", subjectKindList)
		}
		a.SubjectKind = "User"
	}
	if a.useSA {
		if a.SubjectKind != "" {
			return fmt.Errorf("Multiple kind is not allowed. You must specify one of Subject Kind flag in %v", subjectKindList)
		}
		a.SubjectKind = "ServiceAccount"
	}
	return nil
}

func (a AttachDetachOptions) GenerateSubject() (*rbacv1.Subject, error) {
	sub := &rbacv1.Subject{}
	sub.Kind = a.SubjectKind
	sub.Name = a.SubjectName

	switch sub.Kind {
	case "ServiceAccount":
		if a.SubjectNamespace == "" {
			return nil, fmt.Errorf("namespace is requireed when kind is ServiceAccount")
		}
		sub.Namespace = a.SubjectNamespace

	case "Group":
		sub.APIGroup = "rbac.authorization.k8s.io"

	case "User":
		sub.APIGroup = "rbac.authorization.k8s.io"

	default:
		return nil, fmt.Errorf("Kind is allowed in %v", subjectKindList)
	}
	return sub, nil
}
