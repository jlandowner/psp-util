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

package printers

import (
	"io"
	"reflect"

	"github.com/liggitt/tabwriter"
)

var ListHeader = []string{"PSP", "ClusterRole", "ClusterRoleBinding", "NS/Role", "NS/RoleBinding", "Managed"}

type ListPrinterLine struct {
	PSP                string
	ClusterRole        string
	ClusterRoleBinding string
	Role               string
	RoleBinding        string
	PSPUtilManaged     string
}

type ListPrinterOptions struct {
	PSP                bool
	ClusterRole        bool
	ClusterRoleBinding bool
	Role               bool
	RoleBinding        bool
	PSPUtilManaged     bool
}

type ListPrinter struct {
	options ListPrinterOptions
	w       *tabwriter.Writer
}

func NewListPrinter(output io.Writer, options ListPrinterOptions) *ListPrinter {
	l := &ListPrinter{
		options: options,
		w:       GetNewTabWriter(output),
	}
	return l
}

func (l *ListPrinter) PrintHeader() {
	printHeader := []string{}
	lval := reflect.ValueOf(l.options)

	// loop ListPrinter struct properties. Skip the last property "w"
	for i := 0; i < len(ListHeader); i++ {
		if reflect.Indirect(lval).Field(i).Bool() {
			printHeader = append(printHeader, ListHeader[i])
		}
	}
	PrintLine(l.w, printHeader)
}

func (l *ListPrinter) PrintLine(line ListPrinterLine) {
	printLine := []string{}
	lineval := reflect.ValueOf(line)
	lval := reflect.ValueOf(l.options)

	// loop ListPrinter struct properties. Skip the last property "w"
	for i := 0; i < len(ListHeader); i++ {
		if reflect.Indirect(lval).Field(i).Bool() {
			printLine = append(printLine, reflect.Indirect(lineval).Field(i).String())
		}
	}
	PrintLine(l.w, printLine)
}

func (l *ListPrinter) Flush() {
	l.w.Flush()
}
