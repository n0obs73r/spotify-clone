import { Component, OnDestroy, OnInit } from "@angular/core";
import { PlayerService } from "../services/player.service";
import { Subscription } from "rxjs";

@Component({
  selector: 'app-player-controls',
  templateUrl: './player-controls.component.html',
  styleUrls: ['./player-controls.component.scss']
})
export class PlayerControlsComponent implements OnInit, OnDestroy {
  currentProgress = 0;
  startTime = "0:00";
  endTime = "0:00";
  progressSubscription: Subscription | null = null;
  manualSeeking = false; // Flag to indicate manual seeking

  constructor(public playerService: PlayerService) { }

  ngOnInit() {
    // Subscribe to progress updates
    this.progressSubscription = this.playerService.getProgress().subscribe(progress => {
      if (!this.manualSeeking) {
        this.currentProgress = progress;
        this.updateTime(progress);
      }
    });

    // Subscribe to playing status changes
    this.playerService.getPlaying().subscribe(playing => {
      if (!playing) {
        // Reset time and progress when playback stops
        this.currentProgress = 0;
        this.startTime = "0:00";
        this.endTime = "0:00";
      }
    });
  }

  ngOnDestroy() {
    // Unsubscribe to prevent memory leaks
    this.progressSubscription?.unsubscribe();
  }

  updateTime(progress: number) {
    const totalDuration = this.playerService.getTotalDuration();
    if (totalDuration) {
      const currentTimeInSeconds = (progress / 100) * totalDuration;
      this.startTime = this.formatTime(currentTimeInSeconds);
      this.endTime = this.formatTime(totalDuration - currentTimeInSeconds); // Update endTime
    }
  }

  formatTime(seconds: number): string {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = Math.floor(seconds % 60);
    return `${minutes}:${remainingSeconds < 10 ? '0' : ''}${remainingSeconds}`;
  }

  handleProgressBarInput(event: Event) {
    const value = (event.target as HTMLInputElement)?.value;
    if (value !== null && value !== undefined) {
      this.manualSeeking = true; // Set flag to indicate manual seeking
      this.updateTime(Number(value)); 
      this.currentProgress = Number(value); // Update progress bar position immediately
      this.playerService.seekTo(Number(value)); // Seek to the selected position
    }
  }

  handleProgressBarChange() {
    if (!this.manualSeeking) {
      // If manual seeking flag is not set, update the time and progress
      const totalDuration = this.playerService.getTotalDuration();
      if (totalDuration) {
        const newPosition = (this.currentProgress / 100) * totalDuration;
        this.updateTime(newPosition);
        this.playerService.seekTo(newPosition);
      }
    }
    this.manualSeeking = false; // Reset the manual seeking flag
  }
  

  togglePlayPause() {
    this.playerService.togglePlaying();
  }

  playNextSong() {
    this.playerService.playNextSong();
  }

  playPreviousSong() {
    this.playerService.playPreviousSong();
  }
}
