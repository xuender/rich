import { CommonModule } from "@angular/common";
import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";

import { ConcatPipe } from "./concat.pipe"

const pipes = [
  ConcatPipe,
]
@NgModule({
  imports: [
    CommonModule
  ],
  declarations: [...pipes,],
  exports: [...pipes],
  schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class PipeModule { }