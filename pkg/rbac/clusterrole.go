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

package rbac

import (
	"context"

	"github.com/jlandowner/psp-util/pkg/utils"
	policyv1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	APIGroup = "rbac.authorization.k8s.io"
)

func GetClusterRole(ctx context.Context, k8sclient *kubernetes.Clientset, name string) (*rbacv1.ClusterRole, error) {
	return k8sclient.RbacV1().ClusterRoles().Get(ctx, name, metav1.GetOptions{})
}

func CreateClusterRole(ctx context.Context, k8sclient *kubernetes.Clientset, clusterRole *rbacv1.ClusterRole) (*rbacv1.ClusterRole, error) {
	return k8sclient.RbacV1().ClusterRoles().Create(ctx, clusterRole, metav1.CreateOptions{})
}

func DeleteClusterRole(ctx context.Context, k8sclient *kubernetes.Clientset, name string) error {
	return k8sclient.RbacV1().ClusterRoles().Delete(ctx, name, metav1.DeleteOptions{})
}

func CreatePSPRole(ctx context.Context, k8sclient *kubernetes.Clientset, psp *policyv1.PodSecurityPolicy) (*rbacv1.ClusterRole, error) {
	clusterRole := &rbacv1.ClusterRole{
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups:     []string{"policy"},
				ResourceNames: []string{psp.Name},
				Resources:     []string{"podsecuritypolicies"},
				Verbs:         []string{"use"},
			},
		},
	}
	clusterRole.SetName(utils.GenerateName(psp.Name))
	clusterRole.SetAnnotations(utils.GenerateAnotations(psp.Name))

	return CreateClusterRole(ctx, k8sclient, clusterRole)
}

func ListUsePSPRole(ctx context.Context, k8sclient *kubernetes.Clientset) (*rbacv1.ClusterRoleList, error) {
	clusterRoleList, err := k8sclient.RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	pspClusterRoleList := clusterRoleList.DeepCopy()
	pspClusterRoleList.Items = make([]rbacv1.ClusterRole, 0)

	for _, cr := range clusterRoleList.Items {
		pspNames := ExtractPSPFromGenericRole(cr)
		if len(pspNames) != 0 {
			pspClusterRoleList.Items = append(pspClusterRoleList.Items, cr)
		}
	}
	return pspClusterRoleList, nil
}
