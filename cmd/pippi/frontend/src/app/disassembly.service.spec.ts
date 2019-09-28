import { TestBed } from '@angular/core/testing';

import { DisassemblyService } from './disassembly.service';

describe('DisassemblyService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: DisassemblyService = TestBed.get(DisassemblyService);
    expect(service).toBeTruthy();
  });
});
