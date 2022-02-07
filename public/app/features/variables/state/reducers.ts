import { combineReducers } from '@reduxjs/toolkit';

import { optionsPickerReducer } from '../pickers/OptionsPicker/reducer';
import { variableEditorReducer } from '../editor/reducer';
import { variableInspectReducer } from '../inspect/reducer';

import { variablesReducer } from './variablesReducer';
import { transactionReducer } from './transactionReducer';

export const templatingReducers = combineReducers({
  editor: variableEditorReducer,
  variables: variablesReducer,
  optionsPicker: optionsPickerReducer,
  transaction: transactionReducer,
  inspect: variableInspectReducer,
});

export type TemplatingState = ReturnType<typeof templatingReducers>;

export default {
  templating: templatingReducers,
};
