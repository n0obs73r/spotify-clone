import { Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { AlbumComponent } from './albums/albums.component';
import { AlbumDetailComponent } from './album-detail/album-detail.component';

export const routes: Routes = [
    { path: 'home', component: HomeComponent },
    { path: 'albums', component: AlbumComponent },
    { path: 'album-detail/:id', component: AlbumDetailComponent },
];
