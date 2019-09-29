/**
 * A single tab page. It renders the passed template
 * via the @Input properties by using the ngTemplateOutlet
 * and ngTemplateOutletContext directives.
 */

import { Component, Input } from "@angular/core";

@Component({
  selector: "app-tab",
  templateUrl: "./tab.component.html",
  styles: [`
    div {
      height: 100%;
    }
    `],
})

export class TabComponent {
  @Input("title") title: string;
  @Input() active = false;
}
