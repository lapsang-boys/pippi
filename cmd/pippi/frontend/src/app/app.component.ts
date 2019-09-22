import { Component } from '@angular/core';

@Component({
  selector: '[id="app"]',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'my-app';

  public extractedStrings: string[] = [];

  onClickMe() {
    // @ts-ignore
    window.backend.basic('c279338c4e9a9b3b61317d32f8261c123d3c7d714d01b55c90ee43b07b05ad07').then(result => {
      console.log(result);
      this.extractedStrings = result;
    });
  }
}
