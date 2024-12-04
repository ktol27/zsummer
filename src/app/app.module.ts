import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HttpClientModule } from '@angular/common/http'; // 导入 HttpClientModule

import { AppComponent } from './app.component';

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    HttpClientModule // 将 HttpClientModule 添加到 imports 中
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
