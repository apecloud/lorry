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

	"github.com/apecloud/lorry/operations/replica"
)

type GetRoleOptions struct {
	Options
}

func (options *GetRoleOptions) Run() error {
	getRole, ok := getRoleOptions.Operation.(*replica.GetRole)
	if !ok {
		return errors.New("getrole operation not found")

	}

	cluster, err := getRole.DCSStore.GetCluster()
	if err != nil {
		return errors.Wrap(err, "get cluster failed")
	}

	role, err := getRole.DBManager.GetReplicaRole(context.Background(), cluster)
	if err != nil {
		return errors.Wrap(err, "executing getrole failed")
	}
	fmt.Print(role)
	return nil
}

var getRoleOptions = &GetRoleOptions{
	Options: Options{
		Action: "getrole",
	},
}

var GetRoleCmd = &cobra.Command{
	Use:   "getrole",
	Short: "get role of the replica.",
	Example: `
lorry getrole 
  `,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		err := getRoleOptions.Init()
		if err != nil {
			fmt.Printf("getrole init failed: %v\n", err.Error())
			os.Exit(1)
		}

		err = getRoleOptions.Validate()
		if err != nil {
			fmt.Printf("getrole validate failed: %v\n", err.Error())
			os.Exit(1)
		}

		err = getRoleOptions.Run()
		if err != nil {
			fmt.Printf("executing getrole failed: %s\n", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	GetRoleCmd.Flags().BoolP("help", "h", false, "Print this help message")

	RootCmd.AddCommand(GetRoleCmd)
}
