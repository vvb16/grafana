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

func NewAccessControlDashboardGuardian(ctx context.Context, dashboardId int64, user *models.SignedInUser, store *sqlstore.SQLStore, ac accesscontrol.AccessControl) *AccessControlDashboardGuardian {
	return &AccessControlDashboardGuardian{
		ctx:         ctx,
		dashboardID: dashboardId,
		user:        user,
		store:       store,
		ac:          ac,
	}
}

type AccessControlDashboardGuardian struct {
	ctx         context.Context
	dashboardID int64

	dashboard *models.Dashboard
	user      *models.SignedInUser

	store *sqlstore.SQLStore
	ac    accesscontrol.AccessControl
}

func (a *AccessControlDashboardGuardian) CanSave() (bool, error) {
	if err := a.loadDashboard(); err != nil {
		return false, err
	}

	return a.ac.Evaluate(a.ctx, a.user, accesscontrol.EvalAny(
		accesscontrol.EvalPermission(accesscontrol.ActionDashboardsWrite, dashboardScope(a.dashboard.Id)),
		accesscontrol.EvalPermission(accesscontrol.ActionFoldersWrite, folderScope(a.dashboard.FolderId)),
	))
}

func (a AccessControlDashboardGuardian) CanEdit() (bool, error) {
	return a.CanSave()
}

func (a AccessControlDashboardGuardian) CanView() (bool, error) {
	if err := a.loadDashboard(); err != nil {
		return false, err
	}

	return a.ac.Evaluate(a.ctx, a.user, accesscontrol.EvalAny(
		accesscontrol.EvalPermission(accesscontrol.ActionDashboardsRead, dashboardScope(a.dashboard.Id)),
		accesscontrol.EvalPermission(accesscontrol.ActionFoldersRead, folderScope(a.dashboard.FolderId)),
	))
}

func (a AccessControlDashboardGuardian) CanAdmin() (bool, error) {
	if err := a.loadDashboard(); err != nil {
		return false, err
	}

	return a.ac.Evaluate(a.ctx, a.user, accesscontrol.EvalAny(
		accesscontrol.EvalAll(
			accesscontrol.EvalPermission(accesscontrol.ActionDashboardsPermissionsRead, dashboardScope(a.dashboard.Id)),
			accesscontrol.EvalPermission(accesscontrol.ActionDashboardsPermissionsWrite, dashboardScope(a.dashboard.Id)),
		),
		accesscontrol.EvalAll(
			accesscontrol.EvalPermission(accesscontrol.ActionFoldersPermissionsRead, folderScope(a.dashboard.FolderId)),
			accesscontrol.EvalPermission(accesscontrol.ActionFoldersPermissionsWrite, folderScope(a.dashboard.FolderId)),
		),
	))
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

func folderScope(folderID int64) string {
	return accesscontrol.Scope("folders", "id", strconv.FormatInt(folderID, 10))
}
