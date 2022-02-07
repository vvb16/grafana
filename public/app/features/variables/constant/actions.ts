import { ThunkResult } from 'app/types';

import { validateVariableSelectionState } from '../state/actions';
import { toVariablePayload, VariableIdentifier } from '../state/types';

import { createConstantOptionsFromQuery } from './reducer';

export const updateConstantVariableOptions = (identifier: VariableIdentifier): ThunkResult<void> => {
  return async (dispatch) => {
    await dispatch(createConstantOptionsFromQuery(toVariablePayload(identifier)));
    await dispatch(validateVariableSelectionState(identifier));
  };
};
