import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-album-detail',
  templateUrl: './album-detail.component.html',
  styleUrls: ['./album-detail.component.scss']
})
export class AlbumDetailComponent implements OnInit {
  albumId!: string;
  albumDetails: any;

  constructor(private route: ActivatedRoute, private http: HttpClient) { }

  ngOnInit(): void {
    this.route.params.subscribe(params => {
      this.albumId = params['id'];
      this.fetchAlbumDetails();
    });
  }

  fetchAlbumDetails() {
    // Make HTTP request to fetch album details based on albumId
    this.http.get<any>(`/albums/${this.albumId}`).subscribe(response => {
      this.albumDetails = response;
    });
  }
}
