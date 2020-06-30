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

package psp

import (
	"context"

	policyv1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func List(ctx context.Context, k8sclient *kubernetes.Clientset) (*policyv1.PodSecurityPolicyList, error) {
	return k8sclient.PolicyV1beta1().PodSecurityPolicies().List(ctx, metav1.ListOptions{})
}

func Get(ctx context.Context, k8sclient *kubernetes.Clientset, name string) (*policyv1.PodSecurityPolicy, error) {
	return k8sclient.PolicyV1beta1().PodSecurityPolicies().Get(ctx, name, metav1.GetOptions{})
}
