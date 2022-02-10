import { Alert, useStyles2 } from '@grafana/ui';
import React from 'react';
import { config } from '@grafana/runtime/src';
import { css } from '@emotion/css';

export interface Props {
  showPreviews?: boolean;
}

const getText = (requiredImageRendererPluginVersion?: string) => {
  if (requiredImageRendererPluginVersion) {
    return {
      title: 'Image renderer plugin needs to be updated',
      beforeLink: 'You must update the ',
      link: 'Grafana image renderer plugin',
      afterLink: ` to version ${requiredImageRendererPluginVersion} to enable dashboard previews. Please contact your Grafana administrator to update the plugin`,
    };
  }

  return {
    title: 'Image renderer plugin not installed',
    beforeLink: 'You must install the ',
    link: 'Grafana image renderer plugin',
    afterLink: ` to enable dashboard previews. Please contact your Grafana administrator to install the plugin.`,
  };
};

export const PreviewsSystemRequirements = ({ showPreviews }: Props) => {
  const styles = useStyles2(getStyles);

  const previewsEnabled = config.featureToggles.dashboardPreviews;
  const rendererAvailable = config.rendererAvailable;

  const { systemRequirements, thumbnailsExist } = config.dashboardPreviews;

  const arePreviewsEnabled = previewsEnabled && showPreviews;
  const areRequirementsMet = (rendererAvailable && systemRequirements.met) || thumbnailsExist;
  const shouldDisplayRequirements = arePreviewsEnabled && !areRequirementsMet;

  const text = getText(systemRequirements.requiredImageRendererPluginVersion);

  return (
    shouldDisplayRequirements && (
      <div className={styles.wrapper}>
        <Alert className={styles.alert} severity="info" title={text.title}>
          <>{text.beforeLink}</>
          <a
            href="https://grafana.com/grafana/plugins/grafana-image-renderer"
            target="_blank"
            rel="noopener noreferrer"
            className="external-link"
          >
            {text.link}
          </a>
          {text.afterLink}
        </Alert>
      </div>
    )
  );
};

const getStyles = () => {
  return {
    wrapper: css`
      display: flex;
      justify-content: center;
    `,
    alert: css`
      max-width: 800px;
    `,
  };
};
