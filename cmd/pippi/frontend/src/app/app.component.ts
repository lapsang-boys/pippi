import { Component } from '@angular/core';
import { BinaryService } from './binary.service';
import { DisassemblyService } from './disassembly.service';
import { IdService } from './id.service';
import { LoaderService } from './loader.service';

@Component({
  selector: '[id="app"]',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'Pippi front-end';

  binaryData: any;
  disassemblyData: any;
  id = this.idService.id.subscribe(newId => {
    console.log("here", newId)
    if (newId == "") {
      return;
    }
    console.log("Firing!")
    this.resolveBinary(newId);
    this.resolveDisassembly(newId);
    return newId;
  });

  constructor(
    private loaderService: LoaderService,
    private binaryService: BinaryService,
    private disassemblyService: DisassemblyService,
    private idService: IdService
  ) {}

  resolveBinary(id: string) {
    this.loaderService.show();
    this.binaryService.get(id).then(buffer => {
      this.binaryData = buffer;
      this.loaderService.hide();
    });
  }

  resolveDisassembly(id: string) {
    this.loaderService.show();
    this.disassemblyService.get(id).then(output => {
      this.disassemblyData = output;
      this.loaderService.hide();
    });
  }
}
