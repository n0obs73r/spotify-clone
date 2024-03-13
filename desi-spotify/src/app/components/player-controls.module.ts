import {NgModule} from "@angular/core";
import {PlayerControlsComponent} from "./player-controls.component";
import {CommonModule} from "@angular/common";

@NgModule({
  declarations: [PlayerControlsComponent],
  imports: [CommonModule],
  exports: [PlayerControlsComponent]
})
export class PlayerControlsModule { }
