import { Component, OnInit } from '@angular/core';

import { FileUploader } from 'ng2-file-upload/ng2-file-upload';
import { IdService } from '../id.service';

const uploadURL = 'http://localhost:2000/upload';

@Component({
  selector: 'app-upload',
  templateUrl: './upload.component.html',
  styleUrls: ['./upload.component.css']
})
export class UploadComponent implements OnInit {
  public uploader: FileUploader = new FileUploader({url: uploadURL, itemAlias: 'file'});

  constructor(
    private idService: IdService,
  ) { }

  ngOnInit() {
    this.uploader.onAfterAddingFile = (file) => { file.withCredentials = false; }
    this.uploader.onCompleteItem = (item: any, response: any, status: any, header: any) => {
      console.log('FileUpload.uploaded:', item, status, response);
      this.idService.check();
    }
  }
}
