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

package user

import (
	"context"
	"strings"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/apecloud/dbctl/engines"
	"github.com/apecloud/dbctl/engines/models"
	"github.com/apecloud/dbctl/engines/register"
	"github.com/apecloud/dbctl/operations"
	"github.com/apecloud/dbctl/util"
)

type CreateUser struct {
	operations.Base
	DBManager engines.DBManager
	logger    logr.Logger
}

var createUser operations.Operation = &CreateUser{}

func init() {
	err := operations.Register(strings.ToLower(string(util.CreateUserOp)), createUser)
	if err != nil {
		panic(err.Error())
	}
}

func (s *CreateUser) Init(ctx context.Context) error {
	dbManager, err := register.GetDBManager()
	if err != nil {
		return errors.Wrap(err, "get manager failed")
	}
	s.DBManager = dbManager
	s.logger = ctrl.Log.WithName("CreateUser")
	return nil
}

func (s *CreateUser) IsReadonly(ctx context.Context) bool {
	return false
}

func (s *CreateUser) PreCheck(ctx context.Context, req *operations.OpsRequest) error {
	userInfo, err := UserInfoParser(req)
	if err != nil {
		return err
	}

	return userInfo.UserNameAndPasswdValidator()
}

func (s *CreateUser) Do(ctx context.Context, req *operations.OpsRequest) (*operations.OpsResponse, error) {
	userInfo, _ := UserInfoParser(req)
	resp := operations.NewOpsResponse(util.CreateUserOp)

	user, err := s.DBManager.DescribeUser(ctx, userInfo.UserName)
	if err == nil && user != nil {
		return resp.WithSuccess("account already exists")
	}

	// for compatibility with old addons that specify accoutprovision action but not work actually.
	err = s.DBManager.CreateUser(ctx, userInfo.UserName, userInfo.Password)
	if err != nil {
		err = errors.Cause(err)
		s.logger.Info("executing CreateUser error", "error", err.Error())
		return resp, err
	}

	if userInfo.RoleName != "" {
		err := s.DBManager.GrantUserRole(ctx, userInfo.UserName, userInfo.RoleName)
		if err != nil && err != models.ErrNotImplemented {
			s.logger.Info("executing grantRole error", "error", err.Error())
			return resp, err
		}
	}

	return resp.WithSuccess("")
}
