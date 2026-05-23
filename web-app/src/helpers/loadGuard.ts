/** Отменяет устаревшие async-загрузки при смене ключа (например, id маршрута). */
export function createLoadGuard() {
    let generation = 0;

    return {
        isCurrent(gen: number) {
            return gen === generation;
        },
        next() {
            return ++generation;
        },
        cancel() {
            generation++;
        },
    };
}
