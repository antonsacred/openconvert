import type { DownloadEntry, PersistedQueueItem, QueueItem } from './types.ts';

const UPLOAD_QUEUE_STORAGE_KEY = 'openconvert.upload_queue_state';
const ORIGINAL_FILE_STORE = new Map<string, File>();
const DOWNLOAD_STORE = new Map<string, DownloadEntry>();

export function setOriginalFile(itemId: string, file: File): void {
    ORIGINAL_FILE_STORE.set(itemId, file);
}

export function getOriginalFile(itemId: string): File | undefined {
    return ORIGINAL_FILE_STORE.get(itemId);
}

export function deleteOriginalFile(itemId: string): void {
    ORIGINAL_FILE_STORE.delete(itemId);
}

export function hasOriginalFile(itemId: string): boolean {
    return ORIGINAL_FILE_STORE.has(itemId);
}

export function hasDownload(itemId: string): boolean {
    return DOWNLOAD_STORE.has(itemId);
}

export function getDownload(itemId: string): DownloadEntry | undefined {
    return DOWNLOAD_STORE.get(itemId);
}

export function listDownloadsByItemIds(itemIds: string[]): Array<{ itemId: string; download: DownloadEntry }> {
    const entries: Array<{ itemId: string; download: DownloadEntry }> = [];
    itemIds.forEach((itemId) => {
        const download = DOWNLOAD_STORE.get(itemId);
        if (download === undefined) {
            return;
        }

        entries.push({
            itemId,
            download,
        });
    });

    return entries;
}

export function setDownload(itemId: string, fileName: string, mimeType: string, blob: Blob): void {
    revokeDownload(itemId);
    const objectUrl = URL.createObjectURL(blob);
    DOWNLOAD_STORE.set(itemId, {
        blob,
        objectUrl,
        fileName,
        mimeType,
    });
}

export function revokeDownload(itemId: string): void {
    const download = DOWNLOAD_STORE.get(itemId);
    if (download === undefined) {
        return;
    }

    URL.revokeObjectURL(download.objectUrl);
    DOWNLOAD_STORE.delete(itemId);
}

export function hydrateDoneStateFromDownloads(items: QueueItem[]): QueueItem[] {
    return items.map((item) => {
        const download = DOWNLOAD_STORE.get(item.id);
        if (download === undefined) {
            return item;
        }

        return {
            ...item,
            status: 'DONE',
            errorMessage: '',
            outputFileName: download.fileName,
        };
    });
}

export function syncTransientStoresWithItems(items: QueueItem[]): void {
    const activeItemIds = new Set(items.map((item) => item.id));

    Array.from(ORIGINAL_FILE_STORE.keys()).forEach((itemId) => {
        if (!activeItemIds.has(itemId)) {
            ORIGINAL_FILE_STORE.delete(itemId);
        }
    });

    Array.from(DOWNLOAD_STORE.entries()).forEach(([itemId, download]) => {
        if (!activeItemIds.has(itemId)) {
            URL.revokeObjectURL(download.objectUrl);
            DOWNLOAD_STORE.delete(itemId);
        }
    });
}

export function readPersistedQueueItems(): QueueItem[] {
    const rawState = readPersistedState();
    if (rawState === null || rawState === '') {
        return [];
    }

    let decodedState: unknown = null;
    try {
        decodedState = JSON.parse(rawState);
    } catch {
        clearPersistedQueueItems();

        return [];
    }

    if (!Array.isArray(decodedState)) {
        clearPersistedQueueItems();

        return [];
    }

    const parsedItems: QueueItem[] = [];
    decodedState.forEach((value) => {
        if (!isPersistedQueueItem(value)) {
            return;
        }

        parsedItems.push({
            id: value.id,
            fileName: value.fileName,
            source: value.source,
            target: value.target,
            status: 'READY',
            errorMessage: '',
            outputFileName: '',
        });
    });

    return parsedItems;
}

export function persistQueueItems(items: QueueItem[]): void {
    if (items.length === 0) {
        clearPersistedQueueItems();

        return;
    }

    const persistedItems: PersistedQueueItem[] = items.map((item) => ({
        id: item.id,
        fileName: item.fileName,
        source: item.source,
        target: item.target,
    }));

    try {
        window.sessionStorage.setItem(UPLOAD_QUEUE_STORAGE_KEY, JSON.stringify(persistedItems));
    } catch {
        try {
            window.localStorage.setItem(UPLOAD_QUEUE_STORAGE_KEY, JSON.stringify(persistedItems));
        } catch {
        }
    }
}

export function clearPersistedQueueItems(): void {
    try {
        window.sessionStorage.removeItem(UPLOAD_QUEUE_STORAGE_KEY);
    } catch {
    }
    try {
        window.localStorage.removeItem(UPLOAD_QUEUE_STORAGE_KEY);
    } catch {
    }
}

function readPersistedState(): string | null {
    try {
        const sessionState = window.sessionStorage.getItem(UPLOAD_QUEUE_STORAGE_KEY);
        if (sessionState !== null && sessionState !== '') {
            return sessionState;
        }
    } catch {
    }

    try {
        const localState = window.localStorage.getItem(UPLOAD_QUEUE_STORAGE_KEY);
        if (localState !== null && localState !== '') {
            return localState;
        }
    } catch {
    }

    return null;
}

function isPersistedQueueItem(value: unknown): value is PersistedQueueItem {
    if (typeof value !== 'object' || value === null) {
        return false;
    }

    const candidate = value as Record<string, unknown>;

    return typeof candidate.id === 'string'
        && candidate.id.trim() !== ''
        && typeof candidate.fileName === 'string'
        && candidate.fileName.trim() !== ''
        && typeof candidate.source === 'string'
        && candidate.source.trim() !== ''
        && typeof candidate.target === 'string';
}
