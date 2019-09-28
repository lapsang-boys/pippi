import { Component, OnInit } from '@angular/core';
import { LoaderService } from '../loader.service';
import { Subject } from 'rxjs';

@Component({
  selector: 'app-loader',
  templateUrl: './loader.component.html',
  styleUrls: ['./loader.component.css']
})
export class LoaderComponent {
  isLoading: Subject<boolean> = this.loaderService.isLoading;
  constructor(private loaderService: LoaderService){}
}
