import { configState } from '$lib/states/config-state.svelte';
import { types } from '$lib/wailsjs/go/models';
import { formatCatKHzToDottedMHz } from '$lib/utils/frequency';
import { getDateUTC, getTimeUTC } from '$lib/utils/time-date';

const CAT_MAPPINGS: { [K in keyof CatForQsoPayload]: K } = {
    cat_vfoa_freq: 'cat_vfoa_freq',
    cat_vfob_freq: 'cat_vfob_freq',
    cat_main_mode: 'cat_main_mode',
};

export interface QsoTimerState {
    elapsed: number;
    running: boolean;
}

export const qsoTimerState: QsoTimerState = $state({ elapsed: 0, running: false });

let elapsedIntervalID: number | null = null;

export type CatForQsoPayload = Partial<
    Pick<QsoState, 'cat_vfoa_freq' | 'cat_vfob_freq' | 'cat_main_mode'>
>;

export interface CatDrivenFields {
    cat_vfoa_freq: string;
    cat_vfob_freq: string;
    cat_main_mode: string;
}

export interface QsoState extends CatDrivenFields {
    original?: types.Qso;
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
    original: undefined,
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
        rstHelper(this);
        randomQsoHelper(this);
        this.call = '';
        this.name = '';
        this.qth = '';
        this.comment = '';
        this.notes = '';
        this.qso_date = getDateUTC();
        this.time_on = getTimeUTC();
        this.time_off = getTimeUTC();
    },
    fromQso(this: QsoState, qso: types.Qso): void {
        if (!qso) return;
        this.original = qso;
        this.call = qso.call;
        this.name = qso.name ?? '';
        this.qth = qso.qth ?? '';
        rstHelper(this);
        randomQsoHelper(this);
    },
    toQso(this: QsoState): types.Qso {
        const base = this.original ? types.Qso.createFrom(this.original) : new types.Qso({});
        base.call = this.call;
        base.name = this.name;
        base.qth = this.qth;
        base.comment = this.comment;
        base.notes = this.notes;

        base.rst_sent = this.rst_sent;
        base.rst_rcvd = this.rst_rcvd;

        base.mode = this.mode;

        base.qso_date = this.qso_date;
        base.time_on = this.time_on;
        base.time_off = this.time_off;

        return new types.Qso({});
    },
    updateFromCAT(this: QsoState, data: CatForQsoPayload): void {
        if (!data) return;

        (Object.entries(data) as Array<[keyof CatForQsoPayload, string | undefined]>).forEach(
            ([key, value]) => {
                if (value === undefined) return;
                const catKey = CAT_MAPPINGS[key];
                if (!catKey) return;
                switch (catKey) {
                    case 'cat_vfoa_freq':
                    case 'cat_vfob_freq':
                        value = formatCatKHzToDottedMHz(value);
                        break;
                }
                this[catKey] = value;
            }
        );

        this.cat_enabled = true;
    },
    startTimer(this: QsoState): void {
        // If a timer is already active, don't start another.
        if (elapsedIntervalID !== null) {
            return;
        }

        // Ensure we have an initial end time; if it's empty, initialize it
        // to "now" so UI has an immediate value.
        if (!this.time_off) {
            this.time_off = getTimeUTC();
        }

        // Store interval id in a module-scope variable so we can reliably clear it.
        elapsedIntervalID = window.setInterval(() => {
            // Always write through the shared qsoState instance to avoid any
            // confusion around `this` binding inside the interval callback.
            qsoState.time_off = getTimeUTC();
        }, 60_000); // every minute

        // Mark the timer as running so subscribers know.
        qsoTimerState.running = true;
    },
    stopTimer(this: QsoState): void {
        if (elapsedIntervalID !== null) {
            clearInterval(elapsedIntervalID);
            elapsedIntervalID = null;
        }
        qsoTimerState.running = false;
    },
    resetTimer(this: QsoState): void {
        // Ensure no interval continues running.
        this.stopTimer();

        const date = getDateUTC();
        const time = getTimeUTC();

        this.qso_date = date;
        this.time_on = time;
        this.time_off = time;
    },
    isTimerRunning(this: QsoState): boolean {
        return elapsedIntervalID !== null;
    },
});

const rstHelper = (qso: QsoState): void => {
    if (!qso) return;
    if (qso.cat_main_mode === 'CW-U' || qso.cat_main_mode === 'CW-L') {
        qso.rst_rcvd = '599';
        qso.rst_sent = '599';
    } else {
        qso.rst_rcvd = '59';
        qso.rst_sent = '59';
    }
};

const randomQsoHelper = (qso: QsoState): void => {
    if (!qso) return;
    qso.qso_random = configState.default_random_qso ? 'Y' : 'N';
};
