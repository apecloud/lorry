/*
Copyright (C) 2022-2024 ApeCloud Co., Ltd

This file is part of KubeBlocks project

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package ctl

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/apecloud/dbctl/operations"
)

type Options interface {
	Init() error
	Validate() error
	Run() error
	GetAction() string
}

// OptionsBase represents the cmd configuration parameters.
type OptionsBase struct {
	Action  string
	Request *operations.OpsRequest
	operations.Operation
}

func (options *OptionsBase) Init() error {
	ops := operations.Operations()

	operation, ok := ops[strings.ToLower(options.Action)]
	if !ok {
		return errors.New(options.Action + " operation not found")
	}
	err := operation.Init(context.Background())
	if err != nil {
		return errors.Wrap(err, "getrole init failed")
	}
	options.Operation = operation
	return nil
}

func (options *OptionsBase) Validate() error {
	return options.PreCheck(context.Background(), options.Request)
}

func (options *OptionsBase) Run() error {
	return nil
}

func (options *OptionsBase) GetAction() string {
	return options.Action
}

func CmdRunner(options Options) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := options.Init()
		if err != nil {
			fmt.Printf("%s init failed: %s\n", options.GetAction(), err.Error())
			os.Exit(1)
		}

		err = options.Validate()
		if err != nil {
			fmt.Printf("%s validate failed: %s\n", options.GetAction(), err.Error())
			os.Exit(1)
		}

		err = options.Run()
		if err != nil {
			fmt.Printf("%s executing failed: %s\n", options.GetAction(), err.Error())
			os.Exit(1)
		}
	}
}
