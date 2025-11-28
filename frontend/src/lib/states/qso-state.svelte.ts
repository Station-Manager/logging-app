import { types } from '$lib/wailsjs/go/models';

/**
 * The QsoState interface represents the state of a QSO at the frontend UI. It is used to display information about the
 * current QSO to the user and accept information from the user. When the QSO is sent to be logged (stored on the database),
 * the information contained within this object will be used to create a QSO object sent to the backend.
 * Some information is derived from the backed when the QSO object is initially created, while others are user-entered and
 * still others derived from the CAT system (catState object).
 * NOTE: if information derived from the CAT system changes, the corresponding fields in the QsoState interface WILL
 * be updated accordingly. This happens in real-time as the user interacts with the Transceiver.
 */
export interface QsoState {
    original?: types.Qso; // Data derived from the backend when the QSO was created.
    callsign: string;
    rstSent: string;
    rstRcvd: string;
    mainmode: string;
    name: string;
    qth: string;
    comment: string;
    qsoDate: string;
    timeOn: string;
    timeOff: string;
    vfoaFreq: string;
    vfobFreq: string;
    vfoaBand: string;
    vfobBand: string;
}

export const qsoState: QsoState = $state({
    original: undefined,
    callsign: '',
    rstSent: '',
    rstRcvd: '',
    mainmode: '',
    name: '',
    qth: '',
    comment: '',
    qsoDate: '2025-11-28',
    timeOn: '06:57',
    timeOff: '',
    vfoaFreq: '14.320.000',
    vfobFreq: '14.320.000',
    vfoaBand: '20m',
    vfobBand: '20m',
});
