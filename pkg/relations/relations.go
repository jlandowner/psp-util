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

package relations

import (
	"context"
	"fmt"

	"github.com/jlandowner/psp-util/pkg/policy"
	"github.com/jlandowner/psp-util/pkg/rbac"
	"github.com/jlandowner/psp-util/pkg/utils"
	policyv1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/client-go/kubernetes"
)

type RelationalPodSecurityPolicy struct {
	ClusterRoles []RelationalClusterRole
	policyv1.PodSecurityPolicy
}

type RelationalClusterRole struct {
	ClusterRoleBindings []rbacv1.ClusterRoleBinding
	rbacv1.ClusterRole
}

func (r RelationalClusterRole) IsManaged() bool {
	_, ok := r.Annotations[utils.AnnotaionKeyPSPName]
	return ok
}

func GetRelationalPSPs(ctx context.Context, k8sclient *kubernetes.Clientset) ([]RelationalPodSecurityPolicy, error) {
	apiPspList, err := policy.ListPSP(ctx, k8sclient)
	if err != nil {
		return nil, fmt.Errorf("Failed to list PSP: %v", err.Error())
	}
	rpsps := make([]RelationalPodSecurityPolicy, len(apiPspList.Items))

	apiCrList, err := rbac.ListUsePSPRole(ctx, k8sclient)
	if err != nil {
		return nil, fmt.Errorf("Failed to list ClusterRole: %v", err.Error())
	}

	apiCrbList, err := rbac.ListClusterRoleBindings(ctx, k8sclient)
	if err != nil {
		return nil, fmt.Errorf("Failed to list ClusterRole: %v", err.Error())
	}

	for i, p := range apiPspList.Items {
		rpsps[i] = *generateRelationalPSP(&p, apiCrList, apiCrbList)
	}

	return rpsps, nil
}

func GetRelationalPSP(ctx context.Context, k8sclient *kubernetes.Clientset, name string) (*RelationalPodSecurityPolicy, error) {
	apiPsp, err := policy.GetPSP(ctx, k8sclient, name)
	if err != nil {
		return nil, fmt.Errorf("Failed to get PSP: %v", err.Error())
	}

	apiCrList, err := rbac.ListUsePSPRole(ctx, k8sclient)
	if err != nil {
		return nil, fmt.Errorf("Failed to list ClusterRole: %v", err.Error())
	}

	apiCrbList, err := rbac.ListClusterRoleBindings(ctx, k8sclient)
	if err != nil {
		return nil, fmt.Errorf("Failed to list ClusterRole: %v", err.Error())
	}

	return generateRelationalPSP(apiPsp, apiCrList, apiCrbList), nil
}

func generateRelationalPSP(apiPsp *policyv1.PodSecurityPolicy, apiCrList *rbacv1.ClusterRoleList, apiCrbList *rbacv1.ClusterRoleBindingList) *RelationalPodSecurityPolicy {
	rpsp := &RelationalPodSecurityPolicy{PodSecurityPolicy: *apiPsp}

	for _, apicr := range apiCrList.Items {
		pspNames := rbac.ExtractPSPFromGenericRole(apicr)
		for _, pspName := range pspNames {
			if rpsp.Name == pspName {
				rpsp.ClusterRoles = append(rpsp.ClusterRoles, RelationalClusterRole{ClusterRole: apicr})
			}
		}
	}

	for j, cr := range rpsp.ClusterRoles {
		for _, apicrb := range apiCrbList.Items {
			if apicrb.RoleRef.APIGroup == "rbac.authorization.k8s.io" && apicrb.RoleRef.Kind == "ClusterRole" && apicrb.RoleRef.Name == cr.Name {
				rpsp.ClusterRoles[j].ClusterRoleBindings = append(rpsp.ClusterRoles[j].ClusterRoleBindings, apicrb)
			}
		}
	}
	return rpsp
}
