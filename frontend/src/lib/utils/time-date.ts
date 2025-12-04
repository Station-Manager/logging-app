// Strict RFC3339-like format with explicit timezone offset (e.g. 2025-12-02T10:31:04+02:00).
const REMOTE_TIME_REGEX = /^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2}):(\d{2})([+-])(\d{2}):(\d{2})$/;

export function getTimeUTC(): string {
    const data: Date = new Date();
    return (
        data.getUTCHours().toString().padStart(2, '0') +
        ':' +
        data.getUTCMinutes().toString().padStart(2, '0')
    );
}
export function getDateUTC(): string {
    const data: Date = new Date();
    return (
        data.getUTCFullYear().toString() +
        '-' +
        (data.getUTCMonth() + 1).toString().padStart(2, '0') +
        '-' +
        data.getUTCDate().toString().padStart(2, '0')
    );
}
export function extractRemoteTime(timestamp?: string): string {
    if (!timestamp) {
        return '';
    }

    const trimmed = timestamp.trim();
    const match = REMOTE_TIME_REGEX.exec(trimmed);
    if (!match) {
        return '';
    }

    const [, , , , hour, minute, second, , offsetHour, offsetMinute] = match;
    const hourNum = Number(hour);
    const minuteNum = Number(minute);
    const secondNum = Number(second);
    const offsetHourNum = Number(offsetHour);
    const offsetMinuteNum = Number(offsetMinute);
    if (
        Number.isNaN(hourNum) ||
        Number.isNaN(minuteNum) ||
        Number.isNaN(secondNum) ||
        Number.isNaN(offsetHourNum) ||
        Number.isNaN(offsetMinuteNum) ||
        hourNum > 23 ||
        minuteNum > 59 ||
        secondNum > 59 ||
        offsetHourNum > 14 ||
        offsetMinuteNum > 59
    ) {
        return '';
    }

    return `${hour}:${minute}`;
}

export function formatTime(timeStr: string | undefined): string {
    if (timeStr === undefined) {
        return '';
    }
    if (timeStr.length !== 4) {
        throw new Error('Invalid time string length');
    }
    const hours: string = timeStr.slice(0, 2);
    const minutes: string = timeStr.slice(2, 4);
    return `${hours}:${minutes}`;
}

export function formatTimeSecondsToHHColonMMColonSS(seconds: number): string {
    const h = Math.floor(seconds / 3600)
        .toString()
        .padStart(2, '0');
    const m = Math.floor((seconds % 3600) / 60)
        .toString()
        .padStart(2, '0');
    const s = Math.floor(seconds % 60)
        .toString()
        .padStart(2, '0');
    return `${h}:${m}:${s}`;
}
