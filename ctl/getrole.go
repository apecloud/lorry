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

	"github.com/spf13/cobra"

	"github.com/apecloud/kubeblocks/pkg/lorry/client"
)

type GetRoleOptions struct {
	lorryAddr string
	userName  string
	password  string
}

var getRoleOptions = &GetRoleOptions{}

var GetRoleCmd = &cobra.Command{
	Use:   "createuser",
	Short: "create user.",
	Example: `
lorryctl  createuser --username xxx --password xxx 
  `,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		lorryClient, err := client.NewHTTPClientWithURL(getRoleOptions.lorryAddr)
		if err != nil {
			fmt.Printf("new lorry http client failed: %v\n", err)
			return
		}

		err = lorryClient.CreateUser(context.TODO(), getRoleOptions.userName, getRoleOptions.password, "")
		if err != nil {
			fmt.Printf("create user failed: %v\n", err)
			return
		}
		fmt.Printf("create user success")
	},
}

func init() {
	GetRoleCmd.Flags().StringVarP(&getRoleOptions.userName, "username", "", "", "The name of user to create")
	GetRoleCmd.Flags().StringVarP(&getRoleOptions.password, "password", "", "", "The password of user to create")
	GetRoleCmd.Flags().StringVarP(&getRoleOptions.lorryAddr, "lorry-addr", "", "http://localhost:3501/v1.0/", "The addr of lorry to request")
	GetRoleCmd.Flags().BoolP("help", "h", false, "Print this help message")

	RootCmd.AddCommand(GetRoleCmd)
}
