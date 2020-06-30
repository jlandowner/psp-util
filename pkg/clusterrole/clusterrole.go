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

package clusterrole

import (
	"context"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ListUsePSPRole(ctx context.Context, k8sclient *kubernetes.Clientset) (*rbacv1.ClusterRoleList, error) {
	clusterRoleList, err := k8sclient.RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	pspClusterRoleList := clusterRoleList.DeepCopy()
	pspClusterRoleList.Items = make([]rbacv1.ClusterRole, 0)

	for _, cr := range clusterRoleList.Items {
		pspNames := ExtractPSPNamesFromClusterRole(cr)
		if len(pspNames) != 0 {
			pspClusterRoleList.Items = append(pspClusterRoleList.Items, cr)
		}
	}
	return pspClusterRoleList, nil
}

func ExtractPSPNamesFromClusterRole(cr rbacv1.ClusterRole) []string {
	pspNames := make([]string, 0)
	for _, rule := range cr.Rules {
		for i, apiGroups := range rule.APIGroups {
			if apiGroups == "policy" {
				for _, resource := range rule.Resources {
					if resource == "podsecuritypolicies" {
						for _, verb := range rule.Verbs {
							if verb == "use" {
								pspNames = append(pspNames, rule.ResourceNames[i])
							}
						}
					}
				}
			}
		}
	}
	return pspNames
}

func ListBindings(ctx context.Context, k8sclient *kubernetes.Clientset) (*rbacv1.ClusterRoleBindingList, error) {
	return k8sclient.RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{})
}
