import { Observable } from 'rxjs';

import { ObservableMatchers } from './types';
import { toEmitValues } from './toEmitValues';
import { toEmitValuesWith } from './toEmitValuesWith';

export const matchers: ObservableMatchers<void, Observable<any>> = {
  toEmitValues,
  toEmitValuesWith,
};
