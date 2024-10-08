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

package replica

import (
	"context"
	"strings"

	"github.com/go-logr/logr"

	"github.com/apecloud/dbctl/operations"
	"github.com/apecloud/dbctl/util"
)

type dataDump struct {
	operations.Base
	logger  logr.Logger
	Command []string
}

func init() {
	err := operations.Register(strings.ToLower(string(util.DataDumpOperation)), &dataDump{})
	if err != nil {
		panic(err.Error())
	}
}

func (s *dataDump) Do(ctx context.Context, req *operations.OpsRequest) (*operations.OpsResponse, error) {
	return nil, doCommonAction(ctx, s.logger, "dataDump", s.Command)
}

func doCommonAction(ctx context.Context, logger logr.Logger, action string, commands []string) error {
	envs, err := util.GetGlobalSharedEnvs()
	if err != nil {
		return err
	}
	output, err := util.ExecCommand(ctx, commands, envs)
	if output != "" {
		logger.Info(action, "output", output)
	}
	return err
}
