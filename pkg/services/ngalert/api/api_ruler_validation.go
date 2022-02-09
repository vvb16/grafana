package api

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/grafana/grafana/pkg/models"
	apimodels "github.com/grafana/grafana/pkg/services/ngalert/api/tooling/definitions"
	ngmodels "github.com/grafana/grafana/pkg/services/ngalert/models"
	"github.com/grafana/grafana/pkg/services/ngalert/store"
)

func validateRuleNode(
	ruleGroupConfig *apimodels.PostableRuleGroupConfig,
	ruleNode *apimodels.PostableExtendedRuleNode,
	orgId int64,
	namespace *models.Folder,
	conditionValidator func(ngmodels.Condition) error,
) (*ngmodels.AlertRule, error) {
	if ruleNode.GrafanaManagedAlert == nil {
		return nil, fmt.Errorf("not Grafana managed alert rule")
	}

	if ruleNode.GrafanaManagedAlert.Title == "" {
		return nil, errors.New("alert rule title cannot be empty")
	}

	if len(ruleNode.GrafanaManagedAlert.Title) > store.AlertRuleMaxTitleLength {
		return nil, fmt.Errorf("alert rule title is too long. Max length is %d", store.AlertRuleMaxTitleLength)
	}

	noDataState := ngmodels.NoData
	if ruleNode.GrafanaManagedAlert.NoDataState != "" {
		var err error
		noDataState, err = ngmodels.NoDataStateFromString(string(ruleNode.GrafanaManagedAlert.NoDataState))
		if err != nil {
			return nil, err
		}
	}

	errorState := ngmodels.AlertingErrState
	if ruleNode.GrafanaManagedAlert.ExecErrState != "" {
		var err error
		errorState, err = ngmodels.ErrStateFromString(string(ruleNode.GrafanaManagedAlert.ExecErrState))
		if err != nil {
			return nil, err
		}
	}

	if len(ruleNode.GrafanaManagedAlert.Data) == 0 {
		return nil, fmt.Errorf("%w: no queries or expressions are found", ngmodels.ErrAlertRuleFailedValidation)
	}

	cond := ngmodels.Condition{
		Condition: ruleNode.GrafanaManagedAlert.Condition,
		OrgID:     orgId,
		Data:      ruleNode.GrafanaManagedAlert.Data,
	}
	if err := conditionValidator(cond); err != nil {
		return nil, fmt.Errorf("failed to validate condition of alert rule %s: %w", ruleNode.GrafanaManagedAlert.Title, err)
	}

	newAlertRule := ngmodels.AlertRule{
		OrgID:           orgId,
		Title:           ruleNode.GrafanaManagedAlert.Title,
		Condition:       ruleNode.GrafanaManagedAlert.Condition,
		Data:            ruleNode.GrafanaManagedAlert.Data,
		UID:             ruleNode.GrafanaManagedAlert.UID,
		IntervalSeconds: int64(time.Duration(ruleGroupConfig.Interval).Seconds()),
		NamespaceUID:    namespace.Uid,
		RuleGroup:       ruleGroupConfig.Name,
		NoDataState:     noDataState,
		ExecErrState:    errorState,
	}

	if ruleNode.ApiRuleNode != nil {
		newAlertRule.For = time.Duration(ruleNode.ApiRuleNode.For)
		newAlertRule.Annotations = ruleNode.ApiRuleNode.Annotations
		newAlertRule.Labels = ruleNode.ApiRuleNode.Labels

		dashUID := ruleNode.ApiRuleNode.Annotations[ngmodels.DashboardUIDAnnotation]
		panelID := ruleNode.ApiRuleNode.Annotations[ngmodels.PanelIDAnnotation]

		if dashUID != "" && panelID == "" || dashUID == "" && panelID != "" {
			return nil, fmt.Errorf("both annotations %s and %s must be specified", ngmodels.DashboardUIDAnnotation, ngmodels.PanelIDAnnotation)
		}

		if dashUID != "" {
			panelIDValue, err := strconv.ParseInt(panelID, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("annotation %s must be a valid integer Panel ID", ngmodels.PanelIDAnnotation)
			}
			newAlertRule.DashboardUID = &dashUID
			newAlertRule.PanelID = &panelIDValue
		}
	}

	return &newAlertRule, nil
}

func validateRuleGroup(
	ruleGroupConfig *apimodels.PostableRuleGroupConfig,
	orgId int64,
	namespace *models.Folder,
	baseInterval time.Duration,
	conditionValidator func(ngmodels.Condition) error,
) ([]*ngmodels.AlertRule, error) {
	if ruleGroupConfig.Name == "" {
		return nil, errors.New("rule group name cannot be empty")
	}

	if len(ruleGroupConfig.Name) > store.AlertRuleMaxRuleGroupNameLength {
		return nil, fmt.Errorf("rule group name is too long. Max length is %d", store.AlertRuleMaxRuleGroupNameLength)
	}

	if ruleGroupConfig.Interval <= 0 {
		return nil, fmt.Errorf("rule evaluation interval must be positive value greater than")
	}

	if int64(time.Duration(ruleGroupConfig.Interval).Seconds())%int64(baseInterval.Seconds()) != 0 {
		return nil, fmt.Errorf("rule evaluation interval should be multiple of the base interval of %d seconds", int64(baseInterval.Seconds()))
	}

	result := make([]*ngmodels.AlertRule, 0, len(ruleGroupConfig.Rules))
	uids := make(map[string]int, cap(result))
	for idx, ruleNode := range ruleGroupConfig.Rules {
		rule, err := validateRuleNode(ruleGroupConfig, &ruleNode, orgId, namespace, conditionValidator)
		// TODO do not stop on the first failure but return all failures
		if err != nil {
			return nil, fmt.Errorf("invalid rule specification at index [%d]: %w", idx, err)
		}
		if rule.UID != "" {
			if existingIdx, ok := uids[rule.UID]; ok {
				return nil, fmt.Errorf("rule [%d] has UID %s that is already assigned to another rule at index %d", idx, rule.UID, existingIdx)
			}
			uids[rule.UID] = idx
		}
		result = append(result, rule)
	}
	return result, nil
}
