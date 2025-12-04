import { qsoState } from '$lib/states/qso-state.svelte';
import { getTimeUTC } from '$lib/utils/time-date';

const FIVE_SECONDS = 5000;

let intervalID: number = 0;

/**
 * Initializes and starts a timer to update QSO state properties at regular intervals.
 *
 * The function sets up a recurring timer that updates the `time_on` property of the `qsoState`
 * object with the current UTC time every five seconds. If a timer is already running, it stops
 * the existing timer before starting a new one. Additionally, it synchronizes the `time_off`
 * property with the `time_on` property during initialization.
 */
export const startQsoTimer = (): void => {
    // If the interval is already set, clear it before setting a new one
    if (intervalID !== 0) {
        clearInterval(intervalID);
        intervalID = 0;
    }
    intervalID = window.setInterval(() => {
        qsoState.time_on = getTimeUTC();
    }, FIVE_SECONDS);
    qsoState.time_off = qsoState.time_on;
};

/**
 * Stops the currently running QSO timer by clearing the interval associated with it
 * and resetting the interval identifier to 0. Updates the QSO state with the current
 * UTC time for both the start (`time_on`) and stop (`time_off`) times.
 */
export const stopQsoTimer = (): void => {
    if (intervalID !== 0) {
        clearInterval(intervalID);
        intervalID = 0;
        qsoState.time_on = getTimeUTC();
        qsoState.time_off = qsoState.time_on;
    }
};

export const isTimerRunning = (): boolean => {
    return intervalID !== 0;
};
