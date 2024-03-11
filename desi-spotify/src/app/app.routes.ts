import { Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';
// import { AlbumComponent } from './albums/albums.component';
import { AlbumDetailComponent } from './album-details/album-details.component';

export const routes: Routes = [
    { path: 'home', component: HomeComponent },
    { path: 'albums', loadComponent: () => import('./albums/albums.component').then(mod => mod.AlbumComponent) },
    { path: 'albums/:title/songs', component: AlbumDetailComponent }, // Route parameter :title for album title

];
