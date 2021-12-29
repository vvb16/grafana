package resourcepermissions

import (
	"context"
	"strconv"
	"time"

	"github.com/grafana/grafana/pkg/models"

	"github.com/grafana/grafana/pkg/api/routing"
	"github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/sqlstore"
)

func ProvideServices(sql *sqlstore.SQLStore, router routing.RouteRegister, ac accesscontrol.AccessControl, store accesscontrol.ResourcePermissionsStore) (*Services, error) {
	dashboardsService, err := provideDashboardService(sql, router, ac, store)
	if err != nil {
		return nil, err
	}

	return &Services{services: map[string]*Service{"dashboards": dashboardsService}}, nil
}

type Services struct {
	services map[string]*Service
}

func (s *Services) GetDashboardService() *Service {
	return s.services["dashboards"]
}

func provideDashboardService(sql *sqlstore.SQLStore, router routing.RouteRegister, ac accesscontrol.AccessControl, store accesscontrol.ResourcePermissionsStore) (*Service, error) {

	options := Options{
		Resource: "dashboards",
		ResourceValidator: func(ctx context.Context, orgID int64, resourceID string) error {
			id, err := strconv.ParseInt(resourceID, 10, 64)
			if err != nil {
				return err
			}

			if _, err := sql.GetDashboard(id, orgID, "", ""); err != nil {
				return err
			}
			return nil
		},
		Assignments: Assignments{
			Users:        true,
			Teams:        true,
			BuiltInRoles: true,
		},
		PermissionsToActions: map[string][]string{
			"View":  {"dashboards:read"},
			"Edit":  {"dashboards:read", "dashboards:write", "dashboards:delete"},
			"Admin": {"dashboards:read", "dashboards:write", "dashboards:delete", "dashboards.permissions:read", "dashboards.permissions:write"},
		},
		ReaderRoleName: "Dashboard permission reader",
		WriterRoleName: "Dashboard permission writer",
		RoleGroup:      "Dashboards",
		OnSetUser: func(ctx context.Context, orgID, userID int64, resourceID, permission string) error {
			item := models.DashboardAcl{OrgID: orgID, UserID: userID}
			return onDashboardPermissionUpdated(ctx, sql, resourceID, item, permission)
		},
		OnSetTeam: func(ctx context.Context, orgID, teamID int64, resourceID, permission string) error {
			item := models.DashboardAcl{OrgID: orgID, TeamID: teamID}
			return onDashboardPermissionUpdated(ctx, sql, resourceID, item, permission)
		},
		OnSetBuiltInRole: func(ctx context.Context, orgID int64, builtInRole, resourceID, permission string) error {
			role := models.RoleType(builtInRole)
			item := models.DashboardAcl{OrgID: orgID, Role: &role}
			return onDashboardPermissionUpdated(ctx, sql, resourceID, item, permission)
		},
	}

	return New(options, router, ac, store)
}

func onDashboardPermissionUpdated(ctx context.Context, store *sqlstore.SQLStore, resourceID string, item models.DashboardAcl, permission string) error {
	return store.WithTransactionalDbSession(ctx, func(sess *sqlstore.DBSession) error {
		dashboardID, err := strconv.ParseInt(resourceID, 10, 64)
		if err != nil {
			return err
		}

		item.DashboardID = dashboardID
		has, err := sess.Get(&item)
		if err != nil {
			return err
		}

		if permission == "" {
			rawSQL := `
				DELETE FROM dashboard_acl WHERE id = ?
			`
			if _, err := sess.Exec(rawSQL, item.Id); err != nil {
				return err
			}
			return nil
		}

		item.Updated = time.Now()
		item.Permission = translateDashboardPermission(permission)
		if has {
			rawSQL := `
				UPDATE dashboard_acl
				SET permission = ?, updated = ?
				WHERE id = ?
			`
			if _, err := sess.Exec(rawSQL, item.Permission, item.Updated, item.Id); err != nil {
				return err
			}

			return nil
		}

		item.Created = time.Now()
		sess.Nullable("user_id", "team_id")
		if _, err := sess.Insert(&item); err != nil {
			return err
		}

		return nil
	})
}

func translateDashboardPermission(permission string) models.PermissionType {
	permissionLevel := models.PERMISSION_VIEW
	if permission == models.PERMISSION_EDIT.String() {
		permissionLevel = models.PERMISSION_EDIT
	} else if permission == models.PERMISSION_ADMIN.String() {
		permissionLevel = models.PERMISSION_ADMIN
	}
	return permissionLevel
}
