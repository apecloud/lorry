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

	"github.com/apecloud/lorry/operations"
	"github.com/apecloud/lorry/operations/replica"
	"github.com/apecloud/lorry/util"
)

type SwitchOptions struct {
	OptionsBase
	primary   string
	candidate string
	force     bool
}

func (options *SwitchOptions) Validate() error {
	params := map[string]interface{}{}
	if options.primary != "" {
		params["primary"] = options.primary
	}
	if options.candidate != "" {
		params["candidate"] = options.candidate
	}
	req := &operations.OpsRequest{
		Parameters: params,
	}
	options.Request = req
	return options.Operation.PreCheck(context.Background(), req)
}

func (options *SwitchOptions) Run() error {
	switchover, ok := options.Operation.(*replica.Switchover)
	if !ok {
		return errors.Errorf("%s operation not found", options.Action)
	}

	_, err := switchover.Do(context.Background(), options.Request)
	if err != nil {
		return errors.Wrap(err, "executing switchover failed")
	}
	return nil
}

var switchoverOptions = &SwitchOptions{
	OptionsBase: OptionsBase{
		Action: string(util.SwitchoverOperation),
	},
}

var SwitchCmd = &cobra.Command{
	Use:   "switchover",
	Short: "execute a switchover request.",
	Example: `
lorry switchover --primary xxx --candidate xxx
  `,
	Args: cobra.MinimumNArgs(0),
	Run:  CmdRunner(switchoverOptions),
}

func init() {
	SwitchCmd.Flags().StringVarP(&switchoverOptions.primary, "primary", "p", "", "The primary pod name")
	SwitchCmd.Flags().StringVarP(&switchoverOptions.candidate, "candidate", "c", "", "The candidate pod name")
	SwitchCmd.Flags().BoolVarP(&switchoverOptions.force, "force", "f", false, "force to swithover if failed")
	SwitchCmd.Flags().BoolP("help", "h", false, "Print this help message")

	RootCmd.AddCommand(SwitchCmd)
}
