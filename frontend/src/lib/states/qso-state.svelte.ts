import { types } from '$lib/wailsjs/go/models';
import { getDateUTC, getTimeUTC } from '$lib/utils/time-date';

let elapsedIntervalID: number | null = null;

// Helper to reset a QsoState instance to its initial defaults for a new QSO.
// This centralizes defaults (e.g., RST values, initial mode, empty CAT fields)
// so they are applied consistently whenever a new QSO is started.
export function resetQsoStateDefaults(target: QsoState): void {
    // Core QSO defaults
    target.original = undefined;
    target.call = '';
    target.rst_sent = '59';
    target.rst_rcvd = '59';
    target.mode = 'USB';
    target.name = '';
    target.qth = '';
    target.comment = '';
    target.qso_date = getDateUTC();
    target.time_on = getTimeUTC();
    target.time_off = target.time_on;
    target.freq = '';
    target.freq_rx = '';
    target.band = '';
    target.band_rx = '';

    // CAT-only, UI-facing defaults (no CAT available yet)
    target.cat_identity = '';
    target.cat_vfoa_freq = '';
    target.cat_vfob_freq = '';
    target.cat_select = '';
    target.cat_split = '';
    target.cat_main_mode = '';
    target.cat_sub_mode = '';
    target.cat_tx_power = '';
}

// Helper to map backend QSO fields into the mutable QsoState instance.
// If you add more fields to QsoState later (e.g., submode, tx_pwr, etc.), you only need to update applyQsoToState.
// If some fields need transformation (e.g., formatting freq), you can put that logic into applyQsoToState and (if needed) the inverse into toQso
function applyQsoToState(target: QsoState, qso: types.Qso): void {
    target.call = qso.call ?? '';
    target.name = qso.name ?? '';
    target.qth = qso.qth ?? '';
    target.comment = qso.comment ?? '';

    // Preserve frontend defaults for RST if backend does not provide a value.
    if (qso.rst_sent != null && qso.rst_sent !== '') {
        target.rst_sent = qso.rst_sent;
    }
    if (qso.rst_rcvd != null && qso.rst_rcvd !== '') {
        target.rst_rcvd = qso.rst_rcvd;
    }

    if (qso.mode != null && qso.mode !== '') {
        target.mode = qso.mode;
    }

    // Preserve frontend defaults for date/time if backend does not provide a value.
    if (qso.qso_date != null && qso.qso_date !== '') {
        target.qso_date = getDateUTC();
    }
    //    if (qso.time_on != null && qso.time_on !== '') {
    target.time_on = getTimeUTC();
    //    }
    // time_off remains purely backend/explicit; empty means "not set" in UI.
    target.time_off = target.time_on;

    target.freq = qso.freq ?? '';
    target.freq_rx = qso.freq_rx ?? '';
    target.band = qso.band ?? '';
    target.band_rx = qso.band_rx ?? '';

    // NOTE: CAT-only fields are intentionally *not* populated from the backend QSO,
    // as they represent the *current rig state* rather than stored log data.
}

/**
 * Fields in QsoState that are driven directly from catState but do not map
 * 1:1 to ADIF / database schema.
 */
export interface CatDrivenFields {
    /** Current CAT identity / rig identifier. */
    cat_identity: string;
    /** Live VFO A frequency from CAT (formatted for display). */
    cat_vfoa_freq: string;
    /** Live VFO B frequency from CAT (formatted for display). */
    cat_vfob_freq: string;
    /** Which VFO is currently selected (e.g. A/B). */
    cat_select: string;
    /** Current split status (e.g. ON/OFF). */
    cat_split: string;
    /** Main mode as reported by CAT (may differ from ADIF-normalised mode). */
    cat_main_mode: string;
    /** Sub mode as reported by CAT. */
    cat_sub_mode: string;
    /** Transmit power reported by CAT. */
    cat_tx_power: string;
}

/**
 * QsoState is the UI-facing representation of the current QSO.
 *
 * Naming rule:
 * - Where possible, field names are identical to `types.Qso` (database schema)
 *   so that mapping in and out is mostly 1:1.
 * - Only genuinely UI-specific concerns (e.g. helper methods, CAT-only state) deviate.
 */
export interface QsoState extends CatDrivenFields {
    /** Full backend QSO snapshot (database schema aligned). */
    original?: types.Qso;

    /** ADIF / DB-aligned fields (subset exposed in the UI). */
    call: string;
    rst_sent: string;
    rst_rcvd: string;
    mode: string;
    name: string;
    qth: string;
    comment: string;
    qso_date: string; // YYYY-MM-DD
    time_on: string; // HH:mm[:ss]
    time_off: string; // HH:mm[:ss]
    freq: string;
    freq_rx: string;
    band: string;
    band_rx: string;

    /** Populate from backend QSO. */
    createFromQSO(this: QsoState, qso: types.Qso): void;
    /** Apply CAT-derived updates to QSO-related fields and live CAT-only fields. */
    updateFromCAT(this: QsoState, data: CatForQsoPayload): void;
    /** Reset editable fields back to the original QSO snapshot if present. */
    resetToOriginal(this: QsoState): void;
    /** Build a QSO instance suitable for sending back to the backend. */
    toQso(this: QsoState): types.Qso;

    startTimer(this: QsoState): void;
    stopTimer(this: QsoState): void;
    resetTimer(this: QsoState): void;
}

/**
 * The subset of QSO fields that are expected to be driven by CAT data
 * in the log-editing UI. This includes both ADIF-aligned fields and
 * CAT-only UI fields.
 */
export type CatForQsoPayload = Partial<
    Pick<
        QsoState,
        | 'freq'
        | 'freq_rx'
        | 'band'
        | 'band_rx'
        | 'mode'
        | 'cat_identity'
        | 'cat_vfoa_freq'
        | 'cat_vfob_freq'
        | 'cat_select'
        | 'cat_split'
        | 'cat_main_mode'
        | 'cat_sub_mode'
        | 'cat_tx_power'
    >
>;

export const qsoState: QsoState = $state({
    original: undefined,

    // Core QSO fields (aligned with types.Qso / DB schema)
    call: '',
    rst_sent: '59',
    rst_rcvd: '59',
    mode: 'USB',
    name: '',
    qth: '',
    comment: '',
    qso_date: '',
    time_on: '',
    time_off: '',
    freq: '',
    freq_rx: '',
    band: '',
    band_rx: '',

    // CAT-only, UI-facing fields (mirrors of `catState` for the current rig snapshot)
    cat_identity: '',
    cat_vfoa_freq: '',
    cat_vfob_freq: '',
    cat_select: '',
    cat_split: '',
    cat_main_mode: '',
    cat_sub_mode: '',
    cat_tx_power: '',

    createFromQSO(this: QsoState, qso: types.Qso) {
        if (!qso) return;

        // Keep a snapshot of the full backend model for later comparisons / reset.
        this.original = qso;

        applyQsoToState(this, qso);
    },

    updateFromCAT(this: QsoState, data: CatForQsoPayload): void {
        if (!data) return;

        const mappings: { [K in keyof CatForQsoPayload]: K } = {
            freq: 'freq',
            freq_rx: 'freq_rx',
            band: 'band',
            band_rx: 'band_rx',
            mode: 'mode',
            cat_identity: 'cat_identity',
            cat_vfoa_freq: 'cat_vfoa_freq',
            cat_vfob_freq: 'cat_vfob_freq',
            cat_select: 'cat_select',
            cat_split: 'cat_split',
            cat_main_mode: 'cat_main_mode',
            cat_sub_mode: 'cat_sub_mode',
            cat_tx_power: 'cat_tx_power',
        };

        (Object.entries(data) as Array<[keyof CatForQsoPayload, string]>).forEach(
            ([key, value]) => {
                const target = mappings[key];
                if (!target) return;
                this[target] = value;
            }
        );
    },

    resetToOriginal(this: QsoState): void {
        if (!this.original) {
            // No original QSO (new QSO case) -> reset to defaults.
            resetQsoStateDefaults(this);
            return;
        }
        applyQsoToState(this, this.original);
    },

    toQso(this: QsoState): types.Qso {
        // Start from original if present so we preserve all non-UI fields
        // (ids, upload flags, etc.). Otherwise, create a fresh Qso instance.
        const base = this.original ? types.Qso.createFrom(this.original) : new types.Qso({});

        base.call = this.call;
        base.name = this.name;
        base.qth = this.qth;
        base.comment = this.comment;

        base.rst_sent = this.rst_sent;
        base.rst_rcvd = this.rst_rcvd;

        base.mode = this.mode;

        base.qso_date = this.qso_date;
        base.time_on = this.time_on;
        base.time_off = this.time_off;

        base.freq = this.freq;
        base.freq_rx = this.freq_rx;
        base.band = this.band;
        base.band_rx = this.band_rx;

        // NOTE: CAT-only fields are intentionally *not* persisted back to the backend,
        // as they represent live rig state rather than stored QSO data.

        return base;
    },

    // QSO elapsed-time timer helpers
    // startTimer: begin a once-per-minute timer that updates time_off to the
    // current UTC time. Idempotent: if a timer is already running, it does
    // nothing.
    startTimer(this: QsoState): void {
        // If a timer is already active, don't start another.
        if (elapsedIntervalID !== null) {
            return;
        }

        // Ensure we have an initial end time; if it's empty, initialise it
        // to "now" so UI has an immediate value.
        if (!this.time_off) {
            this.time_off = getTimeUTC();
        }

        // Store interval id in module-scope variable so we can reliably clear it.
        elapsedIntervalID = window.setInterval(() => {
            // Always write through the shared qsoState instance to avoid any
            // confusion around `this` binding inside the interval callback.
            qsoState.time_off = getTimeUTC();
        }, 60_000); // every minute
    },

    // stopTimer: stop any running elapsed-time timer but keep the last
    // recorded time_off value intact. Safe to call multiple times.
    stopTimer(this: QsoState): void {
        if (elapsedIntervalID !== null) {
            clearInterval(elapsedIntervalID);
            elapsedIntervalID = null;
        }
    },

    // resetTimer: stop any running timer and reset the timing-related fields
    // for a fresh QSO context.
    resetTimer(this: QsoState): void {
        // Ensure no interval continues running.
        this.stopTimer();

        const date = getDateUTC();
        const time = getTimeUTC();

        this.qso_date = date;
        this.time_on = time;
        this.time_off = time;
    },
});

// Immediately initialize defaults for the shared qsoState instance.
resetQsoStateDefaults(qsoState);
