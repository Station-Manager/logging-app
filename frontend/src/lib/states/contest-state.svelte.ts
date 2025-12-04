export interface ContestState {
    totalQsos: number;
}

export const contestState: ContestState = $state({
    totalQsos: 0,
});
