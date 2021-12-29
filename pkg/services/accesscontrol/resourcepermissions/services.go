package resourcepermissions

import (
	"context"
	"strconv"

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
		ReaderRoleName:   "Dashboard permission reader",
		WriterRoleName:   "Dashboard permission writer",
		RoleGroup:        "Dashboards",
		OnSetUser:        nil,
		OnSetTeam:        nil,
		OnSetBuiltInRole: nil,
	}

	return New(options, router, ac, store)
}
