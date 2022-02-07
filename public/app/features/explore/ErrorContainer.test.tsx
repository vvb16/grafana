import React from 'react';
import { shallow } from 'enzyme';

import { DataQueryError } from '@grafana/data';

import { ErrorContainer } from './ErrorContainer';

const makeError = (propOverrides?: DataQueryError): DataQueryError => {
  const queryError: DataQueryError = {
    data: {
      message: 'Error data message',
      error: 'Error data content',
    },
    message: 'Error message',
    status: 'Error status',
    statusText: 'Error status text',
    refId: 'A',
  };
  Object.assign(queryError, propOverrides);
  return queryError;
};

const setup = (propOverrides?: object) => {
  const props = {
    queryError: makeError(propOverrides),
  };

  const wrapper = shallow(<ErrorContainer {...props} />);
  return wrapper;
};

describe('ErrorContainer', () => {
  it('should render component', () => {
    const wrapper = setup();
    expect(wrapper).toMatchSnapshot();
  });
});
