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
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/apecloud/dbctl/dcs"
	"github.com/apecloud/dbctl/engines"
	"github.com/apecloud/dbctl/engines/register"
	"github.com/apecloud/dbctl/operations"
	"github.com/apecloud/dbctl/util"
)

type CheckStatus struct {
	operations.Base
	LeaderFailedCount          int
	FailureThreshold           int
	dcsStore                   dcs.DCS
	dbManager                  engines.DBManager
	checkFailedCount           int
	failedEventReportFrequency int
	logger                     logr.Logger
}

type FailoverManager interface {
	Failover(ctx context.Context, cluster *dcs.Cluster, candidate string) error
}

var checkstatus operations.Operation = &CheckStatus{}

func init() {
	err := operations.Register(strings.ToLower(string(util.HealthyCheckOperation)), checkstatus)
	if err != nil {
		panic(err.Error())
	}
}

func (s *CheckStatus) Init(ctx context.Context) error {
	s.dcsStore = dcs.GetStore()
	if s.dcsStore == nil {
		return errors.New("dcs store init failed")
	}

	s.failedEventReportFrequency = viper.GetInt("KB_FAILED_EVENT_REPORT_FREQUENCY")
	if s.failedEventReportFrequency < 300 {
		s.failedEventReportFrequency = 300
	} else if s.failedEventReportFrequency > 3600 {
		s.failedEventReportFrequency = 3600
	}

	s.FailureThreshold = 3
	s.logger = ctrl.Log.WithName("checkstatus")
	dbManager, err := register.GetDBManager()
	if err != nil {
		return errors.Wrap(err, "get manager failed")
	}
	s.dbManager = dbManager

	return nil
}

func (s *CheckStatus) IsReadonly(ctx context.Context) bool {
	return true
}

func (s *CheckStatus) Do(ctx context.Context, req *operations.OpsRequest) (*operations.OpsResponse, error) {
	resp := &operations.OpsResponse{
		Data: map[string]any{},
	}
	resp.Data["operation"] = util.HealthyCheckOperation

	k8sStore := s.dcsStore.(*dcs.KubernetesStore)
	cluster := k8sStore.GetClusterFromCache()
	err := s.dbManager.CurrentMemberHealthyCheck(ctx, cluster)
	if err != nil {
		return s.handlerError(ctx, err)
	}

	isLeader, err := s.dbManager.IsLeader(ctx, cluster)
	if err != nil {
		return s.handlerError(ctx, err)
	}

	if isLeader {
		s.LeaderFailedCount = 0
		s.checkFailedCount = 0
		resp.Data["event"] = util.OperationSuccess
		return resp, nil
	}
	err = s.dbManager.LeaderHealthyCheck(ctx, cluster)
	if err != nil {
		s.LeaderFailedCount++
		if s.LeaderFailedCount > s.FailureThreshold {
			err = s.failover(ctx, cluster)
			if err != nil {
				return s.handlerError(ctx, err)
			}
		}
		return s.handlerError(ctx, err)
	}
	s.LeaderFailedCount = 0
	s.checkFailedCount = 0
	resp.Data["event"] = util.OperationSuccess
	return resp, nil
}

func (s *CheckStatus) failover(ctx context.Context, cluster *dcs.Cluster) error {
	failoverManger, ok := s.dbManager.(FailoverManager)
	if !ok {
		return errors.New("failover manager not found")
	}
	err := failoverManger.Failover(ctx, cluster, s.dbManager.GetCurrentMemberName())
	if err != nil {
		return errors.Wrap(err, "failover failed")
	}
	return nil
}

func (s *CheckStatus) handlerError(ctx context.Context, err error) (*operations.OpsResponse, error) {
	resp := &operations.OpsResponse{
		Data: map[string]any{},
	}
	message := err.Error()
	s.logger.Info("healthy checks failed", "error", message)
	resp.Data["event"] = util.OperationFailed
	resp.Data["message"] = message
	if s.checkFailedCount%s.failedEventReportFrequency == 0 {
		s.logger.Info("healthy checks failed continuously", "times", s.checkFailedCount)
		_ = util.SentEventForProbe(ctx, resp.Data)
	}
	s.checkFailedCount++
	return resp, err
}
