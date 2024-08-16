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

package component

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/apecloud/dbctl/engines/models"
	"github.com/apecloud/dbctl/engines/register"
	"github.com/apecloud/dbctl/operations"
	"github.com/apecloud/dbctl/util"
)

type PostProvision struct {
	operations.Base
	Timeout time.Duration
	Command []string
}

type PostProvisionManager interface {
	PostProvision(ctx context.Context, componentNames, podNames, podIPs, podHostNames, podHostIPs string) error
}

var postProvision operations.Operation = &PostProvision{}

func init() {
	err := operations.Register(strings.ToLower(string(util.PostProvisionOperation)), postProvision)
	if err != nil {
		panic(err.Error())
	}
}

func (s *PostProvision) Do(ctx context.Context, req *operations.OpsRequest) (*operations.OpsResponse, error) {
	componentNames := req.GetString("componentNames")
	podNames := req.GetString("podNames")
	podIPs := req.GetString("podIPs")
	podHostNames := req.GetString("podHostNames")
	podHostIPs := req.GetString("podHostIPs")
	manager, err := register.GetDBManager()
	if err != nil {
		return nil, errors.Wrap(err, "get manager failed")
	}

	ppManager, ok := manager.(PostProvisionManager)
	if !ok {
		return nil, models.ErrNotImplemented
	}
	err = ppManager.PostProvision(ctx, componentNames, podNames, podIPs, podHostNames, podHostIPs)
	return nil, err
}
