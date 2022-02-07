import React from 'react';
import { Props as InlineFieldProps } from '@grafana/ui/src/components/Forms/InlineField';

import { InlineField } from '@grafana/ui';

const DEFAULT_LABEL_WIDTH = 18;

export const Field = (props: InlineFieldProps) => {
  return <InlineField labelWidth={DEFAULT_LABEL_WIDTH} {...props} />;
};
