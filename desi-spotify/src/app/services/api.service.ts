import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private apiUrl = 'http://localhost:8080'; // Your API URL here

  constructor(private http: HttpClient) { }

  getPlaylists(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/playlists`);
  }

  getAlbums() {
    return this.http.get<any[]>(`${this.apiUrl}/albums`);
  }

  getArtists(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/artists`);
  }

  getGenres(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/genres`);
  }

  search(query: string): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/search?query=${query}`);
  }

  getTracks(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/tracks`);
  }

  getAlbumArt(title: string): Observable<Blob> {
    return this.http.get(`${this.apiUrl}/albums/${title}/art`, { responseType: 'blob' });
  }
}
