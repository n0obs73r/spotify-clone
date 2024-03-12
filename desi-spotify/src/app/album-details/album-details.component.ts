import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { ApiService } from '../services/api.service';
import { CommonModule } from '@angular/common';
import { Song } from '../models/song.model'; // Import the Song interface


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
  albumArt: string = ''; // Define albumArt property

  constructor(private route: ActivatedRoute, private apiService: ApiService) { }

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
  
        // console.log('Songs:', this.songs); // Log the entire songs array
  
        // Find the first song with artwork available
        const songWithArtwork = this.songs.find(song => song.artwork);
  
        // console.log('Song with artwork:', songWithArtwork); // Log the song with artwork
  
        // If a song with artwork is found, assign its artwork URL to albumArt
        if (songWithArtwork) {
          this.albumArt = 'data:image/jpeg;base64,' + songWithArtwork.artwork;
        } else {
          // Handle the case where no artwork is available
          this.albumArt = this.dummyImage;
        }
      },
      (error) => {
        console.error('Error loading songs:', error);
      }
    );
  }
  
  playSong(song: Song) {
    // const audio = new Audio();
    console.log(song)
    const audioBaseURL = "http://localhost:8080/songs/"
    const encodedFileName = encodeURIComponent(song.fileName).replace(/%26/g, '&');
    const audioURL = audioBaseURL + encodedFileName;
      const audio = new Audio(audioURL);
      audio.play();
    // console.log('Playing song:', audioBaseURL + encodeURIComponent(song.fileName).replace(/%20/g, ''));
}

}
