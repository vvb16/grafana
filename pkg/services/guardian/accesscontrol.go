package guardian

import (
	"context"
	"strconv"

	"github.com/grafana/grafana/pkg/services/accesscontrol/resourcepermissions"

	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/sqlstore"
	"github.com/grafana/grafana/pkg/setting"
)

var permissionMap = map[string]models.PermissionType{
	"View":  models.PERMISSION_VIEW,
	"Edit":  models.PERMISSION_EDIT,
	"Admin": models.PERMISSION_ADMIN,
}

var _ DashboardGuardian = new(AccessControlDashboardGuardian)

func NewAccessControlDashboardGuardian(
	ctx context.Context, dashboardId int64, user *models.SignedInUser,
	store *sqlstore.SQLStore, ac accesscontrol.AccessControl, permissionServices *resourcepermissions.Services,
) *AccessControlDashboardGuardian {
	return &AccessControlDashboardGuardian{
		ctx:                ctx,
		dashboardID:        dashboardId,
		user:               user,
		store:              store,
		ac:                 ac,
		permissionServices: permissionServices,
	}
}

type AccessControlDashboardGuardian struct {
	ctx                context.Context
	dashboardID        int64
	dashboard          *models.Dashboard
	user               *models.SignedInUser
	store              *sqlstore.SQLStore
	ac                 accesscontrol.AccessControl
	permissionServices *resourcepermissions.Services
}

func (a *AccessControlDashboardGuardian) CanSave() (bool, error) {
	// check permissions to create new dashboard / folder
	if a.dashboardID == 0 {
		return a.ac.Evaluate(a.ctx, a.user, accesscontrol.EvalAny(
			accesscontrol.EvalPermission(accesscontrol.ActionFoldersCreate),
			accesscontrol.EvalPermission(accesscontrol.ActionDashboardsCreate),
		))
	}

	return a.CanEdit()
}

func (a *AccessControlDashboardGuardian) CanEdit() (bool, error) {
	if err := a.loadDashboard(); err != nil {
		return false, err
	}

	if a.dashboard.IsFolder {
		return a.ac.Evaluate(a.ctx, a.user, accesscontrol.EvalPermission(accesscontrol.ActionFoldersWrite, folderScope(a.dashboardID)))
	}

	return a.ac.Evaluate(a.ctx, a.user, accesscontrol.EvalAny(
		accesscontrol.EvalPermission(accesscontrol.ActionDashboardsWrite, dashboardScope(a.dashboard.Id)),
		accesscontrol.EvalPermission(accesscontrol.ActionFoldersWrite, folderScope(a.dashboard.FolderId)),
	))
}

func (a *AccessControlDashboardGuardian) CanView() (bool, error) {
	if err := a.loadDashboard(); err != nil {
		return false, err
	}

	if a.dashboard.IsFolder {
		return a.ac.Evaluate(a.ctx, a.user, accesscontrol.EvalPermission(accesscontrol.ActionFoldersRead, folderScope(a.dashboardID)))
	}

	return a.ac.Evaluate(a.ctx, a.user, accesscontrol.EvalAny(
		accesscontrol.EvalPermission(accesscontrol.ActionDashboardsRead, dashboardScope(a.dashboard.Id)),
		accesscontrol.EvalPermission(accesscontrol.ActionFoldersRead, folderScope(a.dashboard.FolderId)),
	))
}

func (a *AccessControlDashboardGuardian) CanAdmin() (bool, error) {
	if err := a.loadDashboard(); err != nil {
		return false, err
	}

	if a.dashboard.IsFolder {
		return a.ac.Evaluate(a.ctx, a.user, accesscontrol.EvalAll(
			accesscontrol.EvalPermission(accesscontrol.ActionFoldersPermissionsRead, folderScope(a.dashboard.Id)),
			accesscontrol.EvalPermission(accesscontrol.ActionFoldersPermissionsWrite, folderScope(a.dashboard.Id)),
		))
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

func (a *AccessControlDashboardGuardian) CheckPermissionBeforeUpdate(permission models.PermissionType, updatePermissions []*models.DashboardAcl) (bool, error) {
	// not used with access control
	return false, nil
}

func (a *AccessControlDashboardGuardian) GetAcl() ([]*models.DashboardAclInfoDTO, error) {
	// translate access control permissions to dashboard acl info
	if err := a.loadDashboard(); err != nil {
		return nil, err
	}

	svc := a.permissionServices.GetDashboardService()
	if a.dashboard.IsFolder {
		svc = a.permissionServices.GetFolderService()
	}

	permissions, err := svc.GetPermissions(a.ctx, a.dashboard.OrgId, strconv.FormatInt(a.dashboard.Id, 10))
	if err != nil {
		return nil, err
	}

	acl := make([]*models.DashboardAclInfoDTO, 0, len(permissions))
	for _, p := range permissions {
		var role *models.RoleType
		if p.BuiltInRole != "" {
			tmp := models.RoleType(p.BuiltInRole)
			role = &tmp
		}

		acl = append(acl, &models.DashboardAclInfoDTO{
			OrgId:          a.dashboard.OrgId,
			DashboardId:    a.dashboard.Id,
			FolderId:       a.dashboard.FolderId,
			Created:        p.Created,
			Updated:        p.Updated,
			UserId:         p.UserId,
			UserLogin:      p.UserLogin,
			UserEmail:      p.UserEmail,
			TeamId:         p.TeamId,
			TeamEmail:      p.TeamEmail,
			Team:           p.Team,
			Role:           role,
			Permission:     permissionMap[svc.MapActions(p)],
			PermissionName: permissionMap[svc.MapActions(p)].String(),
			Uid:            a.dashboard.Uid,
			Title:          a.dashboard.Title,
			Slug:           a.dashboard.Slug,
			IsFolder:       a.dashboard.IsFolder,
			Url:            a.dashboard.GetUrl(),
			Inherited:      false,
		})
	}

	return acl, nil
}

func (a *AccessControlDashboardGuardian) GetACLWithoutDuplicates() ([]*models.DashboardAclInfoDTO, error) {
	return a.GetAcl()
}

func (a *AccessControlDashboardGuardian) GetHiddenACL(cfg *setting.Cfg) ([]*models.DashboardAcl, error) {
	// not used with access control
	return nil, nil
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
