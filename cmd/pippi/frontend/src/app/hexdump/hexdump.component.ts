import { Component, OnInit, ViewChild, Input, OnChanges, SimpleChanges, SimpleChange } from '@angular/core';

import hexdump from "hexdump-nodejs";

import { SelectionService } from '../selection.service';

@Component({
  selector: 'app-hexdump',
  templateUrl: './hexdump.component.html',
  styleUrls: ['./hexdump.component.css']
})
export class HexdumpComponent implements OnChanges {
  @ViewChild('editor', {static: false}) editor;

  @Input() data: string;
  private _data: string = "";

  constructor(
    private selectionService: SelectionService,
  ) {}

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

    this.selectionService.selection.subscribe(newSel => {
      if (newSel.start == 0 && newSel.end == 0) {
        return;
      }
      console.log(this);
      console.log(newSel);

      let row = Math.floor(newSel.start/16)+1;
      // let pos = this.editor.getEditor().getSession().doc.indexToPosition(newSel.start, 0);
      // console.log(pos);
      console.log(row)
      this.editor.getEditor().moveCursorToPosition({column: 0, row: row});
      this.editor.getEditor().scrollToLine(row, true, false, null);
    })

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
