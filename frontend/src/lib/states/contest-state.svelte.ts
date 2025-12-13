export interface ContestState {
    isIncremental: boolean;
    totalQsos: number;
    currentStx: string;
    setCurrent(stx: string): void;
    getNext(): string;
}

export const contestState: ContestState = $state<ContestState>({
    isIncremental: true,
    // The total number of QSOs logged so far
    totalQsos: 0,
    currentStx: '001', // Default starting value
    // Set the current STX value
    setCurrent(this: ContestState, stx: string): void {
        let increment = true;
        let next = stx.trim();
        if (next.startsWith('!')) {
            increment = false;
            next = next.substring(1);
        }
        this.isIncremental = increment;
        this.currentStx = next;
    },
    // Get the next STX value
    getNext(this: ContestState): string {
        const stx: string = this.currentStx;
        if (this.isIncremental) {
            this.currentStx = doIncrement(this.currentStx);
        }
        return stx;
    }
});

const doIncrement = (cur: string): string => {
    const next = parseInt(cur, 10) + 1;
    if (Number.isNaN(next)) {
        return cur;
    }
    const nextStr = next.toString();
    const width = Math.max(cur.length, nextStr.length);
    return nextStr.padStart(width, "0");
}
