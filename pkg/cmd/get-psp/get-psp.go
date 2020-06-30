package getpsp

import (
	"context"
	"fmt"

	policyv1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"

	"github.com/jlandowner/psp-util/pkg/clusterrole"
	"github.com/jlandowner/psp-util/pkg/psp"
	"k8s.io/client-go/kubernetes"
)

type PodSecurityPolicy struct {
	ClusterRoles []ClusterRole
	policyv1.PodSecurityPolicy
}

type ClusterRole struct {
	ClusterRoleBindings []rbacv1.ClusterRoleBinding
	rbacv1.ClusterRole
}

func GetPSP(ctx context.Context, k8sclient *kubernetes.Clientset) ([]PodSecurityPolicy, error) {
	apiPspList, err := psp.List(ctx, k8sclient)
	if err != nil {
		return nil, fmt.Errorf("Failed to list PSP %v", err.Error())
	}
	psps := make([]PodSecurityPolicy, len(apiPspList.Items))
	for i, p := range apiPspList.Items {
		psps[i].PodSecurityPolicy = p
	}

	apiCrList, err := clusterrole.ListUsePSPRole(ctx, k8sclient)
	if err != nil {
		return nil, fmt.Errorf("Failed to list ClusterRole %v", err.Error())
	}

	apiCrbList, err := clusterrole.ListBindings(ctx, k8sclient)
	if err != nil {
		return nil, fmt.Errorf("Failed to list ClusterRole %v", err.Error())
	}

	for i := range psps {
		for _, apicr := range apiCrList.Items {
			pspNames := clusterrole.ExtractPSPNamesFromClusterRole(apicr)
			for _, pspName := range pspNames {
				if psps[i].Name == pspName {
					psps[i].ClusterRoles = append(psps[i].ClusterRoles, ClusterRole{ClusterRole: apicr})
				}
			}
		}

		for j, cr := range psps[i].ClusterRoles {
			for _, apicrb := range apiCrbList.Items {
				if apicrb.RoleRef.APIGroup == "rbac.authorization.k8s.io" && apicrb.RoleRef.Kind == "ClusterRole" && apicrb.RoleRef.Name == cr.Name {
					psps[i].ClusterRoles[j].ClusterRoleBindings = append(psps[i].ClusterRoles[j].ClusterRoleBindings, apicrb)
				}
			}
		}
	}

	return psps, nil
}
