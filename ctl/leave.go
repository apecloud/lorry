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

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/apecloud/dbctl/operations/replica"
	"github.com/apecloud/dbctl/util"
)

type LeaveOptions struct {
	OptionsBase
}

func (options *LeaveOptions) Run() error {
	leave, ok := options.Operation.(*replica.Leave)
	if !ok {
		return errors.Errorf("%s operation not found", options.Action)
	}

	_, err := leave.Do(context.Background(), nil)
	if err != nil {
		return errors.Wrap(err, "executing leave failed")
	}
	return nil
}

var leaveOptions = &LeaveOptions{
	OptionsBase: OptionsBase{
		Action: string(util.LeaveMemberOperation),
	},
}

var LeaveCmd = &cobra.Command{
	Use:   "leavemember",
	Short: "execute a leave member request.",
	Example: `
dbctl leavemember
  `,
	Args: cobra.MinimumNArgs(0),
	Run:  CmdRunner(leaveOptions),
}

func init() {
	LeaveCmd.Flags().BoolP("help", "h", false, "Print this help message")

	DatabaseCmd.AddCommand(LeaveCmd)
}
