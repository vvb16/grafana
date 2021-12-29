import React from 'react';
import { Permissions } from 'app/core/components/AccessControl';

interface Props {
  id: number;
}

export const AccessControlDashboardPermissions = ({ id }: Props) => (
  <Permissions resource={'dashboards'} resourceId={id} canListUsers={true} canSetPermissions={true} />
);
