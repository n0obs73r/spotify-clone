import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { ApiService } from '../services/api.service';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-album',
  templateUrl: './albums.component.html',
  styleUrls: ['./albums.component.scss'],
  standalone:true,
  imports:[CommonModule]
})

export class AlbumComponent implements OnInit {
  albums: any[] = [];
  dummyImage: string = 'assets/images/dummy.jpg'; 

  constructor(private apiService: ApiService, private router: Router, private route: ActivatedRoute) { }

  ngOnInit(): void {
    this.route.queryParams.subscribe(params => {
      this.getAlbums();
    });
  }
  
  getAlbums() {
    this.apiService.getAlbums().subscribe((response: any[]) => {
      this.albums = response.map((album: any) => {
        let artworkSrc = album.artwork ? 'data:image/jpeg;base64,' + album.artwork : this.dummyImage;

        return { 
          title: album.title,
          artist: album.artist,
          tracks: album.tracks,
          artwork: artworkSrc,
        };
      });
    });
  }
  
  showAlbumDetails(album: any) {
    if (album.title) {
      const encodedTitle = encodeURIComponent(album.title);
      this.router.navigate(['/album-detail', encodedTitle]);
    } else {
      console.error('Album title is undefined');
    }
  }
}