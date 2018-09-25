import { NgZone } from '@angular/core';
import { HttpClient } from "@angular/common/http";
import { find, pull } from 'lodash'

import { Obj } from './obj';
import { URL } from './init'
import { Paging } from './paging';

export abstract class ObjService<T extends Obj> {
  private datas: T[] = []
  private txt: string = ''
  private tags: string[] = []
  private page = 0
  constructor(
    protected http: HttpClient,
    protected ngZone: NgZone,
  ) { }

  abstract path(): string

  get all$() {
    return this.http.get<T[]>(`${URL}${this.path()}`, { params: this.params })
  }

  get paging$() {
    return this.http.get<Paging<T>>(`${URL}${this.path()}`, { params: this.params })
  }

  get nextPaging$() {
    this.page += 1
    return this.paging$
  }

  has(datas: T[]): boolean {
    if (datas) {
      this.datas = datas
      return datas.length > 0
    }
    return false;
  }

  async save(o: T) {
    if (o.id) {
      const t = await this.http.put<T>(`${URL}${this.path()}/${o.id}`, o).toPromise()
      Object.assign(find(this.datas, (d) => d.id == t.id), t)
      return Object.assign(o, t)
    } else {
      const t = await this.http.post<T>(`${URL}${this.path()}`, o).toPromise()
      this.datas.push(t)
      return t
    }
  }

  async delete(o: T) {
    await this.http.delete(`${URL}${this.path()}/${o.id}`).toPromise()
    pull(this.datas, o)
    this.ngZone.run(() => false)
  }

  search(txt: string) {
    this.txt = txt
    return this.all$
  }

  select(tags: string[]) {
    this.tags = tags
    return this.all$
  }

  searchPaging(txt: string) {
    this.page = 0
    this.txt = txt
    return this.paging$
  }

  selectPaging(tags: string[]) {
    this.page = 0
    this.tags = tags
    return this.paging$
  }


  private get params() {
    const params = {}
    if (this.txt) {
      params['search'] = this.txt
    }
    if (this.tags && this.tags.length > 0) {
      params['tags'] = JSON.stringify(this.tags)
    }
    if (this.page > 0) {
      params['page'] = this.page
    }
    return params
  }
}