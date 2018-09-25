import { Injectable, NgZone } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { ObjService } from '../api/obj.service';
import { Item } from './item';

@Injectable({
  providedIn: 'root'
})
export class ItemService extends ObjService<Item>{
  constructor(
    protected http: HttpClient,
    protected ngZone: NgZone,
  ) {
    super(http, ngZone)
  }

  path() { return '/api/items' }
}
