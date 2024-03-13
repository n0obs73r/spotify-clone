import { Injectable} from "@angular/core";
import {BehaviorSubject} from "rxjs";

@Injectable({
  providedIn: 'root'
})
export class PlayerService {
  playing$ = new BehaviorSubject<boolean>(false);
  currentAudio: HTMLAudioElement | null = null;

  constructor() {
    this.playing$.subscribe(playing => {
      if (this.currentAudio) {
        if (playing) {
          this.currentAudio.play();
        } else {
          this.currentAudio.pause();
        }
      }
    });
  }

  setPlaying(playing: boolean) {
    this.playing$.next(playing);
  }

  togglePlaying() {
    this.playing$.next(!this.playing$.getValue());
  }

  getPlaying() {
    return this.playing$.asObservable();
  }

  setCurrentAudio(audio: HTMLAudioElement) {
    this.currentAudio = audio;
  }
}
