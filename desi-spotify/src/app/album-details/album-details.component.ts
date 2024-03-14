import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { ApiService } from '../services/api.service';
import { CommonModule } from '@angular/common';
import { Song } from '../models/song.model';
import {PlayerService} from "../services/player.service"; // Import the Song interface


@Component({
  selector: 'app-album-details',
  standalone: true,
  imports:[CommonModule],
  templateUrl: './album-details.component.html',
  styleUrl: './album-details.component.scss'
})

export class AlbumDetailsComponent implements OnInit {
  albumTitle: string = '';
  songs: any[] = [];
  dummyImage: string = 'assets/images/dummy.jpg';
  albumArt: string = '';

  constructor(private route: ActivatedRoute, private apiService: ApiService,
              private playerService: PlayerService
  ) {
  }

  ngOnInit(): void {
    this.route.params.subscribe(params => {
      this.albumTitle = params['title'];
      this.getAlbumSongs();
    });
  }

  getAlbumSongs() {
    this.apiService.getAlbumSongs(this.albumTitle).subscribe(
      (response: any[]) => {
        this.songs = response;
        const songWithArtwork = this.songs.find(song => song.artwork);
        if (songWithArtwork) {
          this.albumArt = 'data:image/jpeg;base64,' + songWithArtwork.artwork;
        } else {
          this.albumArt = this.dummyImage;
        }
      },
      (error) => {
        console.error('Error loading songs:', error);
      }
    );
  }

  playSong(song: Song) {
    this.playerService.playSongService(song);
  }
}
