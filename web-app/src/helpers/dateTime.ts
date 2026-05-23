import dayjs from 'dayjs';
import utc from 'dayjs/plugin/utc';

dayjs.extend(utc);

export const DEFAULT_START_TIME = '00:00:00';
export const DEFAULT_COMPLETED_TIME = '00:01:00';

export function splitIsoDateTime(iso: string | null): { date: string; time: string } {
    if (!iso) {
        return { date: dayjs().format('YYYY-MM-DD'), time: '' };
    }
    const d = dayjs.utc(iso);
    return {
        date: d.format('YYYY-MM-DD'),
        time: d.format('HH:mm:ss'),
    };
}

export function resolveTimePart(time: string, fallback: string): string {
    return time || fallback;
}

export function toUtcIsoString(date: string, time: string, fallbackTime: string): string {
    return `${date}T${resolveTimePart(time, fallbackTime)}Z`;
}

export function parseUtcIso(date: string, time: string, fallbackTime: string) {
    return dayjs.utc(toUtcIsoString(date, time, fallbackTime));
}

export function isValidUtcDateTime(date: string, time: string, fallbackTime: string): boolean {
    return parseUtcIso(date, time, fallbackTime).isValid();
}

export function isCompletedAfterStart(
    startDate: string,
    startTime: string,
    completedDate: string,
    completedTime: string,
): boolean {
    const start = parseUtcIso(startDate, startTime, DEFAULT_START_TIME);
    const completed = parseUtcIso(completedDate, completedTime, DEFAULT_COMPLETED_TIME);
    return completed.isAfter(start);
}

export function formatDisplayDate(dateString: string): string {
    return dayjs(dateString).format('DD/MM/YYYY');
}

export function formatLocaleDate(dateString: string): string {
    return dayjs(dateString).toDate().toLocaleDateString();
}

export function sortByStartedAtDesc<T extends { started_at: string }>(games: T[]): T[] {
    return [...games].sort((a, b) => dayjs(b.started_at).valueOf() - dayjs(a.started_at).valueOf());
}
