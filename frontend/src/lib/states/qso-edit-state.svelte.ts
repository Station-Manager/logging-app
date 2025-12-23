import { types } from '$lib/wailsjs/go/models';
import { formatDate, formatTime } from '$lib/utils/time-date';

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
    notes: string;
    rig: string;
    rx_pwr: string;
    tx_pwr: string;
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
    notes: '',
    rig: '',
    rx_pwr: '',
    tx_pwr: '',
    fromQso(this: QsoEditState, qso: types.Qso): void {
        this.call = qso.call;
        this.rst_sent = qso.rst_sent;
        this.rst_rcvd = qso.rst_rcvd;
        this.mode = qso.mode;
        this.submode = qso.submode;
        this.freq = qso.freq;
        this.freq_rx = qso.freq_rx;
        this.name = qso.name;
        this.qth = qso.qth;
        this.comment = qso.comment;
        this.notes = qso.notes;
        this.qso_date = formatDate(qso.qso_date);
        this.time_on = formatTime(qso.time_on);
        this.time_off = formatTime(qso.time_off);
        this.rig = qso.rig;
        this.rx_pwr = qso.rx_pwr;
    },
    toQso(this: QsoEditState): types.Qso {
        return new types.Qso({
            call: this.call,
            rst_sent: this.rst_sent,
            rst_rcvd: this.rst_rcvd,
            submode: this.submode,
            mode: this.mode,
            freq: this.freq,
            freq_rx: this.freq_rx,
            name: this.name,
            qth: this.qth,
            comment: this.comment,
            notes: this.notes,
            qso_date: this.qso_date,
            time_on: this.time_on,
            time_off: this.time_off,
            rig: this.rig,
            rx_pwr: this.rx_pwr,
        });
    },
});
