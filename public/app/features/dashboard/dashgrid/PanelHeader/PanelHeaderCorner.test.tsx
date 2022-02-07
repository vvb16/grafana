import React from 'react';
import { shallow } from 'enzyme';

import { PanelModel } from '../../state';

import { PanelHeaderCorner } from './PanelHeaderCorner';

describe('Render', () => {
  it('should render component', () => {
    const panel = new PanelModel({});
    const wrapper = shallow(<PanelHeaderCorner panel={panel} />);
    const instance = wrapper.instance() as PanelHeaderCorner;

    expect(instance.getInfoContent()).toBeDefined();
  });
});
