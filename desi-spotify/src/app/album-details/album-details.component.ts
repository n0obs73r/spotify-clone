import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { ApiService } from '../services/api.service';
import { CommonModule } from '@angular/common';

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
      },
      (error) => {
        console.error('Error loading songs:', error);
      }
    );
  }
}