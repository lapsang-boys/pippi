import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class StringsService {
  get(id: string): Promise<any> {
    // @ts-ignore
    return window.backend.strings(id).then(result => {
      console.log(result)
      return result;
    });
  }
}
