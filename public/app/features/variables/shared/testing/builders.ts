import { initialQueryVariableModelState } from '../../query/reducer';
import { initialAdHocVariableModelState } from '../../adhoc/reducer';
import { initialDataSourceVariableModelState } from '../../datasource/reducer';
import { initialIntervalVariableModelState } from '../../interval/reducer';
import { initialTextBoxVariableModelState } from '../../textbox/reducer';
import { initialCustomVariableModelState } from '../../custom/reducer';
import { initialConstantVariableModelState } from '../../constant/reducer';

import { MultiVariableBuilder } from './multiVariableBuilder';
import { OptionsVariableBuilder } from './optionsVariableBuilder';
import { DatasourceVariableBuilder } from './datasourceVariableBuilder';
import { IntervalVariableBuilder } from './intervalVariableBuilder';
import { AdHocVariableBuilder } from './adHocVariableBuilder';
import { QueryVariableBuilder } from './queryVariableBuilder';
import { TextBoxVariableBuilder } from './textboxVariableBuilder';

export const adHocBuilder = () => new AdHocVariableBuilder(initialAdHocVariableModelState);
export const intervalBuilder = () => new IntervalVariableBuilder(initialIntervalVariableModelState);
export const datasourceBuilder = () => new DatasourceVariableBuilder(initialDataSourceVariableModelState);
export const queryBuilder = () => new QueryVariableBuilder(initialQueryVariableModelState);
export const textboxBuilder = () => new TextBoxVariableBuilder(initialTextBoxVariableModelState);
export const customBuilder = () => new MultiVariableBuilder(initialCustomVariableModelState);
export const constantBuilder = () => new OptionsVariableBuilder(initialConstantVariableModelState);
