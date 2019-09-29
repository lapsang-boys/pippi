import { Component } from '@angular/core';
import { IdService } from '../id.service';

@Component({
  selector: 'app-select-id',
  templateUrl: './select-id.component.html',
  styleUrls: ['./select-id.component.css']
})
export class SelectIdComponent {
  ids: string[] = [];

  currentId: string = "";
  constructor(
    private idService: IdService,
  ) {
    this.idService.ids.subscribe(ids => {
      this.ids = ids;
    })
    this.idService.id.subscribe(newId => {
      this.currentId = newId
    });
  }

  short(id: string): string {
    // For 1 thousand binaries, there's a 2.93% probability of collision.
    return id.substring(0, 6);
  }

  setId(id: string): void {
    this.idService.set(id);
  }
}
