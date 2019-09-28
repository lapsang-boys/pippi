import { TestBed } from '@angular/core/testing';

import { IdService } from './id.service';

describe('IdService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: IdService = TestBed.get(IdService);
    expect(service).toBeTruthy();
  });
});
