export interface CatState {
    identity: string;
    vfoaFreq: string;
    vfobFreq: string;
    select: string;
    split: string;
}

export const catState: CatState = $state({
    identity: '',
    vfoaFreq: '',
    vfobFreq: '',
    select: 'VFO-A',
    split: 'OFF',
});
