import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SelectIdComponent } from './select-id.component';

describe('SelectIdComponent', () => {
  let component: SelectIdComponent;
  let fixture: ComponentFixture<SelectIdComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SelectIdComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SelectIdComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
