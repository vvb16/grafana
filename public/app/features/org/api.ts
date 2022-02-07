import { getBackendSrv } from '@grafana/runtime';
import { appEvents } from 'app/core/core';
import { AppEvents } from '@grafana/data';

import { FormModel } from './UserInviteForm';

export const userInviteSubmit = async (formData: FormModel) => {
  try {
    await getBackendSrv().post('/api/org/invites', formData);
  } catch (err) {
    appEvents.emit(AppEvents.alertError, ['Failed to send invitation.', err.message]);
  }
};
