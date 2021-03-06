import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ModalController } from '@ionic/angular';

import { Ext } from './ext';
import { Config } from '../api/config'
import { ExtPage } from './ext.page';

@Injectable({
  providedIn: 'root'
})
export class ExtService {
  constructor(
    private modalCtrl: ModalController,
    private http: HttpClient
  ) {
  }

  get extsByCustomer$() {
    return this.http.get<Ext[]>(`${Config.URL}/api/exts/E-C`)
  }
  get extsByItem$() {
    return this.http.get<Ext[]>(`${Config.URL}/api/exts/E-I`)
  }

  async saveByCustomer(exts: Ext[]) {
    return this.http.put<Ext[]>(`${Config.URL}/api/exts/E-C`, exts).toPromise()
  }

  async saveByItem(exts: Ext[]) {
    return this.http.put<Ext[]>(`${Config.URL}/api/exts/E-I`, exts).toPromise()
  }

  async extCustomer() {
    const modal = await this.modalCtrl.create({
      component: ExtPage,
      componentProps: {
        exts: await this.extsByCustomer$.toPromise(),
        service: this,
        code: 'C',
      }
    })
    await modal.present()
    return await modal.onDidDismiss()
  }

  async extItem() {
    const modal = await this.modalCtrl.create({
      component: ExtPage,
      componentProps: {
        exts: await this.extsByItem$.toPromise(),
        service: this,
        code: 'I',
      }
    })
    return await modal.present()
  }
}
