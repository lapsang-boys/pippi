import { Component } from '@angular/core';
import { Buffer } from "buffer";
import  hexdump from "hexdump-nodejs";

@Component({
  selector: '[id="app"]',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'my-app';

  public binary: any;
  getBinary() {
    // @ts-ignore
    window.backend.binary('c279338c4e9a9b3b61317d32f8261c123d3c7d714d01b55c90ee43b07b05ad07').then(result => {
      let buffer = new Buffer(result, 'base64');
      console.log(buffer);
      this.binary = hexdump(buffer);
    });
  }

  public disassembly: any;
  getDisassembly() {
    // @ts-ignore
    window.backend.disassembly('c279338c4e9a9b3b61317d32f8261c123d3c7d714d01b55c90ee43b07b05ad07').then(result => {
      console.log(result);
      this.disassembly = result;
    });
  }

  public file: any;
  getFile() {
    // @ts-ignore
    window.backend.sections('c279338c4e9a9b3b61317d32f8261c123d3c7d714d01b55c90ee43b07b05ad07').then(result => {
      console.log(result);
      this.file = result;
    });
  }
  public extractedStrings: string[] = [];
  getStrings() {
    // @ts-ignore
    window.backend.strings('c279338c4e9a9b3b61317d32f8261c123d3c7d714d01b55c90ee43b07b05ad07').then(result => {
      console.log(result);
      this.extractedStrings = result;
    });
  }
}
