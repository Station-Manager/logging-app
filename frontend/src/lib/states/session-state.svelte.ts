import { types } from '$lib/wailsjs/go/models';

export interface SessionState {
    total: number;
    list: types.Qso[];
    update(this: SessionState, list: types.Qso[]): void;
    start(this: SessionState): void;
    stop(this: SessionState): void;
    reset(this: SessionState): void;
}

export interface SessionTimeState {
    elapsed: number;
}

let sessionIntervalID: number | null = null;

// Cleanup function to be called when the state is no longer needed
export const cleanupSessionState = (): void => {
    if (sessionIntervalID !== null) {
        window.clearInterval(sessionIntervalID);
        sessionIntervalID = null;
    }
};

// Declare sessionTimeState before sessionState to avoid forward reference
export const sessionTimeState: SessionTimeState = $state({
    elapsed: 0,
});

export const sessionState: SessionState = $state({
    total: 0,
    list: [],
    update(list: types.Qso[]): void {
        this.list = list;
        this.total = list.length;
    },
    start(): void {
        cleanupSessionState();
        sessionTimeState.elapsed = 0;
        sessionIntervalID = window.setInterval(() => {
            sessionTimeState.elapsed += 1;
        }, 1000);
    },
    stop(): void {
        cleanupSessionState();
    },
    reset(): void {
        cleanupSessionState();
        sessionTimeState.elapsed = 0;
        this.list = [];
        this.total = 0;
    },
});

// Helper function to format seconds as HH:MM:SS
export const formatSeconds = (totalSeconds: number): string => {
    const seconds: number = Math.max(0, totalSeconds);
    const hours: number = Math.floor(seconds / 3600);
    const minutes: number = Math.floor((seconds % 3600) / 60);
    const secs: number = seconds % 60;

    return [hours, minutes, secs].map((v) => v.toString().padStart(2, '0')).join(':');
};

// Derived reactive value for formatted elapsed time
export const sessionElapsedTime = $derived.by((): string => {
    return formatSeconds(sessionTimeState.elapsed);
});
