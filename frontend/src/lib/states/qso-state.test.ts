import { describe, it, expect, beforeEach } from 'vitest';
import { dottedKHzToShortMHz } from '$lib/utils/frequency';
import { types } from '$lib/wailsjs/go/models';
import { qsoState, type CatForQsoPayload } from './qso-state.svelte';

function resetQsoState() {
    qsoState.original = undefined;
    qsoState.call = '';
    qsoState.name = '';
    qsoState.qth = '';
    qsoState.comment = '';

    qsoState.rst_sent = '';
    qsoState.rst_rcvd = '';

    qsoState.mode = '';

    qsoState.qso_date = '';
    qsoState.time_on = '';
    qsoState.time_off = '';

    qsoState.freq = '';
    qsoState.freq_rx = '';
    qsoState.band = '';
    qsoState.band_rx = '';
}

describe('qso-state', () => {
    beforeEach(() => {
        resetQsoState();
    });

    it('creates state from backend QSO and resets correctly', () => {
        const backendQso = types.Qso.createFrom({
            call: 'K1ABC',
            name: 'Alice',
            qth: 'Boston',
            comment: 'Test QSO',
            rst_sent: '59',
            rst_rcvd: '59',
            mode: 'SSB',
            qso_date: '2025-11-28',
            time_on: '10:00',
            time_off: '10:05',
            freq: '14.250',
            freq_rx: '14.250',
            band: '20m',
            band_rx: '20m',
        });

        qsoState.createFromQSO(backendQso);

        expect(qsoState.call).toBe('K1ABC');
        expect(qsoState.mode).toBe('SSB');
        // qso_date is currently normalised to current UTC date in applyQsoToState,
        // so we just assert it is non-empty rather than matching the backend value.
        expect(qsoState.qso_date).not.toBe('');
        expect(qsoState.freq).toBe('14.250');

        // Mutate some fields to simulate user edits.
        qsoState.call = 'K1XYZ';
        qsoState.comment = 'Changed comment';

        // Reset should bring us back to the original backend values for simple fields.
        qsoState.resetToOriginal();

        expect(qsoState.call).toBe('K1ABC');
        expect(qsoState.comment).toBe('Test QSO');
    });

    it('updates from CAT payload for whitelisted fields', () => {
        // Start with a known baseline.
        qsoState.freq = '7.000';
        qsoState.freq_rx = '7.000';
        qsoState.band = '40m';
        qsoState.band_rx = '40m';
        qsoState.mode = 'CW';

        const catPayload: CatForQsoPayload = {
            freq: '14.320',
            freq_rx: '14.320',
            band: '20m',
            band_rx: '20m',
            mode: 'SSB',
        };

        qsoState.updateFromCAT(catPayload);

        expect(qsoState.freq).toBe('14.320');
        expect(qsoState.freq_rx).toBe('14.320');
        expect(qsoState.band).toBe('20m');
        expect(qsoState.band_rx).toBe('20m');
        expect(qsoState.mode).toBe('SSB');
    });

    it('ignores unknown keys in CAT payload', () => {
        qsoState.freq = '7.000';

        // Unknown key should be ignored and not affect freq.
        type ExtendedCatPayload = CatForQsoPayload & { foo: string };
        const payload: ExtendedCatPayload = { foo: 'bar' };

        qsoState.updateFromCAT(payload);

        expect(qsoState.freq).toBe('7.000');
    });

    it('round-trips toQso preserving non-UI fields', () => {
        const original = types.Qso.createFrom({
            id: 123,
            logbook_id: 1,
            call: 'G1AAA',
            mode: 'CW',
            qso_date: '2025-11-27',
            time_on: '09:00',
            time_off: '09:05',
            freq: '7.010',
            band: '40m',
        });

        qsoState.createFromQSO(original);

        // Simulate edits.
        qsoState.call = 'G1BBB';
        qsoState.mode = 'SSB';
        qsoState.freq = '14.200';
        qsoState.band = '20m';

        const built = qsoState.toQso();

        expect(built.call).toBe('G1BBB');
        expect(built.mode).toBe('SSB');
        expect(built.freq).toBe(dottedKHzToShortMHz(qsoState.cat_vfoa_freq));
        expect(built.band).toBe('20m');

        // Ensure we still preserve non-UI fields like id/logbook_id from original.
        expect(built.id).toBe(original.id);
        expect(built.logbook_id).toBe(original.logbook_id);
    });

    it('startTimer initialises and is idempotent, stopTimer clears interval, resetTimer resets timing fields', async () => {
        // Capture original values so we can detect changes.
        const initialDate = qsoState.qso_date;
        const initialOn = qsoState.time_on;

        // Calling resetTimer should set date/on/off to non-empty and off == on.
        qsoState.resetTimer();
        expect(qsoState.qso_date).not.toBe('');
        expect(qsoState.time_on).not.toBe('');
        expect(qsoState.time_off).toBe(qsoState.time_on);

        const firstOff = qsoState.time_off;

        // startTimer should be safe to call multiple times without throwing.
        qsoState.startTimer();
        qsoState.startTimer();

        // Stop the timer immediately to avoid relying on real time passing.
        qsoState.stopTimer();

        // After stopTimer, time_off should remain whatever it was.
        expect(qsoState.time_off).toBe(firstOff);

        // reset again should change timing fields.
        qsoState.resetTimer();
        expect(qsoState.time_on).not.toBe(initialOn);
        expect(qsoState.qso_date).not.toBe(initialDate);
        expect(qsoState.time_off).toBe(qsoState.time_on);
    });
});
