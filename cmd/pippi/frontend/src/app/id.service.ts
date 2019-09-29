import { Injectable } from '@angular/core';
import { from, Subject, BehaviorSubject, Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class IdService {
  public readonly _id: BehaviorSubject<string> = new BehaviorSubject("");
  public readonly id: Observable<string> = this._id.asObservable();

  //@ts-ignore
  public readonly _ids: BehaviorSubject<string[]> = new BehaviorSubject([]);
  public readonly ids: Observable<string[]> = this._ids.asObservable();

  constructor() {
    //@ts-ignore
    window.backend.listIds().then(ids => {
      this._ids.next(ids);
    })
  }

  set(id: string): void {
    this._id.next(id);
  }

  check() {
    //@ts-ignore
    window.backend.listIds().then(ids => {
      this._ids.next(ids);
    })
  }
}
