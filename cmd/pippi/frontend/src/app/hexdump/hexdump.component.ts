import { Component, OnInit, ViewChild, Input, OnChanges, SimpleChanges, SimpleChange } from '@angular/core';
import hexdump from "hexdump-nodejs";

@Component({
  selector: 'app-hexdump',
  templateUrl: './hexdump.component.html',
  styleUrls: ['./hexdump.component.css']
})
export class HexdumpComponent implements OnChanges {
  @ViewChild('editor', {static: false}) editor;

  @Input() data: string;
  private _data: string = "";

  ngOnChanges(changes: SimpleChanges) {
    const data: SimpleChange = changes.data;
    if (!data.currentValue) {
      return;
    }
    let buffer = data.currentValue;
    this._data = hexdump(buffer);
  }

  ngAfterViewInit() {
    this.editor.setTheme("github");

    this.editor.getEditor().commands.addCommand({
        name: "xref",
        bindKey: "Ctrl-.",
        exec: function(editor) {
          var selection = editor.getSelection();
          var anchor = selection.getSelectionAnchor();
          var token = editor.session.getTokenAt(anchor.row, anchor.column+1);
          console.log(token);
        },
        readOnly: true
    })
  }

  getData() {
    return this._data;
  }
}
