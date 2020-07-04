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

	"github.com/spf13/cobra"
)

type CleanOptions struct {
	PSPName string
}

func (o *CleanOptions) PreRunE(cmd *cobra.Command, args []string) error {
	if err := o.Validate(cmd, args); err != nil {
		return err
	}
	if err := o.Complete(cmd, args); err != nil {
		return err
	}
	return nil
}

func (o *CleanOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Args is invalid. Required: `PSP-NAME`")
	}
	return nil
}

func (o *CleanOptions) Complete(cmd *cobra.Command, args []string) error {
	o.PSPName = args[0]
	return nil
}
