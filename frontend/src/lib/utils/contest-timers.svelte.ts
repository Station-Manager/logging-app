export interface ContestTimers {
    sinceStartTimerID: number;
    sinceLastQsoTimerID: number;
    elapsedSinceStart: number;
    elapasedSinceLastQso: number;
    start(): void;
    reset(): void;
    stop(): void;
}

export class ContestTimersClass implements ContestTimers {
    sinceStartTimerID: number;
    sinceLastQsoTimerID: number;
    elapsedSinceStart: number;
    elapasedSinceLastQso: number;

    constructor() {
        this.sinceStartTimerID = 0;
        this.sinceLastQsoTimerID = 0;
        this.elapsedSinceStart = 0;
        this.elapasedSinceLastQso = 0;
    }
    start(): void {
        if (this.sinceStartTimerID === 0) {
            this.sinceStartTimerID = window.setInterval(() => {
                this.elapsedSinceStart += 1;
            }, 1000);
        }
    }
    reset(): void {
        const lastQso = this.elapasedSinceLastQso;
        if (lastQso !== 0) {
            window.clearInterval(this.sinceLastQsoTimerID);
            this.elapasedSinceLastQso = 0;
        }
        this.sinceLastQsoTimerID = window.setInterval(() => {
            this.elapasedSinceLastQso += 1;
        }, 1000);
    }
    stop(): void {
        window.clearInterval(this.sinceStartTimerID);
        this.sinceStartTimerID = 0;
        window.clearInterval(this.sinceLastQsoTimerID);
        this.sinceLastQsoTimerID = 0;
        this.elapsedSinceStart = 0;
        this.elapasedSinceLastQso = 0;
    }
}

export const contestTimers: ContestTimers = new ContestTimersClass();
