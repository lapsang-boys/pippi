import { Component, OnInit } from '@angular/core';
import { BinaryService } from './binary.service';
import { DisassemblyService } from './disassembly.service';
import { IdService } from './id.service';
import { LoaderService } from './loader.service';
import { StringsService } from './strings.service';
import { SelectionService } from './selection.service';

@Component({
  selector: '[id="app"]',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  title = 'Pippi front-end';

  binaryData: any;
  disassemblyData: any;
  stringsData: any;
  disassemblyPos: number = 200;
  constructor(
    private selectionService: SelectionService,
    private loaderService: LoaderService,
    private binaryService: BinaryService,
    private disassemblyService: DisassemblyService,
    private stringsService: StringsService,
    private idService: IdService
  ) {}

  ngOnInit() {
    this.selectionService.selection.subscribe(newSelection => {
      console.log(newSelection);
    });

    this.idService.id.subscribe(newId => {
      if (newId == "") {
        return;
      }
      this.resolveBinary(newId);
      this.resolveDisassembly(newId);
      this.resolveStrings(newId);
      return newId;
    });
  }

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

  resolveStrings(id: string) {
    this.loaderService.show();
    this.stringsService.get(id).then(output => {
      this.stringsData = output;
      this.loaderService.hide();
    });
  }
}
