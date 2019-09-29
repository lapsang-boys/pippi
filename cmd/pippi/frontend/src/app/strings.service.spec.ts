import { TestBed } from '@angular/core/testing';

import { StringsService } from './strings.service';

describe('StringsService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: StringsService = TestBed.get(StringsService);
    expect(service).toBeTruthy();
  });
});
