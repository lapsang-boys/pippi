import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { StringsComponent } from './strings.component';

describe('StringsComponent', () => {
  let component: StringsComponent;
  let fixture: ComponentFixture<StringsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ StringsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StringsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
