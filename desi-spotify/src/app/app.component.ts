import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import {PlayerControlsModule} from "./components/player-controls.module";

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, PlayerControlsModule],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss'
})
export class AppComponent {
  title = 'desi-spotify';
}
