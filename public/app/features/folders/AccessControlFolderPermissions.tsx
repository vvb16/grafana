import React from 'react';
import { Permissions } from 'app/core/components/AccessControl';
import { connect, ConnectedProps } from 'react-redux';

import Page from 'app/core/components/Page/Page';
import { getNavModel } from 'app/core/selectors/navModel';
import { contextSrv } from 'app/core/core';
import { getLoadingNav } from './state/navModel';
import { AccessControlAction, StoreState } from 'app/types';
import { GrafanaRouteComponentProps } from 'app/core/navigation/types';

interface RouteProps extends GrafanaRouteComponentProps<{ uid: string }> {}

function mapStateToProps(state: StoreState, props: RouteProps) {
  const uid = props.match.params.uid;
  return {
    resourceId: state.folder.id,
    navModel: getNavModel(state.navIndex, `folder-permissions-${uid}`, getLoadingNav(1)),
  };
}

const connector = connect(mapStateToProps);
export type Props = ConnectedProps<typeof connector>;

export const AccessControlFolderPermissions = ({ resourceId, navModel }: Props) => {
  const canListUsers = contextSrv.hasPermission(AccessControlAction.OrgUsersRead);
  const canSetPermissions = contextSrv.hasPermission(AccessControlAction.DashboardsPermissionsWrite);

  return (
    <Page navModel={navModel}>
      <Page.Contents>
        <Permissions
          resource="dashboards"
          resourceId={resourceId}
          canListUsers={canListUsers}
          canSetPermissions={canSetPermissions}
        />
      </Page.Contents>
    </Page>
  );
};

export default connector(AccessControlFolderPermissions);
