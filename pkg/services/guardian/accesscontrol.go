package guardian

import (
	"context"
	"strconv"

	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/sqlstore"
	"github.com/grafana/grafana/pkg/setting"
)

var _ DashboardGuardian = new(AccessControlDashboardGuardian)

func NewAccessControlDashboardGuardian() *AccessControlDashboardGuardian {
	return &AccessControlDashboardGuardian{}
}

type AccessControlDashboardGuardian struct {
	ctx         context.Context
	dashboardID int64
	dashboard   *models.Dashboard
	user        *models.SignedInUser

	store *sqlstore.SQLStore
	ac    accesscontrol.AccessControl
}

func (a *AccessControlDashboardGuardian) CanSave() (bool, error) {
	if err := a.loadDashboard(); err != nil {
		return false, err
	}

	evaluators := []accesscontrol.Evaluator{
		accesscontrol.EvalPermission(accesscontrol.ActionDashboardsWrite, dashboardScope(a.dashboard.Id)),
	}

	if a.dashboard.FolderId != 0 {
		evaluators = append(evaluators, accesscontrol.EvalPermission(accesscontrol.ActionDashboardsWrite, dashboardScope(a.dashboard.FolderId)))
	}

	return a.ac.Evaluate(a.ctx, a.user, accesscontrol.EvalAny(evaluators...))
}

func (a AccessControlDashboardGuardian) CanEdit() (bool, error) {
	return a.CanSave()
}

func (a AccessControlDashboardGuardian) CanView() (bool, error) {
	if err := a.loadDashboard(); err != nil {
		return false, err
	}

	evaluators := []accesscontrol.Evaluator{
		accesscontrol.EvalPermission(accesscontrol.ActionDashboardsRead, dashboardScope(a.dashboard.Id)),
	}

	if a.dashboard.FolderId != 0 {
		evaluators = append(evaluators, accesscontrol.EvalPermission(accesscontrol.ActionDashboardsRead, dashboardScope(a.dashboard.FolderId)))
	}

	return a.ac.Evaluate(a.ctx, a.user, accesscontrol.EvalAny(evaluators...))
}

func (a AccessControlDashboardGuardian) CanAdmin() (bool, error) {
	if err := a.loadDashboard(); err != nil {
		return false, err
	}

	evaluators := []accesscontrol.Evaluator{
		accesscontrol.EvalAll(
			accesscontrol.EvalPermission(accesscontrol.ActionDashboardsPermissionRead, dashboardScope(a.dashboard.Id)),
			accesscontrol.EvalPermission(accesscontrol.ActionDashboardsPermissionWrite, dashboardScope(a.dashboard.Id)),
		),
	}

	if a.dashboard.FolderId != 0 {
		evaluators = append(evaluators,
			accesscontrol.EvalAll(
				accesscontrol.EvalPermission(accesscontrol.ActionDashboardsPermissionRead, dashboardScope(a.dashboard.FolderId)),
				accesscontrol.EvalPermission(accesscontrol.ActionDashboardsPermissionWrite, dashboardScope(a.dashboard.FolderId)),
			),
		)
	}

	return a.ac.Evaluate(a.ctx, a.user, accesscontrol.EvalAny(evaluators...))
}

func (a AccessControlDashboardGuardian) CheckPermissionBeforeUpdate(permission models.PermissionType, updatePermissions []*models.DashboardAcl) (bool, error) {
	panic("implement me")
}

func (a AccessControlDashboardGuardian) GetAcl() ([]*models.DashboardAclInfoDTO, error) {
	panic("implement me")
}

func (a AccessControlDashboardGuardian) GetACLWithoutDuplicates() ([]*models.DashboardAclInfoDTO, error) {
	panic("implement me")
}

func (a AccessControlDashboardGuardian) GetHiddenACL(cfg *setting.Cfg) ([]*models.DashboardAcl, error) {
	panic("implement me")
}

func (a *AccessControlDashboardGuardian) loadDashboard() error {
	if a.dashboard == nil {
		dashboard, err := a.store.GetDashboard(a.dashboardID, a.user.OrgId, "", "")
		if err != nil {
			return err
		}
		a.dashboard = dashboard
	}
	return nil
}

func dashboardScope(dashboardID int64) string {
	return accesscontrol.Scope("dashboards", "id", strconv.FormatInt(dashboardID, 10))
}
