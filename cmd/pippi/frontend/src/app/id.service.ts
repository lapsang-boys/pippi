import { Injectable } from '@angular/core';
import { from, Subject, BehaviorSubject, Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class IdService {
  public readonly _id: BehaviorSubject<string> = new BehaviorSubject("");
  public readonly id: Observable<string> = this._id.asObservable();

  constructor() { }

  set(id: string): void {
    this._id.next(id);
  }
}
