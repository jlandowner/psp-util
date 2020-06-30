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

package structured

import (
	"context"
	"fmt"

	"github.com/jlandowner/psp-util/pkg/clusterrole"
	"github.com/jlandowner/psp-util/pkg/psp"
	policyv1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/client-go/kubernetes"
)

type StructuredPodSecurityPolicy struct {
	ClusterRoles []StructuredClusterRole
	policyv1.PodSecurityPolicy
}

type StructuredClusterRole struct {
	ClusterRoleBindings []rbacv1.ClusterRoleBinding
	rbacv1.ClusterRole
}

func GetStructuredPSPs(ctx context.Context, k8sclient *kubernetes.Clientset) ([]StructuredPodSecurityPolicy, error) {
	apiPspList, err := psp.List(ctx, k8sclient)
	if err != nil {
		return nil, fmt.Errorf("Failed to list PSP: %v", err.Error())
	}
	spsps := make([]StructuredPodSecurityPolicy, len(apiPspList.Items))

	apiCrList, err := clusterrole.ListUsePSPRole(ctx, k8sclient)
	if err != nil {
		return nil, fmt.Errorf("Failed to list ClusterRole: %v", err.Error())
	}

	apiCrbList, err := clusterrole.ListBindings(ctx, k8sclient)
	if err != nil {
		return nil, fmt.Errorf("Failed to list ClusterRole: %v", err.Error())
	}

	for i, p := range apiPspList.Items {
		spsps[i] = *generateStructuredPSP(&p, apiCrList, apiCrbList)
	}

	return spsps, nil
}

func GetStructuredPSP(ctx context.Context, k8sclient *kubernetes.Clientset, name string) (*StructuredPodSecurityPolicy, error) {
	apiPsp, err := psp.Get(ctx, k8sclient, name)
	if err != nil {
		return nil, fmt.Errorf("Failed to get PSP: %v", err.Error())
	}

	apiCrList, err := clusterrole.ListUsePSPRole(ctx, k8sclient)
	if err != nil {
		return nil, fmt.Errorf("Failed to list ClusterRole: %v", err.Error())
	}

	apiCrbList, err := clusterrole.ListBindings(ctx, k8sclient)
	if err != nil {
		return nil, fmt.Errorf("Failed to list ClusterRole: %v", err.Error())
	}

	return generateStructuredPSP(apiPsp, apiCrList, apiCrbList), nil
}

func generateStructuredPSP(apiPsp *policyv1.PodSecurityPolicy, apiCrList *rbacv1.ClusterRoleList, apiCrbList *rbacv1.ClusterRoleBindingList) *StructuredPodSecurityPolicy {
	spsp := &StructuredPodSecurityPolicy{PodSecurityPolicy: *apiPsp}

	for _, apicr := range apiCrList.Items {
		pspNames := clusterrole.ExtractPSPNamesFromClusterRole(apicr)
		for _, pspName := range pspNames {
			if spsp.Name == pspName {
				spsp.ClusterRoles = append(spsp.ClusterRoles, StructuredClusterRole{ClusterRole: apicr})
			}
		}
	}

	for j, cr := range spsp.ClusterRoles {
		for _, apicrb := range apiCrbList.Items {
			if apicrb.RoleRef.APIGroup == "rbac.authorization.k8s.io" && apicrb.RoleRef.Kind == "ClusterRole" && apicrb.RoleRef.Name == cr.Name {
				spsp.ClusterRoles[j].ClusterRoleBindings = append(spsp.ClusterRoles[j].ClusterRoleBindings, apicrb)
			}
		}
	}
	return spsp
}
