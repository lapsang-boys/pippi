import { Injectable } from '@angular/core';
import { Buffer } from "buffer";

@Injectable({
  providedIn: 'root'
})
export class BinaryService {
  get(id: string): Promise<Buffer> {
    console.log(id)
    // @ts-ignore
    return window.backend.binary(id).then(result => {
      let buffer = new Buffer(result, 'base64');
      return buffer;
    });
  }
}
