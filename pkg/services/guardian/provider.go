package guardian

import (
	"context"

	"github.com/grafana/grafana/pkg/services/featuremgmt"

	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/accesscontrol/resourceservices"
	"github.com/grafana/grafana/pkg/services/sqlstore"
)

type Provider struct{}

func ProvideService(features featuremgmt.FeatureToggles, store *sqlstore.SQLStore, ac accesscontrol.AccessControl, permissionServices *resourceservices.ResourceServices) *Provider {
	// TODO: remove this hack
	if features.IsEnabled(featuremgmt.FlagAccesscontrol) {
		New = func(ctx context.Context, dashId int64, orgId int64, user *models.SignedInUser) DashboardGuardian {
			return NewAccessControlDashboardGuardian(ctx, dashId, user, store, ac, permissionServices)
		}
	}
	return &Provider{}
}
