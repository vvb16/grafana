package guardian

import (
	"context"

	"github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/sqlstore"

	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/setting"
)

type Provider struct{}

func ProvideService(cfg *setting.Cfg, store *sqlstore.SQLStore, ac accesscontrol.AccessControl) *Provider {
	// TODO: remove this hack
	if cfg.FeatureToggles["accesscontrol"] {
		New = func(ctx context.Context, dashId int64, orgId int64, user *models.SignedInUser) DashboardGuardian {
			return &AccessControlDashboardGuardian{ctx: ctx, dashboardID: dashId, user: user, ac: ac, store: store}
		}
	}
	return &Provider{}
}
