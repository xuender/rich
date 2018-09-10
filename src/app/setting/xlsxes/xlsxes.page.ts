import { Component, OnInit } from '@angular/core';
import { XlsxService } from '../../api/xlsx.service';
import { ModalController } from '@ionic/angular';
import { XlsxPage } from './xlsx/xlsx.page';
import { Xlsx } from '../../api/xlsx';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-xlsxes',
  templateUrl: './xlsxes.page.html',
  styleUrls: ['./xlsxes.page.scss'],
})
export class XlsxesPage implements OnInit {
  xlsxes: Observable<Xlsx[]>
  constructor(
    private xlsx: XlsxService,
    private modalCtrl: ModalController,
  ) {
    this.xlsxes = xlsx.xlsesx
  }

  cancel() {
    this.modalCtrl.dismiss();
  }
  ngOnInit() {
  }

  async create() {
    const modal = await this.modalCtrl.create({
      component: XlsxPage,
      componentProps: {
        xlsx: {
          name: '新Excel定义',
          map: {},
        }
      }
    });
    modal.onDidDismiss().then(d => {
      if (d.data) {
        this.xlsx.post(d.data)
          .subscribe(() => {
            this.xlsxes = this.xlsx.xlsesx
          })
      }
    })
    return await modal.present()
  }
  async update(x: Xlsx) {
    const modal = await this.modalCtrl.create({
      component: XlsxPage,
      componentProps: { xlsx: Object.assign({}, x) },
    });
    modal.onDidDismiss().then(d => {
      if (d.data) {
        this.xlsx.put(d.data)
          .subscribe(nx => {
            Object.assign(x, nx)
          })
      }
    })
    return await modal.present()
  }
}