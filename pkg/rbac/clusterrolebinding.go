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
	"reflect"

	"github.com/jlandowner/psp-util/pkg/utils"
	policyv1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetClusterRoleBinding(ctx context.Context, k8sclient *kubernetes.Clientset, name string) (*rbacv1.ClusterRoleBinding, error) {
	return k8sclient.RbacV1().ClusterRoleBindings().Get(ctx, name, metav1.GetOptions{})
}

func CreateClusterRoleBinding(ctx context.Context, k8sclient *kubernetes.Clientset, clusterRoleBinding *rbacv1.ClusterRoleBinding) (*rbacv1.ClusterRoleBinding, error) {
	return k8sclient.RbacV1().ClusterRoleBindings().Create(ctx, clusterRoleBinding, metav1.CreateOptions{})
}

func UpdateClusterRoleBinding(ctx context.Context, k8sclient *kubernetes.Clientset, clusterRoleBinding *rbacv1.ClusterRoleBinding) (*rbacv1.ClusterRoleBinding, error) {
	return k8sclient.RbacV1().ClusterRoleBindings().Update(ctx, clusterRoleBinding, metav1.UpdateOptions{})
}

func ListClusterRoleBindings(ctx context.Context, k8sclient *kubernetes.Clientset) (*rbacv1.ClusterRoleBindingList, error) {
	return k8sclient.RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{})
}

func CreatePSPRoleBinding(ctx context.Context, k8sclient *kubernetes.Clientset, psp *policyv1.PodSecurityPolicy) (*rbacv1.ClusterRoleBinding, error) {
	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     utils.GenerateName(psp.Name),
		},
	}
	clusterRoleBinding.SetName(utils.GenerateName(psp.Name))
	clusterRoleBinding.SetAnnotations(utils.GenerateAnotations(psp.Name))

	return CreateClusterRoleBinding(ctx, k8sclient, clusterRoleBinding)
}

func AttachSubjectToClusterRoleBinding(clusterRoleBinding *rbacv1.ClusterRoleBinding, subject rbacv1.Subject) (hasGivenSubject bool) {
	hasGivenSubject = false
	for _, s := range clusterRoleBinding.Subjects {
		if reflect.DeepEqual(s, subject) {
			hasGivenSubject = true
		}
	}
	if !hasGivenSubject {
		clusterRoleBinding.Subjects = append(clusterRoleBinding.Subjects, subject)
	}
	return hasGivenSubject
}

func DetachSubjectToClusterRoleBinding(clusterRoleBinding *rbacv1.ClusterRoleBinding, subject rbacv1.Subject) (hasGivenSubject bool) {
	hasGivenSubject = false
	pos := -1
	for i, s := range clusterRoleBinding.Subjects {
		if reflect.DeepEqual(s, subject) {
			pos = i
		}
	}
	if pos > 0 {
		hasGivenSubject = true
		clusterRoleBinding.Subjects = append(clusterRoleBinding.Subjects[:pos], clusterRoleBinding.Subjects[pos+1:]...)
	}
	return hasGivenSubject
}
