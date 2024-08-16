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
	"github.com/apecloud/dbctl/constant"
	"github.com/apecloud/dbctl/dcs"
	"github.com/apecloud/dbctl/engines/register"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var DatabaseCmd = &cobra.Command{
	Use:     "mongodb",
	Aliases: []string{"mysql", "postgresql"},
	Short:   "specify database.",
	Example: `
dbctl mongodb createuser --username root --password password
  `,
	Args: cobra.MinimumNArgs(0),
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		calledAs := cmd.CalledAs()
		viper.SetDefault(constant.KBEnvEngineType, calledAs)
		// Initialize DCS (Distributed Control System)
		err := dcs.InitStore()
		if err != nil {
			return errors.Wrap(err, "DCS initialize failed")
		}

		// Initialize DB Manager
		err = register.InitDBManager(configDir)
		if err != nil {
			return errors.Wrap(err, "DB manager initialize failed")
		}
		return nil
	},

	Run: func(cmd *cobra.Command, _ []string) {
		_ = cmd.Help()
	},
}

func init() {
	RootCmd.AddCommand(DatabaseCmd)
}
