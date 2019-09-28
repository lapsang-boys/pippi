import { Component, OnInit } from '@angular/core';

import { FileUploader } from 'ng2-file-upload/ng2-file-upload';

const uploadURL = 'http://localhost:2000/upload';

@Component({
  selector: 'app-control-bar',
  templateUrl: './control-bar.component.html',
  styleUrls: ['./control-bar.component.css']
})
export class ControlBarComponent {
  public uploader: FileUploader = new FileUploader({url: uploadURL, itemAlias: 'file'});

  ngOnInit() {
    this.uploader.onAfterAddingFile = (file) => { file.withCredentials = false; }
    this.uploader.onCompleteItem = (item: any, response: any, status: any, header: any) => {
      console.log('FileUpload.uploaded:', item, status, response);
    }
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
