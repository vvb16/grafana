package guardian

import (
	"context"
	"strconv"

	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/sqlstore"
	"github.com/grafana/grafana/pkg/setting"
)

// TODO: temp location
const (
	ActionRead            = "dashboards:read"
	ActionWrite           = "dashboards:write"
	ActionDelete          = "dashboards:delete"
	ActionPermissionRead  = "dashboards.permissions:read"
	ActionPermissionWrite = "dashboards.permissions:write"
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
		accesscontrol.EvalPermission(ActionWrite, dashboardScope(a.dashboard.Id)),
	}

	if a.dashboard.FolderId != 0 {
		evaluators = append(evaluators, accesscontrol.EvalPermission(ActionWrite, dashboardScope(a.dashboard.FolderId)))
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
		accesscontrol.EvalPermission(ActionRead, dashboardScope(a.dashboard.Id)),
	}

	if a.dashboard.FolderId != 0 {
		evaluators = append(evaluators, accesscontrol.EvalPermission(ActionRead, dashboardScope(a.dashboard.FolderId)))
	}

	return a.ac.Evaluate(a.ctx, a.user, accesscontrol.EvalAny(evaluators...))
}

func (a AccessControlDashboardGuardian) CanAdmin() (bool, error) {
	if err := a.loadDashboard(); err != nil {
		return false, err
	}

	evaluators := []accesscontrol.Evaluator{
		accesscontrol.EvalAll(
			accesscontrol.EvalPermission(ActionPermissionRead, dashboardScope(a.dashboard.Id)),
			accesscontrol.EvalPermission(ActionPermissionWrite, dashboardScope(a.dashboard.Id)),
		),
	}

	if a.dashboard.FolderId != 0 {
		evaluators = append(evaluators,
			accesscontrol.EvalAll(
				accesscontrol.EvalPermission(ActionPermissionRead, dashboardScope(a.dashboard.FolderId)),
				accesscontrol.EvalPermission(ActionPermissionRead, dashboardScope(a.dashboard.FolderId)),
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
