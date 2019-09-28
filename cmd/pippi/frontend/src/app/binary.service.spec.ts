import { TestBed } from '@angular/core/testing';

import { BinaryService } from './binary.service';

describe('BinaryService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: BinaryService = TestBed.get(BinaryService);
    expect(service).toBeTruthy();
  });
});
