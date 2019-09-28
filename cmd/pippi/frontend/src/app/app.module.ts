import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { APP_BASE_HREF } from '@angular/common';

import { FormsModule } from '@angular/forms';
import { FileUploadModule } from 'ng2-file-upload';

import { AceEditorModule } from 'ng2-ace-editor';

import { AngularSplitModule } from "angular-split";

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { HexdumpComponent } from './hexdump/hexdump.component';
import { ControlBarComponent } from './control-bar/control-bar.component';
import { DisassemblyComponent } from './disassembly/disassembly.component';
import { SelectIdComponent } from './select-id/select-id.component';
import { LoaderComponent } from './loader/loader.component';
import { LoaderService } from './loader.service';

@NgModule({
  declarations: [
    AppComponent,
    HexdumpComponent,
    ControlBarComponent,
    DisassemblyComponent,
    SelectIdComponent,
    LoaderComponent,
  ],
  imports: [
    AngularSplitModule.forRoot(),
    BrowserModule,
    FileUploadModule,
    FormsModule,
    AppRoutingModule,
    AceEditorModule
  ],
  providers: [LoaderService, {provide: APP_BASE_HREF, useValue : '/' }],
  bootstrap: [AppComponent]
})
export class AppModule { }
