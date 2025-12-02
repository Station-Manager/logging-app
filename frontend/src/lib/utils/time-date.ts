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
    const match = /^\d{4}-\d{2}-\d{2}T(\d{2}):(\d{2}):\d{2}[+-]\d{2}:\d{2}$/.exec(trimmed);
    if (!match) {
        return '';
    }

    return `${match[1]}:${match[2]}`;
}
