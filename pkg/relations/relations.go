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
	ClusterRoles []*RelationalClusterRole
	policyv1.PodSecurityPolicy
}

type RelationalClusterRole struct {
	ClusterRoleBindings []*rbacv1.ClusterRoleBinding
	rbacv1.ClusterRole
}

func (r RelationalClusterRole) IsManaged() bool {
	_, ok := r.Annotations[utils.AnnotaionKeyPSPName]
	return ok
}

func GetRelationalPSPs(ctx context.Context, k8sclient *kubernetes.Clientset) ([]RelationalPodSecurityPolicy, error) {
	psps, err := policy.ListPSP(ctx, k8sclient)
	if err != nil {
		return nil, fmt.Errorf("Failed to list PSP: %v", err.Error())
	}

	crs, err := rbac.ListClusterRolesWithPSP(ctx, k8sclient)
	if err != nil {
		return nil, fmt.Errorf("Failed to list ClusterRole: %v", err.Error())
	}

	crbs, err := rbac.ListClusterRoleBindings(ctx, k8sclient)
	if err != nil {
		return nil, fmt.Errorf("Failed to list ClusterRoleBindings: %v", err.Error())
	}

	rpsps := generateRelationalPSP(psps, crs, crbs)

	return rpsps, nil
}

func generateRelationalPSP(psps *policyv1.PodSecurityPolicyList, crs *rbacv1.ClusterRoleList, crbs *rbacv1.ClusterRoleBindingList) []RelationalPodSecurityPolicy {
	rpsps := make([]RelationalPodSecurityPolicy, len(psps.Items))

	rpspByName := make(map[string]*RelationalPodSecurityPolicy)
	for i, psp := range psps.Items {
		rpsp := RelationalPodSecurityPolicy{PodSecurityPolicy: psp}
		rpsps[i] = rpsp
		rpspByName[rpsp.Name] = &rpsps[i]
	}

	// build PSP to ClusterRole references
	crByName := make(map[string]*RelationalClusterRole)
	for _, cr := range crs.Items {
		pspNames := rbac.ExtractPSPFromGenericRole(cr)
		for _, pspName := range pspNames {
			if rpsp, ok := rpspByName[pspName]; ok {
				rcr := &RelationalClusterRole{ClusterRole: cr}
				rpsp.ClusterRoles = append(rpsp.ClusterRoles, rcr)
				crByName[cr.Name] = rcr
			}
		}
	}

	for i, crb := range crbs.Items {
		cr, ok := crByName[crb.RoleRef.Name]
		if !ok {
			continue
		}
		if crb.RoleRef.APIGroup != "rbac.authorization.k8s.io" || crb.RoleRef.Kind != "ClusterRole" {
			continue
		}

		cr.ClusterRoleBindings = append(cr.ClusterRoleBindings, &crbs.Items[i])
	}

	return rpsps
}
