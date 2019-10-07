import { Component, Input, SimpleChange, SimpleChanges, ViewChild } from '@angular/core';

@Component({
  selector: 'app-disassembly',
  templateUrl: './disassembly.component.html',
  styleUrls: ['./disassembly.component.css']
})
export class DisassemblyComponent {
  @ViewChild('editor', {static: false}) editor;

  @Input() data: string;
  private _data: string = "";

  ngOnChanges(changes: SimpleChanges) {
    const data: SimpleChange = changes.data;
    if (!data.currentValue) {
      return;
    }
    this._data = data.currentValue.map(inst => `0x${inst.addr.toString(16)}: ${inst.inst_str}`).join("\n");
    this.editor.getEditor().getSession().setUseWrapMode(true);
  }

  getData() {
    return this._data;
  }
}
