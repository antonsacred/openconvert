export type FormatsBySource = Record<string, string[]>;

export type QueueItemStatus = 'READY' | 'PROCESSING' | 'DONE' | 'ERROR';

export type QueueItem = {
    id: string;
    fileName: string;
    source: string;
    target: string;
    status: QueueItemStatus;
    errorMessage: string;
    outputFileName: string;
};

export type PersistedQueueItem = {
    id: string;
    fileName: string;
    source: string;
    target: string;
};

export type DownloadEntry = {
    blob: Blob;
    objectUrl: string;
    fileName: string;
    mimeType: string;
};

export type ConvertSuccessPayload = {
    fileName: string;
    mimeType: string;
    contentBase64: string;
};
