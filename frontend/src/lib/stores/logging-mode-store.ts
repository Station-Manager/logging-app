import { writable, type Writable } from 'svelte/store';

// Define a narrow union for keys
export type LoggingModeKey = 'normal' | 'contest';

export type LoggingModeMap = Record<LoggingModeKey, string>;

// Central list of modes
export const loggingModes: LoggingModeMap = {
    normal: 'Normal',
    contest: 'Contest',
};

// Central default key and value
export const DEFAULT_LOGGING_MODE_KEY: LoggingModeKey = 'normal';
export const DEFAULT_LOGGING_MODE_LABEL = loggingModes[DEFAULT_LOGGING_MODE_KEY];

// Store holds just the key; components can derive label as needed
export const loggingModeStore: Writable<LoggingModeKey> = writable(DEFAULT_LOGGING_MODE_KEY);
writable(loggingModes.normal);
