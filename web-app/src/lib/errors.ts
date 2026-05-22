type AnyRecord = Record<string, unknown>;

const errorDictionary: Record<string, string> = {
    Unauthorized: 'Не авторизовано',
    'username or password is incorrect': 'Неверное имя пользователя или пароль',
    'player cannot access this entity': 'Нет доступа к этой записи',
    'player id is invalid': 'Некорректный идентификатор игрока',
    'player not found': 'Игрок не найден',
    'could not issue token': 'Не удалось выдать токен',
    'set password': 'Ошибка смены пароля',
    register: 'Ошибка регистрации',
};

function translateMessage(message: string): string {
    const direct = errorDictionary[message];
    if (direct) return direct;
    const lower = message.toLowerCase();
    for (const [key, value] of Object.entries(errorDictionary)) {
        if (lower.includes(key.toLowerCase())) return value;
    }
    return message;
}

export function parseApiError(status: number, statusText: string, rawBody: string): string {
    let fallback = `HTTP ${status}: ${statusText}`;
    if (!rawBody) return fallback;

    try {
        const data = JSON.parse(rawBody) as AnyRecord | string;
        if (typeof data === 'string') {
            return translateMessage(data);
        }

        const detail = typeof data.detail === 'string' ? data.detail : undefined;
        const message = typeof data.message === 'string' ? data.message : undefined;
        const title = typeof data.title === 'string' ? data.title : undefined;

        const candidate = detail || message || title;
        if (candidate) {
            return translateMessage(candidate);
        }
    } catch {
        return fallback;
    }

    return fallback;
}

