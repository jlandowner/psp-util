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
	rbacv1 "k8s.io/api/rbac/v1"
)

func ExtractPSPFromGenericRole(r interface{}) []string {
	var rules []rbacv1.PolicyRule

	switch r := r.(type) {
	case rbacv1.ClusterRole:
		rules = r.Rules
	case rbacv1.Role:
		rules = r.Rules
	}

	pspNames := make([]string, 0)
	for _, rule := range rules {
		if hasAPIGroupsPolicy(rule) && hasResourcePSP(rule) && hasVerbUse(rule) {
			for _, resourceName := range rule.ResourceNames {
				pspNames = append(pspNames, resourceName)
			}
		}
	}
	return pspNames
}

func hasAPIGroupsPolicy(rule rbacv1.PolicyRule) bool {
	for _, apiGroups := range rule.APIGroups {
		if apiGroups == "policy" || apiGroups == "extensions" {
			return true
		}
	}
	return false
}

func hasResourcePSP(rule rbacv1.PolicyRule) bool {
	for _, resource := range rule.Resources {
		if resource == "podsecuritypolicies" {
			return true
		}
	}
	return false
}

func hasVerbUse(rule rbacv1.PolicyRule) bool {
	for _, verb := range rule.Verbs {
		if verb == "use" {
			return true
		}
	}
	return false
}
