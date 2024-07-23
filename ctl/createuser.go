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

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/apecloud/lorry/operations"
	"github.com/apecloud/lorry/operations/user"
	"github.com/apecloud/lorry/util"
)

type CreateUserOptions struct {
	Options
	userName string
	password string
	roleName string
}

func (options *CreateUserOptions) Validate() error {
	parameters := map[string]any{
		"userName": options.userName,
		"password": options.password,
	}
	if options.roleName != "" {
		parameters["roleName"] = options.roleName
	}

	req := &operations.OpsRequest{
		Parameters: parameters,
	}
	options.Request = req
	return options.Operation.PreCheck(context.Background(), req)
}

func (options *CreateUserOptions) Run() error {
	createUser, ok := createUserOptions.Operation.(*user.CreateUser)
	if !ok {
		return errors.New("createUser operation not found")
	}

	_, err := createUser.Do(context.Background(), options.Request)
	if err != nil {
		return errors.Wrap(err, "executing createUser failed")
	}
	return nil
}

var createUserOptions = &CreateUserOptions{
	Options: Options{
		Action: string(util.CreateUserOp),
	},
}

var CreateUserCmd = &cobra.Command{
	Use:   "createuser",
	Short: "create user.",
	Example: `
lorryctl  createuser --username xxx --password xxx 
  `,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		err := createUserOptions.Init()
		if err != nil {
			fmt.Printf("%s init failed: %s\n", createUserOptions.Action, err.Error())
			os.Exit(1)
		}

		err = createUserOptions.Validate()
		if err != nil {
			fmt.Printf("%s validate failed: %s\n", createUserOptions.Action, err.Error())
			os.Exit(1)
		}

		err = createUserOptions.Run()
		if err != nil {
			fmt.Printf("%s executing failed: %s\n", createUserOptions.Action, err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	CreateUserCmd.Flags().StringVarP(&createUserOptions.userName, "username", "", "", "The name of user to create")
	CreateUserCmd.Flags().StringVarP(&createUserOptions.password, "password", "", "", "The password of user to create")
	CreateUserCmd.Flags().StringVarP(&createUserOptions.roleName, "rolename", "", "", "The role of user to create")
	CreateUserCmd.Flags().BoolP("help", "h", false, "Print this help message")

	RootCmd.AddCommand(CreateUserCmd)
}
