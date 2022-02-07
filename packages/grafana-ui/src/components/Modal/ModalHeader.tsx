import React from 'react';

import { IconName } from '../../types';
import { useStyles2 } from '../../themes';

import { getModalStyles } from './getModalStyles';

interface Props {
  title: string;
  /** @deprecated */
  icon?: IconName;
  /** @deprecated */
  iconTooltip?: string;
}

/** @internal */
export const ModalHeader: React.FC<Props> = ({ icon, iconTooltip, title, children }) => {
  const styles = useStyles2(getModalStyles);

  return (
    <>
      <h2 className={styles.modalHeaderTitle}>{title}</h2>
      {children}
    </>
  );
};
