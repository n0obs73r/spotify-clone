import {ChangeDetectionStrategy, ChangeDetectorRef, Component} from "@angular/core";
import {PlayerService} from "../services/player.service";

@Component({
  selector: 'app-player-controls',
  templateUrl: './player-controls.component.html',
  styleUrls: ['./player-controls.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class PlayerControlsComponent {
  currentProgress = 0;

  constructor(
    private cdr: ChangeDetectorRef,
    public playerService: PlayerService
  ) { }

  handleProgressBarClick(e: MouseEvent) {
    const target = e.target as HTMLElement;
    if (target.classList.contains('thumb')) {
      e.stopPropagation();
      return;
    }
    this.recalculateCurrentProgress(e.offsetX);
  }

  recalculateCurrentProgress(offsetX: number) {
    const progressBarWidth = 440;
    this.currentProgress = (offsetX / progressBarWidth) * 100;
    this.cdr.markForCheck();
  }

  togglePlayPause() {
    this.playerService.togglePlaying();
  }
}
