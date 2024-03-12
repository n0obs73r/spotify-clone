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
      const uniqueAlbumsMap = new Map<string, any>();
      
      // Group albums by title
      response.forEach((album: any) => {
        const key = album.title;
        if (!uniqueAlbumsMap.has(key)) {
          uniqueAlbumsMap.set(key, {
            title: album.title,
            tracks: [],
            artwork: album.artwork ? 'data:image/jpeg;base64,' + album.artwork : this.dummyImage
          });
        }
        
        const uniqueAlbum = uniqueAlbumsMap.get(key);
        uniqueAlbum.tracks.push(album.tracks);
      });
      
      this.albums = Array.from(uniqueAlbumsMap.values());
    });
  }
  
  showAlbumDetails(album: any) {
    if (album.title) {
      const titleWithSpaces = decodeURIComponent(album.title); // Decode URI component to replace %20 with spaces
      this.router.navigate(['/albums', titleWithSpaces, 'songs']); // Navigate to the desired route
    } else {
      console.error('Album title is undefined');
    }
  }
}