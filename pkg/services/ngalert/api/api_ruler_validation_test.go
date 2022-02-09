package api

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"

	models2 "github.com/grafana/grafana/pkg/models"
	apimodels "github.com/grafana/grafana/pkg/services/ngalert/api/tooling/definitions"
	"github.com/grafana/grafana/pkg/services/ngalert/models"
	"github.com/grafana/grafana/pkg/services/ngalert/store"
	"github.com/grafana/grafana/pkg/util"
)

var allNoData = []apimodels.NoDataState{
	apimodels.OK,
	apimodels.NoData,
	apimodels.Alerting,
}

var allExecError = []apimodels.ExecutionErrorState{
	apimodels.ErrorErrState,
	apimodels.AlertingErrState,
}

var baseInterval = time.Duration(rand.Intn(99)+1) * time.Second

func validRule() apimodels.PostableExtendedRuleNode {
	return apimodels.PostableExtendedRuleNode{
		ApiRuleNode: &apimodels.ApiRuleNode{
			For: model.Duration(rand.Int63n(1000)),
			Labels: map[string]string{
				"test-label": "data",
			},
			Annotations: map[string]string{
				"test-annotation": "data",
			},
		},
		GrafanaManagedAlert: &apimodels.PostableGrafanaRule{
			Title:     fmt.Sprintf("TEST-ALERT-%d", rand.Int63()),
			Condition: "A",
			Data: []models.AlertQuery{
				{
					RefID:     "A",
					QueryType: "TEST",
					RelativeTimeRange: models.RelativeTimeRange{
						From: 10,
						To:   0,
					},
					DatasourceUID: "DATASOURCE_TEST",
					Model:         nil,
				},
			},
			UID:          util.GenerateShortUID(),
			NoDataState:  allNoData[rand.Intn(len(allNoData)-1)],
			ExecErrState: allExecError[rand.Intn(len(allExecError)-1)],
		},
	}
}

func validGroup(rules ...apimodels.PostableExtendedRuleNode) apimodels.PostableRuleGroupConfig {
	return apimodels.PostableRuleGroupConfig{
		Name:     "TEST-ALERTS",
		Interval: model.Duration(baseInterval),
		Rules:    rules,
	}
}

func randFolder() *models2.Folder {
	return &models2.Folder{
		Id:        rand.Int63(),
		Uid:       util.GenerateShortUID(),
		Title:     "TEST-FOLDER",
		Url:       "",
		Version:   0,
		Created:   time.Time{},
		Updated:   time.Time{},
		UpdatedBy: 0,
		CreatedBy: 0,
		HasAcl:    false,
	}
}

func TestValidateRuleGroup(t *testing.T) {
	orgId := rand.Int63()
	folder := randFolder()
	t.Run("should validate struct and rules", func(t *testing.T) {
		rules := make([]apimodels.PostableExtendedRuleNode, 0, rand.Intn(4)+1)
		for i := 0; i < cap(rules); i++ {
			rules = append(rules, validRule())
		}
		g := validGroup(rules...)
		conditionValidations := 0

		alerts, err := validateRuleGroup(&g, orgId, folder, baseInterval, func(condition models.Condition) error {
			conditionValidations++
			return nil
		})
		require.NoError(t, err)
		require.Len(t, alerts, len(rules))
		require.Equal(t, len(rules), conditionValidations)
	})
}

func TestValidateRuleGroupFailures(t *testing.T) {
	orgId := rand.Int63()
	folder := randFolder()

	testCases := []struct {
		name   string
		group  func() *apimodels.PostableRuleGroupConfig
		assert func(t *testing.T, apiModel *apimodels.PostableRuleGroupConfig, err error)
	}{
		{
			name: "fail if title is empty",
			group: func() *apimodels.PostableRuleGroupConfig {
				g := validGroup()
				g.Name = ""
				return &g
			},
		},
		{
			name: "fail if title is too long",
			group: func() *apimodels.PostableRuleGroupConfig {
				g := validGroup()
				for len(g.Name) < store.AlertRuleMaxRuleGroupNameLength {
					g.Name += g.Name
				}
				return &g
			},
		},
		{
			name: "fail if interval is 0",
			group: func() *apimodels.PostableRuleGroupConfig {
				g := validGroup()
				g.Interval = model.Duration(0)
				return &g
			},
		},
		{
			name: "fail if interval is not aligned with base interval",
			group: func() *apimodels.PostableRuleGroupConfig {
				g := validGroup()
				g.Interval = model.Duration(baseInterval + time.Duration(rand.Intn(10)+1)*time.Second)
				return &g
			},
		},
		{
			name: "fail if two rules have same UID",
			group: func() *apimodels.PostableRuleGroupConfig {
				r1 := validRule()
				r2 := validRule()
				uid := util.GenerateShortUID()
				r1.GrafanaManagedAlert.UID = uid
				r2.GrafanaManagedAlert.UID = uid
				g := validGroup(r1, r2)
				return &g
			},
			assert: func(t *testing.T, apiModel *apimodels.PostableRuleGroupConfig, err error) {
				require.Contains(t, err.Error(), apiModel.Rules[0].GrafanaManagedAlert.UID)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			g := testCase.group()
			_, err := validateRuleGroup(g, orgId, folder, baseInterval, func(condition models.Condition) error {
				return nil
			})
			require.Error(t, err)
			if testCase.assert != nil {
				testCase.assert(t, g, err)
			}
		})
	}
}

func TestValidateRuleNode(t *testing.T) {
	orgId := rand.Int63()
	folder := randFolder()
	g := validGroup()

	testCases := []struct {
		name   string
		rule   func() *apimodels.PostableExtendedRuleNode
		assert func(t *testing.T, model *apimodels.PostableExtendedRuleNode, rule *models.AlertRule)
	}{
		{
			name: "coverts api model to AlertRule",
			rule: func() *apimodels.PostableExtendedRuleNode {
				r := validRule()
				return &r
			},
			assert: func(t *testing.T, api *apimodels.PostableExtendedRuleNode, alert *models.AlertRule) {
				require.Equal(t, int64(0), alert.ID)
				require.Equal(t, orgId, alert.OrgID)
				require.Equal(t, api.GrafanaManagedAlert.Title, alert.Title)
				require.Equal(t, api.GrafanaManagedAlert.Condition, alert.Condition)
				require.Equal(t, api.GrafanaManagedAlert.Data, alert.Data)
				require.Equal(t, time.Time{}, alert.Updated)
				require.Equal(t, int64(time.Duration(g.Interval).Seconds()), alert.IntervalSeconds)
				require.Equal(t, int64(0), alert.Version)
				require.Equal(t, api.GrafanaManagedAlert.UID, alert.UID)
				require.Equal(t, folder.Uid, alert.NamespaceUID)
				require.Nil(t, alert.DashboardUID)
				require.Nil(t, alert.PanelID)
				require.Equal(t, g.Name, alert.RuleGroup)
				require.Equal(t, models.NoDataState(api.GrafanaManagedAlert.NoDataState), alert.NoDataState)
				require.Equal(t, models.ExecutionErrorState(api.GrafanaManagedAlert.ExecErrState), alert.ExecErrState)
				require.Equal(t, time.Duration(api.ApiRuleNode.For), alert.For)
				require.Equal(t, api.ApiRuleNode.Annotations, alert.Annotations)
				require.Equal(t, api.ApiRuleNode.Labels, alert.Labels)
			},
		},
		{
			name: "coverts api without ApiRuleNode",
			rule: func() *apimodels.PostableExtendedRuleNode {
				r := validRule()
				r.ApiRuleNode = nil
				return &r
			},
			assert: func(t *testing.T, api *apimodels.PostableExtendedRuleNode, alert *models.AlertRule) {
				require.Equal(t, time.Duration(0), alert.For)
				require.Nil(t, alert.Annotations)
				require.Nil(t, alert.Labels)
			},
		},
		{
			name: "defaults to NoData if NoDataState is empty",
			rule: func() *apimodels.PostableExtendedRuleNode {
				r := validRule()
				r.GrafanaManagedAlert.NoDataState = ""
				return &r
			},
			assert: func(t *testing.T, api *apimodels.PostableExtendedRuleNode, alert *models.AlertRule) {
				require.Equal(t, models.NoData, alert.NoDataState)
			},
		},
		{
			name: "defaults to Alerting if ExecErrState is empty",
			rule: func() *apimodels.PostableExtendedRuleNode {
				r := validRule()
				r.GrafanaManagedAlert.ExecErrState = ""
				return &r
			},
			assert: func(t *testing.T, api *apimodels.PostableExtendedRuleNode, alert *models.AlertRule) {
				require.Equal(t, models.AlertingErrState, alert.ExecErrState)
			},
		},
		{
			name: "extracts Dashboard UID and Panel Id from annotations",
			rule: func() *apimodels.PostableExtendedRuleNode {
				r := validRule()
				r.ApiRuleNode.Annotations = map[string]string{
					models.DashboardUIDAnnotation: util.GenerateShortUID(),
					models.PanelIDAnnotation:      strconv.Itoa(rand.Int()),
				}
				return &r
			},
			assert: func(t *testing.T, api *apimodels.PostableExtendedRuleNode, alert *models.AlertRule) {
				require.Equal(t, api.ApiRuleNode.Annotations[models.DashboardUIDAnnotation], *alert.DashboardUID)
				panelId, err := strconv.Atoi(api.ApiRuleNode.Annotations[models.PanelIDAnnotation])
				require.NoError(t, err)
				require.Equal(t, int64(panelId), *alert.PanelID)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			r := testCase.rule()
			alert, err := validateRuleNode(&g, r, orgId, folder, func(condition models.Condition) error {
				return nil
			})
			require.NoError(t, err)
			testCase.assert(t, r, alert)
		})
	}
}

func TestValidateRuleNodeFailures(t *testing.T) {
	orgId := rand.Int63()
	folder := randFolder()
	g := validGroup()
	successValidation := func(condition models.Condition) error {
		return nil
	}

	testCases := []struct {
		name                string
		rule                func() *apimodels.PostableExtendedRuleNode
		conditionValidation func(condition models.Condition) error
		assert              func(t *testing.T, model *apimodels.PostableExtendedRuleNode, err error)
	}{
		{
			name: "fail if GrafanaManagedAlert is not specified",
			rule: func() *apimodels.PostableExtendedRuleNode {
				r := validRule()
				r.GrafanaManagedAlert = nil
				return &r
			},
		},
		{
			name: "fail if title is empty",
			rule: func() *apimodels.PostableExtendedRuleNode {
				r := validRule()
				r.GrafanaManagedAlert.Title = ""
				return &r
			},
		},
		{
			name: "fail if title is too long",
			rule: func() *apimodels.PostableExtendedRuleNode {
				r := validRule()
				for len(r.GrafanaManagedAlert.Title) < store.AlertRuleMaxTitleLength {
					r.GrafanaManagedAlert.Title += r.GrafanaManagedAlert.Title
				}
				return &r
			},
		},
		{
			name: "fail if NoDataState is not known",
			rule: func() *apimodels.PostableExtendedRuleNode {
				r := validRule()
				r.GrafanaManagedAlert.NoDataState = apimodels.NoDataState(util.GenerateShortUID())
				return &r
			},
		},
		{
			name: "fail if ExecErrState is not known",
			rule: func() *apimodels.PostableExtendedRuleNode {
				r := validRule()
				r.GrafanaManagedAlert.ExecErrState = apimodels.ExecutionErrorState(util.GenerateShortUID())
				return &r
			},
		},
		{
			name: "fail if there are not data (nil)",
			rule: func() *apimodels.PostableExtendedRuleNode {
				r := validRule()
				r.GrafanaManagedAlert.Data = nil
				return &r
			},
		},
		{
			name: "fail if there are not data (empty)",
			rule: func() *apimodels.PostableExtendedRuleNode {
				r := validRule()
				r.GrafanaManagedAlert.Data = make([]models.AlertQuery, 0, 1)
				return &r
			},
		},
		{
			name: "fail if validator function returns error",
			rule: func() *apimodels.PostableExtendedRuleNode {
				r := validRule()
				return &r
			},
			conditionValidation: func(condition models.Condition) error {
				return errors.New("BAD alert condition")
			},
		},
		{
			name: "fail if Dashboard UID is specified but not Panel ID",
			rule: func() *apimodels.PostableExtendedRuleNode {
				r := validRule()
				r.ApiRuleNode.Annotations = map[string]string{
					models.DashboardUIDAnnotation: util.GenerateShortUID(),
				}
				return &r
			},
		},
		{
			name: "fail if Dashboard UID is specified and Panel ID is NaN",
			rule: func() *apimodels.PostableExtendedRuleNode {
				r := validRule()
				r.ApiRuleNode.Annotations = map[string]string{
					models.DashboardUIDAnnotation: util.GenerateShortUID(),
					models.PanelIDAnnotation:      util.GenerateShortUID(),
				}
				return &r
			},
		},
		{
			name: "fail if PanelID is specified but not Dashboard UID ",
			rule: func() *apimodels.PostableExtendedRuleNode {
				r := validRule()
				r.ApiRuleNode.Annotations = map[string]string{
					models.PanelIDAnnotation: "0",
				}
				return &r
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			r := testCase.rule()
			f := successValidation
			if testCase.conditionValidation != nil {
				f = testCase.conditionValidation
			}
			_, err := validateRuleNode(&g, r, orgId, folder, f)
			require.Error(t, err)
			if testCase.assert != nil {
				testCase.assert(t, r, err)
			}
		})
	}
}
