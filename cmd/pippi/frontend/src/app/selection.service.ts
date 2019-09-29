import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { Selection } from './selection';

@Injectable({
  providedIn: 'root'
})
export class SelectionService {
  public readonly _selection: BehaviorSubject<Selection> = new BehaviorSubject(new Selection());
  public readonly selection: Observable<Selection> = this._selection.asObservable();

  constructor() {}

  set(
    start?,
    end?,
    content?,
  ) {
    this._selection.next(new Selection(start, end, content));
  }
}
