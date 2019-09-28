import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DisassemblyComponent } from './disassembly.component';

describe('DisassemblyComponent', () => {
  let component: DisassemblyComponent;
  let fixture: ComponentFixture<DisassemblyComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ DisassemblyComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DisassemblyComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
