import { types } from '$lib/wailsjs/go/models';

export interface QsoEditState {
    call: string;
    freq: string;
    freq_rx: string;
    qso_date: string;
    time_on: string;
    time_off: string;
    rst_sent: string;
    rst_rcvd: string;
    mode: string;
    submode: string;
    name: string;
    qth: string;
    comment: string;
    fromQso(this: QsoEditState, qso: types.Qso): void;
    toQso(this: QsoEditState): types.Qso;
}

export const qsoEditState: QsoEditState = $state({
    call: '',
    freq: '',
    freq_rx: '',
    qso_date: '',
    time_on: '',
    time_off: '',
    rst_sent: '',
    rst_rcvd: '',
    mode: '',
    submode: '',
    name: '',
    qth: '',
    comment: '',
    fromQso(this: QsoEditState, qso: types.Qso): void {
        this.call = qso.call;
        this.rst_sent = qso.rst_sent;
        this.rst_rcvd = qso.rst_rcvd;
        this.submode = qso.submode;
        this.freq = qso.freq;
        this.freq_rx = qso.freq_rx;
        this.name = qso.name;
        this.qth = qso.qth;
        this.comment = qso.comment;
        this.qso_date = qso.qso_date;
        this.time_on = qso.time_on;
        this.time_off = qso.time_off;
    },
    toQso(this: QsoEditState): types.Qso {
        return new types.Qso({});
    },
});
