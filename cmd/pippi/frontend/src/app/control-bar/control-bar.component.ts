import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-control-bar',
  templateUrl: './control-bar.component.html',
  styleUrls: ['./control-bar.component.css']
})
export class ControlBarComponent {
  public file: any;
  getFile() {
    // @ts-ignore
    window.backend.sections('c279338c4e9a9b3b61317d32f8261c123d3c7d714d01b55c90ee43b07b05ad07').then(result => {
      this.file = result;
    });
  }
  public extractedStrings: string[] = [];
  getStrings() {
    // @ts-ignore
    window.backend.strings('c279338c4e9a9b3b61317d32f8261c123d3c7d714d01b55c90ee43b07b05ad07').then(result => {
      this.extractedStrings = result;
    });
  }
}
