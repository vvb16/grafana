import { ThunkResult } from 'app/types';

import { validateVariableSelectionState } from '../state/actions';
import { toVariablePayload, VariableIdentifier } from '../state/types';

import { createCustomOptionsFromQuery } from './reducer';

export const updateCustomVariableOptions = (identifier: VariableIdentifier): ThunkResult<void> => {
  return async (dispatch) => {
    await dispatch(createCustomOptionsFromQuery(toVariablePayload(identifier)));
    await dispatch(validateVariableSelectionState(identifier));
  };
};
