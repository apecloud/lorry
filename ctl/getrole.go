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

	"github.com/spf13/cobra"

	"github.com/apecloud/lorry/operations"
	"github.com/apecloud/lorry/operations/replica"
)

type GetRoleOptions struct {
}

var GetRoleCmd = &cobra.Command{
	Use:   "getrole",
	Short: "get role of the replica.",
	Example: `
lorry getrole 
  `,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		options := operations.Operations()

		option, ok := options["getrole"]
		if !ok {
			fmt.Printf("getrole operation not found")
			os.Exit(1)
		}
		err := option.Init(context.Background())
		if err != nil {
			fmt.Printf("getrole init failed: %v\n", err)
			os.Exit(1)
		}

		err = option.PreCheck(context.Background(), nil)
		if err != nil {
			fmt.Printf("getrole precheck failed: %v\n", err)
			os.Exit(1)
		}

		getRole, ok := option.(*replica.GetRole)
		if !ok {
			fmt.Printf("getrole operation not found")
			os.Exit(1)
		}

		cluster, err := getRole.DCSStore.GetCluster()
		if err != nil {
			fmt.Printf("executing getrole failed: %s\n", err.Error())
			os.Exit(1)
		}

		role, err := getRole.DBManager.GetReplicaRole(context.Background(), cluster)
		if err != nil {
			fmt.Printf("executing getrole failed: %s\n", err.Error())
			os.Exit(1)
		}

		fmt.Print(role)
	},
}

func init() {
	GetRoleCmd.Flags().BoolP("help", "h", false, "Print this help message")

	RootCmd.AddCommand(GetRoleCmd)
}
