import { types } from '$lib/wailsjs/go/models';

export interface QsoState {
    original?: types.Qso;
    callsign: string;
    rstSent: string;
    rstRcvd: string;
    mainmode: string;
}

export const qsoState: QsoState = $state({
    original: undefined,
    callsign: '',
    rstSent: '',
    rstRcvd: '',
    mainmode: '',
});
