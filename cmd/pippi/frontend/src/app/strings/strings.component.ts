import { Component, Input, SimpleChange, SimpleChanges, ViewChild } from '@angular/core';
import { SelectionService } from '../selection.service';

@Component({
  selector: 'app-strings',
  templateUrl: './strings.component.html',
  styleUrls: ['./strings.component.css']
})
export class StringsComponent {
  @ViewChild('editor', {static: false}) editor;

  @Input() data: string;
  private _data: string = "";

  constructor(
    private selectionService: SelectionService,
  ) { }

  ngOnChanges(changes: SimpleChanges) {
    const data: SimpleChange = changes.data;
    if (!data.currentValue) {
      return;
    }
    this._data = data.currentValue.map(s => `${s.location}: ${s.raw_string}`).join("\n");
    this.editor.getEditor().getSession().setUseWrapMode(true);
  }

  ngAfterViewInit() {
    this.editor.setTheme("github");

    var that = this;
    this.editor.getEditor().commands.addCommand({
        name: "xref",
        bindKey: "Ctrl-.",
        exec: function(editor) {
          var selection = editor.getSelection();
          var anchor = selection.getSelectionAnchor();
          // var token = editor.session.getTokenAt(anchor.row, anchor.column+1);
          var numberToken = editor.session.getTokenAt(anchor.row, 0);
          var stringToken = editor.session.getTokenAt(anchor.row, numberToken.value.length+1);
          if (numberToken.type == "constant.character.decimal.assembly") {
            console.log(numberToken);
            console.log(stringToken);
            let content = stringToken.value.substring(2);
            let start = parseInt(numberToken.value);
            that.selectionService.set(start, start+content.length, content);
          }
        },
        readOnly: true
    })
  }

  getData() {
    return this._data;
  }
}
