import { describe, it, expect } from 'vitest';
import {
    parseCatKHzToMHz,
    parseUserMHz,
    frequencyToBandMHz,
    frequencyToBandFromCat,
    frequencyToBandFromMHz,
    frequencyToBandFromDottedMHz,
} from './frequency';

describe('frequency utilities', () => {
    describe('parseCatKHzToMHz', () => {
        it('parses valid 9-digit CAT kHz string to MHz', () => {
            expect(parseCatKHzToMHz('007101000')).toBeCloseTo(7.101, 6);
            expect(parseCatKHzToMHz('014320000')).toBeCloseTo(14.32, 6);
        });

        it('returns null for invalid or short strings', () => {
            expect(parseCatKHzToMHz('')).toBeNull();
            expect(parseCatKHzToMHz('123')).toBeNull();
            expect(parseCatKHzToMHz('ABC123456')).toBeNull();
            expect(parseCatKHzToMHz('000000000')).toBeNull();
        });
    });

    describe('parseUserMHz', () => {
        it('parses valid MHz strings', () => {
            expect(parseUserMHz('7.101')).toBeCloseTo(7.101, 6);
            expect(parseUserMHz(' 14.070 ')).toBeCloseTo(14.07, 6);
            expect(parseUserMHz('21,200')).toBeCloseTo(21.2, 6);
        });

        it('returns null for invalid input', () => {
            expect(parseUserMHz('')).toBeNull();
            expect(parseUserMHz('abc')).toBeNull();
            expect(parseUserMHz('-7.1')).toBeNull();
            expect(parseUserMHz('1001')).toBeNull();
        });
    });

    describe('frequencyToBandMHz', () => {
        it('maps typical HF/6m frequencies to correct bands', () => {
            expect(frequencyToBandMHz(1.9)).toBe('160m');
            expect(frequencyToBandMHz(3.7)).toBe('80m');
            expect(frequencyToBandMHz(5.3)).toBe('60m');
            expect(frequencyToBandMHz(7.1)).toBe('40m');
            expect(frequencyToBandMHz(10.12)).toBe('30m');
            expect(frequencyToBandMHz(14.1)).toBe('20m');
            expect(frequencyToBandMHz(18.1)).toBe('17m');
            expect(frequencyToBandMHz(21.2)).toBe('15m');
            expect(frequencyToBandMHz(24.93)).toBe('12m');
            expect(frequencyToBandMHz(28.5)).toBe('10m');
            expect(frequencyToBandMHz(50.5)).toBe('6m');
        });

        it('returns empty string for out-of-range or invalid values', () => {
            expect(frequencyToBandMHz(0.1)).toBe('');
            expect(frequencyToBandMHz(200)).toBe('');
            expect(frequencyToBandMHz(null)).toBe('');
            expect(frequencyToBandMHz(Number.NaN)).toBe('');
        });
    });

    describe('frequencyToBandFromCat', () => {
        it('maps CAT kHz strings to correct bands', () => {
            expect(frequencyToBandFromCat('007101000')).toBe('40m');
            expect(frequencyToBandFromCat('014320000')).toBe('20m');
            expect(frequencyToBandFromCat('000000000')).toBe('');
            expect(frequencyToBandFromCat('')).toBe('');
        });
    });

    describe('frequencyToBandFromMHz', () => {
        it('maps MHz strings to correct bands', () => {
            expect(frequencyToBandFromMHz('7.101')).toBe('40m');
            expect(frequencyToBandFromMHz('14.32')).toBe('20m');
            expect(frequencyToBandFromMHz('50.5')).toBe('6m');
        });

        it('returns empty string for invalid MHz strings', () => {
            expect(frequencyToBandFromMHz('')).toBe('');
            expect(frequencyToBandFromMHz('abc')).toBe('');
        });
    });

    describe('frequencyToBandFromDottedMHz', () => {
        it('maps dotted MHz strings to correct bands', () => {
            expect(frequencyToBandFromDottedMHz('7.200.000')).toBe('40m');
            expect(frequencyToBandFromDottedMHz('14.320.000')).toBe('20m');
            expect(frequencyToBandFromDottedMHz('50.500.000')).toBe('6m');
            expect(frequencyToBandFromDottedMHz(' 14.074.000 ')).toBe('20m');
            expect(frequencyToBandFromDottedMHz('07.200.000')).toBe('40m');
        });

        it('returns empty string for invalid or out-of-range dotted MHz strings', () => {
            expect(frequencyToBandFromDottedMHz('')).toBe('');
            expect(frequencyToBandFromDottedMHz('abc')).toBe('');
            expect(frequencyToBandFromDottedMHz('14,320.000')).toBe('');
            expect(frequencyToBandFromDottedMHz('00.000.000')).toBe('');
            expect(frequencyToBandFromDottedMHz('0.100.000')).toBe('');
            expect(frequencyToBandFromDottedMHz('200.000.000')).toBe('');
            expect(frequencyToBandFromDottedMHz('.7.200.000')).toBe('');
            expect(frequencyToBandFromDottedMHz('7..200.000')).toBe('');
        });
    });
});
