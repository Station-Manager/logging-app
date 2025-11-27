const DEFAULT_MAX_LENGTH = 10;
const SANE_MAX_LENGTH = 30;

export interface ClipboardState {
    /**
     * Read-only view of the clipboard list; mutate via methods only.
     */
    readonly list: readonly string[];
    /**
     * Read-only view of the maximum history length; change via setMaxLength.
     */
    readonly maxLength: number;
    add(this: ClipboardState, item: string): void;
    get(this: ClipboardState, index: number): string | undefined;
    setMaxLength(this: ClipboardState, len: number): void;
}

// Internal mutable implementation that backs the exported, read-only interface.
interface ClipboardStateInternal extends Omit<ClipboardState, 'list' | 'maxLength'> {
    list: string[];
    maxLength: number;
}

export const clipboardState: ClipboardState = $state<ClipboardStateInternal>({
    // Internal mutable state; exposed as read-only via the interface above.
    list: [],
    maxLength: DEFAULT_MAX_LENGTH,
    add(this: ClipboardStateInternal, item: string): void {
        const value = item.trim();
        if (!value) return;

        // If the value already exists, move it to the most recent position
        // instead of skipping it, to reflect recency in the history.
        const index = this.list.indexOf(value);
        if (index !== -1) {
            this.list.splice(index, 1);
        }

        this.list.push(value);

        // Efficiently prune oldest entries in a single operation when over limit.
        if (this.list.length > this.maxLength) {
            const excess = this.list.length - this.maxLength;
            this.list.splice(0, excess);
        }
    },
    get(this: ClipboardStateInternal, index: number): string | undefined {
        if (!Number.isInteger(index)) return undefined;
        if (index < 0 || index >= this.list.length) return undefined;
        return this.list[index];
    },
    setMaxLength(this: ClipboardStateInternal, len: number): void {
        // Treat non-finite values as 0 and ensure integer >= 0, then clamp to sane upper bound.
        const n = Number.isFinite(len) ? Math.floor(len) : 0;
        this.maxLength = Math.min(SANE_MAX_LENGTH, Math.max(0, n));
        if (this.list.length > this.maxLength) {
            const excess: number = this.list.length - this.maxLength;
            this.list.splice(0, excess);
        }
    },
}) as ClipboardState;
