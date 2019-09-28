import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';

@Injectable()
export class LoaderService {
    private i = 0;
    isLoading = new Subject<boolean>();
    show() {
      this.i += 1;
      console.log("Showing!")
      this.isLoading.next(true);
    }
    hide() {
      this.i -= 1;
      if (this.i == 0) {
        console.log("Hiding!")
        this.isLoading.next(false);
      }
    }
}
