import React, { FunctionComponent } from 'react';
import { css } from '@emotion/css';

import { AsyncMultiSelect, InlineField, SegmentAsync, Select } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';

import { useDispatch } from '../../../../hooks/useStatelessReducer';
import { useFields } from '../../../../hooks/useFields';
import { TopMetrics } from '../aggregations';
import { changeMetricSetting } from '../state/actions';
import { orderOptions } from '../../BucketAggregationsEditor/utils';

interface Props {
  metric: TopMetrics;
}

const toMultiSelectValue = (value: string): SelectableValue<string> => ({ value, label: value });

export const TopMetricsSettingsEditor: FunctionComponent<Props> = ({ metric }) => {
  const dispatch = useDispatch();
  const getOrderByOptions = useFields(['number', 'date']);
  const getMetricsOptions = useFields(metric.type);

  return (
    <>
      <InlineField label="Metrics" labelWidth={16}>
        <AsyncMultiSelect
          menuShouldPortal
          onChange={(e) =>
            dispatch(
              changeMetricSetting({
                metric,
                settingName: 'metrics',
                newValue: e.map((v) => v.value!),
              })
            )
          }
          loadOptions={getMetricsOptions}
          value={metric.settings?.metrics?.map(toMultiSelectValue)}
          closeMenuOnSelect={false}
          defaultOptions
        />
      </InlineField>
      <InlineField label="Order" labelWidth={16}>
        <Select
          menuShouldPortal
          onChange={(e) => dispatch(changeMetricSetting({ metric, settingName: 'order', newValue: e.value }))}
          options={orderOptions}
          value={metric.settings?.order}
        />
      </InlineField>
      <InlineField
        label="Order By"
        labelWidth={16}
        className={css`
          & > div {
            width: 100%;
          }
        `}
      >
        <SegmentAsync
          className={css`
            margin-right: 0;
          `}
          loadOptions={getOrderByOptions}
          onChange={(e) => dispatch(changeMetricSetting({ metric, settingName: 'orderBy', newValue: e.value }))}
          placeholder="Select Field"
          value={metric.settings?.orderBy}
        />
      </InlineField>
    </>
  );
};
