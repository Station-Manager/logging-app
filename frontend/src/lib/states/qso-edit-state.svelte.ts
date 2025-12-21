import {types} from "$lib/wailsjs/go/models";

export interface QsoEditState {
    call: string;
    qso_date: string;
    time_on: string;
    time_off: string;
    rst_sent: string;
    rst_rcvd: string;
    mode: string;
    name: string;
    qth: string;
    comment: string;
    fromQso(this: QsoEditState, qso: types.Qso): void;
    toQso(this: QsoEditState): types.Qso;
}

export const qsoEditState: QsoEditState = $state({
    call: '',
    qso_date: '',
    time_on: '',
    time_off: '',
    rst_sent: '',
    rst_rcvd: '',
    mode: '',
    name: '',
    qth: '',
    comment: '',
    fromQso(this: QsoEditState, qso: types.Qso): void {

    },
    toQso(this: QsoEditState): types.Qso {
        return new types.Qso({});
    }
});
