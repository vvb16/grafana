import { MutableDataFrame } from '@grafana/data';

import { transformResponse, transformToZipkin } from './transforms';
import { traceFrameFields, zipkinResponse } from './testData';

describe('transformResponse', () => {
  it('transforms response', () => {
    const dataFrame = transformResponse(zipkinResponse);

    expect(dataFrame.fields).toMatchObject(traceFrameFields);
  });
  it('converts dataframe to ZipkinSpan[]', () => {
    const dataFrame = transformResponse(zipkinResponse);
    const response = transformToZipkin(new MutableDataFrame(dataFrame));
    expect(response).toMatchObject(zipkinResponse);
  });
});
