import { describe, it, expect, beforeEach } from 'vitest';
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
        expect(qsoState.qso_date).toBe('2025-11-28');
        expect(qsoState.freq).toBe('14.250');

        // Mutate some fields to simulate user edits.
        qsoState.call = 'K1XYZ';
        qsoState.comment = 'Changed comment';

        // Reset should bring us back to the original backend values.
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
        expect(built.freq).toBe('14.200');
        expect(built.band).toBe('20m');

        // Ensure we still preserve non-UI fields like id/logbook_id from original.
        expect(built.id).toBe(original.id);
        expect(built.logbook_id).toBe(original.logbook_id);
    });
});
