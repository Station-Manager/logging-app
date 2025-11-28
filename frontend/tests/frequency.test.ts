import { describe, it, expect } from 'vitest';
import { formatCatKHzToDottedMHz, parseCatKHzToMHz, frequencyToBandFromCat } from '$lib/utils/frequency';

describe('frequency utils', () => {
    describe('parseCatKHzToMHz', () => {
        it('parses valid 9-digit CAT kHz strings', () => {
            expect(parseCatKHzToMHz('014320000')).toBe(14320 / 1000);
            expect(parseCatKHzToMHz('007200000')).toBe(7200 / 1000);
        });

        it('returns null for invalid inputs', () => {
            expect(parseCatKHzToMHz('')).toBeNull();
            expect(parseCatKHzToMHz('abc')).toBeNull();
            expect(parseCatKHzToMHz('12345')).toBeNull();
        });
    });

    describe('formatCatKHzToDottedMHz', () => {
        it('formats valid CAT kHz strings to dotted MHz', () => {
            expect(formatCatKHzToDottedMHz('014320000')).toBe('14.320.000');
            expect(formatCatKHzToDottedMHz('007200000')).toBe('7.200.000');
        });

        it('returns empty string for invalid input', () => {
            expect(formatCatKHzToDottedMHz('')).toBe('');
            expect(formatCatKHzToDottedMHz('abc')).toBe('');
            expect(formatCatKHzToDottedMHz('12345')).toBe('');
        });
    });

    describe('frequencyToBandFromCat', () => {
        it('maps known bands correctly', () => {
            // 14.320 MHz -> 20m
            expect(frequencyToBandFromCat('014320000')).toBe('20m');
            // 7.200 MHz -> 40m
            expect(frequencyToBandFromCat('007200000')).toBe('40m');
        });
    });
});

