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

package options

import (
	"fmt"
	"reflect"

	"github.com/jlandowner/psp-util/pkg/client"
	"github.com/jlandowner/psp-util/pkg/rbac"
	"github.com/spf13/cobra"
	rbacv1 "k8s.io/api/rbac/v1"
)

var (
	subjectKindFlags = "[ --group | --user | --sa ]"
	subjectKindList  = []string{"Group", "User", "ServiceAccount"}
)

type AttachDetachOptions struct {
	PSPName          string
	SubjectKind      string
	SubjectName      string
	SubjectNamespace string
	SubjectAPIGroup  string

	// Same field name as kind in `subjectKindList`
	Group          string
	User           string
	ServiceAccount string
}

func (o *AttachDetachOptions) PreRunE(cmd *cobra.Command, args []string) error {
	if err := o.Validate(cmd, args); err != nil {
		return err
	}
	if err := o.Complete(cmd, args); err != nil {
		return err
	}
	return nil
}

func (o *AttachDetachOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Args is invalid. Required: `PSP-NAME`")
	}
	_, _, kindFlagCount := getValuesFromKindFlags(o)

	if use(o.SubjectKind) {
		if kindFlagCount != 0 {
			return fmt.Errorf("Using both --kind and %s is not allowed", subjectKindFlags)
		}

		if !use(o.SubjectName) {
			return fmt.Errorf("--name is required when using --kind")
		}

	} else {
		if kindFlagCount == 0 {
			return fmt.Errorf("You must specify Subject's Kind. Use --kind or %s", subjectKindFlags)
		}

		if kindFlagCount > 1 {
			return fmt.Errorf("Cannot use multiple kind flag in %s", subjectKindFlags)
		}

		if use(o.SubjectAPIGroup) || use(o.SubjectName) {
			return fmt.Errorf("--api-group or --name is not allowed when using kind flags in %s", subjectKindFlags)
		}

	}
	return nil
}

func (o *AttachDetachOptions) Complete(cmd *cobra.Command, args []string) error {
	o.PSPName = args[0]
	kind, name, _ := getValuesFromKindFlags(o)
	if kind != "" {
		o.SubjectKind = kind
		o.SubjectName = name
	}
	return nil
}

func (o *AttachDetachOptions) GenerateSubject(kubeconfigPath *string) (*rbacv1.Subject, error) {
	sub := &rbacv1.Subject{}
	sub.Kind = o.SubjectKind
	sub.Name = o.SubjectName

	if use(o.ServiceAccount) {
		if o.SubjectNamespace == "" {
			namespace, err := client.GetDefaultNamespace(kubeconfigPath)
			if err != nil {
				return nil, err
			}
			sub.Namespace = namespace
		} else {
			sub.Namespace = o.SubjectNamespace
		}
	}

	if use(o.Group) || use(o.User) {
		sub.APIGroup = rbac.APIGroup
	}

	return sub, nil
}

func use(v string) bool {
	return v != ""
}

func getValuesFromKindFlags(o *AttachDetachOptions) (kind, name string, kindFlagCount int) {
	option := reflect.Indirect(reflect.ValueOf(o))
	for i := 0; i < option.Type().NumField(); i++ {
		fieldName := option.Type().Field(i).Name

		for _, subKind := range subjectKindList {
			if subKind == fieldName {
				subName := option.Field(i).String()
				if subName != "" {
					kind = subKind
					name = subName
					kindFlagCount++
				}
			}
		}
	}
	return kind, name, kindFlagCount
}
