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

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ListRolesWithPSP(ctx context.Context, k8sclient *kubernetes.Clientset) (*rbacv1.RoleList, error) {
	roleList, err := k8sclient.RbacV1().Roles("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	pspRoleList := roleList.DeepCopy()
	pspRoleList.Items = make([]rbacv1.Role, 0)

	for _, r := range roleList.Items {
		pspNames := ExtractPSPFromGenericRole(r)
		if len(pspNames) != 0 {
			pspRoleList.Items = append(pspRoleList.Items, r)
		}
	}
	return pspRoleList, nil
}

func ListRoleBindings(ctx context.Context, k8sclient *kubernetes.Clientset) (*rbacv1.RoleBindingList, error) {
	return k8sclient.RbacV1().RoleBindings("").List(ctx, metav1.ListOptions{})
}
