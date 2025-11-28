import { types } from '$lib/wailsjs/go/models';

// Helper to map backend QSO fields into the mutable QsoState instance.
// If you add more fields to QsoState later (e.g., submode, tx_pwr, etc.), you only need to update applyQsoToState.
// If some fields need transformation (e.g., formatting freq), you can put that logic into applyQsoToState and (if needed) the inverse into toQso
function applyQsoToState(target: QsoState, qso: types.Qso): void {
    target.call = qso.call ?? '';
    target.name = qso.name ?? '';
    target.qth = qso.qth ?? '';
    target.comment = qso.comment ?? '';

    target.rst_sent = qso.rst_sent ?? '';
    target.rst_rcvd = qso.rst_rcvd ?? '';

    target.mode = qso.mode ?? '';

    target.qso_date = qso.qso_date ?? '';
    target.time_on = qso.time_on ?? '';
    target.time_off = qso.time_off ?? '';

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
        if (!this.original) return;
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
});
