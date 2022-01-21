package resourcepermissions

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/grafana/grafana/pkg/models"

	"github.com/grafana/grafana/pkg/api/routing"
	ac "github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/sqlstore"
)

var dashboardsView = []string{ac.ActionDashboardsRead}
var dashboardsEdit = append(dashboardsView, []string{ac.ActionDashboardsWrite, ac.ActionDashboardsDelete}...)
var dashboardsAdmin = append(dashboardsEdit, []string{ac.ActionDashboardsPermissionsRead, ac.ActionDashboardsPermissionsWrite}...)
var foldersView = []string{ac.ActionFoldersRead}
var foldersEdit = append(foldersView, []string{ac.ActionFoldersWrite, ac.ActionFoldersDelete}...)
var foldersAdmin = append(foldersEdit, []string{ac.ActionFoldersPermissionsRead, ac.ActionFoldersPermissionsWrite}...)

func ProvideServices(sql *sqlstore.SQLStore, router routing.RouteRegister, ac ac.AccessControl, store ac.ResourcePermissionsStore) (*Services, error) {
	dashboardsService, err := provideDashboardService(sql, router, ac, store)
	if err != nil {
		return nil, err
	}

	folderService, err := provideFolderService(sql, router, ac, store)
	if err != nil {
		return nil, err
	}

	return &Services{services: map[string]*Service{
		"folders":    folderService,
		"dashboards": dashboardsService,
	}}, nil
}

type Services struct {
	services map[string]*Service
}

func (s *Services) GetDashboardService() *Service {
	return s.services["dashboards"]
}

func (s *Services) GetFolderService() *Service {
	return s.services["folders"]
}

func provideDashboardService(sql *sqlstore.SQLStore, router routing.RouteRegister, accesscontrol ac.AccessControl, store ac.ResourcePermissionsStore) (*Service, error) {
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
			"View":  dashboardsView,
			"Edit":  dashboardsEdit,
			"Admin": dashboardsAdmin,
		},
		ReaderRoleName: "Dashboard permission reader",
		WriterRoleName: "Dashboard permission writer",
		RoleGroup:      "Dashboards",
	}

	return New(options, router, accesscontrol, store)
}

func provideFolderService(sql *sqlstore.SQLStore, router routing.RouteRegister, accesscontrol ac.AccessControl, store ac.ResourcePermissionsStore) (*Service, error) {
	options := Options{
		Resource: "folders",
		ResourceValidator: func(ctx context.Context, orgID int64, resourceID string) error {
			id, err := strconv.ParseInt(resourceID, 10, 64)
			if err != nil {
				return err
			}
			if dashboard, err := sql.GetDashboard(id, orgID, "", ""); err != nil {
				return err
			} else if !dashboard.IsFolder {
				return errors.New("not found")
			}

			return nil
		},
		Assignments: Assignments{
			Users:        true,
			Teams:        true,
			BuiltInRoles: true,
		},
		PermissionsToActions: map[string][]string{
			"View":  append(dashboardsView, foldersView...),
			"Edit":  append(dashboardsEdit, foldersEdit...),
			"Admin": append(dashboardsAdmin, foldersAdmin...),
		},
		ReaderRoleName: "Folder permission reader",
		WriterRoleName: "Folder permission writer",
		RoleGroup:      "Folders",
	}

	return New(options, router, accesscontrol, store)
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
