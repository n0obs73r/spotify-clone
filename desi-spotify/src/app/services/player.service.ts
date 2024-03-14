import { Injectable } from "@angular/core";
import { BehaviorSubject, Subject } from "rxjs";
import { Song } from "../models/song.model";
import { Howl, Howler } from 'howler';

@Injectable({
  providedIn: 'root'
})
export class PlayerService {
  playing$ = new BehaviorSubject<boolean>(false);
  private currentHowl: Howl | null = null;
  private currentSongUrl: string | null = null;
  private playlist: Song[] = [];
  private currentIndex: number = 0;
  private progress$ = new Subject<number>(); // Subject to emit playback progress
  private duration: number = 0;
  private progressInterval: any;

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
  
    // Reset duration when loading a new song
    this.duration = 0;
  
    // Emit an event indicating that a new song is loading
    this.setPlaying(false);
  
    // Check if the song to play is the same as the currently playing song
    if (this.currentSongUrl === songUrl && this.currentHowl) {
      // If the song is the same and there's an existing Howl instance, just seek to the beginning
      this.currentHowl.seek(0);
      this.setPlaying(true); // Ensure the player is set to playing
      return;
    }
  
    // If a different song is requested or there's no existing Howl instance, create a new one
    if (this.currentHowl) {
      // If a song is playing and it's not the requested song, stop it before playing the new song
      this.stopCurrentSong();
    }
  
    const howl = new Howl({
      src: songUrl ? [songUrl] : [''],
      onload: () => {
        console.log('Audio file loaded successfully.');
        console.log('Is current Howl loaded?', howl.state() === 'loaded');
        this.setPlaying(true);
        this.duration = howl.duration(); // Update duration when the song is loaded
        console.log('Song started playing.');
        // Start updating progress when the song starts playing
        this.startProgressInterval();
      },
      onloaderror: (id, error) => {
        console.error('Error loading audio:', error);
        this.stopCurrentSong();
      },
      onplayerror: (id, error) => {
        console.error('Error playing audio:', error);
        this.stopCurrentSong();
      },
      onend: () => { // When the song ends, reset progress
        this.progress$.next(0);
        this.setPlaying(false); // Ensure the player is set to not playing when the song ends
        this.resetProgressInterval();
      },
      onseek: (position) => { // When the song is seeked, update progress
        this.progress$.next(position);
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
      this.resetProgressInterval();
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

  getProgress() {
    return this.progress$.asObservable();
  }

  getTotalDuration() {
    return this.duration;
  }
  
  seekTo(position: number) {
    if (this.currentHowl) {
      this.currentHowl.seek(position);
      if (!this.playing$.getValue()) {
        // If the player is not playing, start the progress interval after seeking
        this.startProgressInterval();
      }
    }
  }
  
  

  private startProgressInterval() {
    this.progressInterval = setInterval(() => {
      const seek = this.currentHowl?.seek() || 0;
      this.progress$.next(seek);
    }, 1000);
  }

  private resetProgressInterval() {
    clearInterval(this.progressInterval);
  }
}
