import { Injectable } from "@angular/core";
import { BehaviorSubject } from "rxjs";
import { Song } from "../models/song.model";
import { Howl } from 'howler';

@Injectable({
  providedIn: 'root'
})
export class PlayerService {
  playing$ = new BehaviorSubject<boolean>(false);
  private currentHowl: Howl | null = null;
  private currentSongUrl: string | null = null;
  private playlist: Song[] = [];
  private currentIndex: number = 0;

  constructor() {
    this.playing$.subscribe(playing => {
      if (this.currentHowl) {
        if (playing) {
          this.currentHowl.play();
        } else {
          this.currentHowl.pause();
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

  setCurrentHowl(howl: Howl) {
    this.currentHowl = howl;
  }

  setPlaylist(playlist: Song[]) {
    this.playlist = playlist;
  }

  playSongService(song: Song) {
    const songUrl = this.getAudioUrl(song);
    console.log('Attempting to play song with URL:', songUrl);

    // Check if a song is currently playing
    if (this.currentHowl) {
        // If a song is playing, stop it before playing the new song
        this.stopCurrentSong();
    }
    
    const howl = new Howl({
        src: songUrl ? [songUrl] : [''],
        onload: () => {
            console.log('Audio file loaded successfully.');
            console.log('Is current Howl loaded?', howl.state() === 'loaded');
            this.setPlaying(true);
            howl.play();
            console.log('Song started playing.');
        },
        onloaderror: (id, error) => {
            console.error('Error loading audio:', error);
            this.stopCurrentSong();
        },
        onplayerror: (id, error) => {
            console.error('Error playing audio:', error);
            this.stopCurrentSong();
        }
    });

    console.log('Howl instance:', howl);
    console.log('Current Howl error:', howl.state() === 'unloaded' ? 'No error' : 'Error');

    this.setCurrentHowl(howl);
    this.currentSongUrl = songUrl;
}



  stopCurrentSong() {
    if (this.currentHowl) {
      this.currentHowl.stop();
      this.currentHowl.unload();
      this.currentHowl = null;
      this.currentSongUrl = null;
      this.setPlaying(false);
      console.log('Song stopped.');
    }
  }

  private getAudioUrl(song: Song): string | null {
    if (!song || !song.fileName) return null;
    const audioBaseURL = "http://localhost:8080/songs/";
    const encodedFileName = encodeURIComponent(song.fileName).replace(/%26/g, '&');
    const timestamp = new Date().getTime();
    return `${audioBaseURL}${encodedFileName}?t=${timestamp}`;
  }

  playNextSong() {
    this.currentIndex = (this.currentIndex + 1) % this.playlist.length;
    this.loadCurrentSong();
  }

  playPreviousSong() {
    this.currentIndex = (this.currentIndex - 1 + this.playlist.length) % this.playlist.length;
    this.loadCurrentSong();
  }

  private loadCurrentSong() {
    const song = this.playlist[this.currentIndex];
    this.playSongService(song);
  }
}
