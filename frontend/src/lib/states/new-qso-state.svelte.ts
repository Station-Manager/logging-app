import {configState} from '$lib/states/config-state.svelte';
import {types} from '$lib/wailsjs/go/models';
import {formatCatKHzToDottedMHz, formatDottedKHzToDottedMHz} from '$lib/utils/frequency';

const CAT_MAPPINGS: { [K in keyof CatForQsoPayload]: K } = {
    cat_vfoa_freq: 'cat_vfoa_freq',
    cat_vfob_freq: 'cat_vfob_freq',
};

export interface QsoTimerState {
    elapsed: number;
    running: boolean;
}

export const qsoTimerState: QsoTimerState = $state({elapsed: 0, running: false});

export type CatForQsoPayload = Partial<Pick<QsoState, 'cat_vfoa_freq' | 'cat_vfob_freq'>>;

export interface CatDrivenFields {
    cat_vfoa_freq: string;
    cat_vfob_freq: string;
}

export interface QsoState extends CatDrivenFields {
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
    rx_pwr: string;
    tx_pwr: string;
    notes: string;
    cqz: string;
    rig: string;
    qso_random: string;

    country_name: string;
    ccode: string;
    ant_path: string;
    short_path_distance: string;
    long_path_distance: string;
    remote_time: string;
    remote_offset: string;

    contact_history: types.ContactHistory[];

    web: string;
    email: string;

    cat_enabled: boolean;
    cat_vfoa_freq: string;
    cat_vfob_freq: string;
    cat_main_mode: string;

    reset(this: QsoState): void;

    fromQso(this: QsoState, qso: types.Qso): void;

    toQso(this: QsoState): types.Qso;

    updateFromCAT(this: QsoState, data: CatForQsoPayload): void;

    startTimer(this: QsoState): void;

    stopTimer(this: QsoState): void;

    resetTimer(this: QsoState): void;

    isTimerRunning(this: QsoState): boolean;
}

export const qsoState: QsoState = $state({
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
    rx_pwr: '',
    tx_pwr: '',
    notes: '',
    cqz: '',
    web: '',
    email: '',
    rig: '',
    qso_random: '',

    country_name: '',
    ccode: '',
    ant_path: '',
    short_path_distance: '',
    long_path_distance: '',
    remote_time: '',
    remote_offset: '',

    contact_history: [],

    cat_enabled: false,
    cat_vfoa_freq: '',
    cat_vfob_freq: '',
    cat_main_mode: '',

    reset(this: QsoState): void {
        if (!this.cat_enabled) {
            // These checks allow us to maintain user entered values when CAT is disabled
            // as CAT values are used directly in the UI.
            if (this.cat_vfoa_freq === '') {
                this.cat_vfoa_freq = configState.default_freq;
            }
            if (this.cat_main_mode === '') {
                this.cat_main_mode = configState.default_mode;
            }
        }
        if (this.cat_main_mode === 'CW-U' || this.cat_main_mode === 'CW-L') {
            this.rst_rcvd = '599';
            this.rst_sent = '599';
        } else {
            this.rst_rcvd = '59';
            this.rst_sent = '59';
        }
    },
    fromQso(this: QsoState, qso: types.Qso): void {
        console.log('fromQso', qso);
    },
    toQso(this: QsoState): types.Qso {
        console.log('toQso');
        return new types.Qso({});
    },
    updateFromCAT(this: QsoState, data: CatForQsoPayload): void {
        console.log('updateFromCAT', data);
        if (!data) return;

        (Object.entries(data) as Array<[keyof CatForQsoPayload, string | undefined]>).forEach(
            ([key, value]) => {
                if (value === undefined) return;
                const catKey = CAT_MAPPINGS[key];
                if (!catKey) return;
                switch (catKey) {
                    case 'cat_vfoa_freq':
                    case 'cat_vfob_freq':
                        value = value.includes('.')
                            ? formatDottedKHzToDottedMHz(value)
                            : formatCatKHzToDottedMHz(value);
                        break;
                }
                this[catKey] = value;
                console.log('>>', value);
            }
        );

        this.cat_enabled = true;
    },
    startTimer(this: QsoState): void {
        console.log('startTimer');
    },
    stopTimer(this: QsoState): void {
        console.log('stopTimer');
    },
    resetTimer(this: QsoState): void {
        console.log('resetTimer');
    },
    isTimerRunning(this: QsoState): boolean {
        return false;
    },
});
