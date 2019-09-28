import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class DisassemblyService {
  get(id: string): Promise<any> {
    // @ts-ignore
    return window.backend.disassembly(id).then(result => {
      return result;
    });
  }
}
