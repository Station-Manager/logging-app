export interface ContestTimers {
    sinceStartTimerID: number;
    sinceLastQsoTimerID: number;
    elapsedSinceStart: number;
    elapsedSinceLastQso: number;
    isRunning: boolean;
    start(): void;
    reset(): void;
    stop(): void;
}

export class ContestTimersClass implements ContestTimers {
    sinceStartTimerID: number;
    sinceLastQsoTimerID: number;
    elapsedSinceStart: number = $state(0);
    elapsedSinceLastQso: number = $state(0);
    isRunning: boolean = $state(false); // Add a reactive flag

    constructor() {
        this.sinceStartTimerID = 0;
        this.sinceLastQsoTimerID = 0;
        this.elapsedSinceStart = 0;
        this.elapsedSinceLastQso = 0;
    }
    start(): void {
        if (this.sinceStartTimerID === 0) {
            this.sinceStartTimerID = window.setInterval(() => {
                this.elapsedSinceStart += 1;
            }, 1000);
        }
    }
    reset(): void {
        const lastQso = this.elapsedSinceLastQso;
        if (lastQso !== 0) {
            window.clearInterval(this.sinceLastQsoTimerID);
            this.elapsedSinceLastQso = 0;
        }
        this.sinceLastQsoTimerID = window.setInterval(() => {
            this.elapsedSinceLastQso += 1;
        }, 1000);
    }
    stop(): void {
        window.clearInterval(this.sinceStartTimerID);
        this.sinceStartTimerID = 0;
        window.clearInterval(this.sinceLastQsoTimerID);
        this.sinceLastQsoTimerID = 0;
        this.elapsedSinceStart = 0;
        this.elapsedSinceLastQso = 0;
    }
}

export const contestTimers: ContestTimers = new ContestTimersClass();
