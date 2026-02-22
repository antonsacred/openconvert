import type { ConvertSuccessPayload } from './types.ts';

type ConvertRequestInput = {
    from: string;
    to: string;
    fileName: string;
    file: File;
};

type ConvertResponse = {
    ok: boolean;
    payload: unknown;
};

export async function submitConversionRequest(convertUrl: string, input: ConvertRequestInput): Promise<ConvertResponse> {
    const formData = new FormData();
    formData.append('from', input.from);
    formData.append('to', input.to);
    formData.append('file', input.file, input.fileName);

    const response = await fetch(convertUrl, {
        method: 'POST',
        headers: {
            Accept: 'application/json',
        },
        body: formData,
    });

    let payload: unknown = null;
    try {
        payload = await response.json();
    } catch {
        payload = null;
    }

    return {
        ok: response.ok,
        payload,
    };
}

export function extractErrorMessage(payload: unknown, fallback: string): string {
    if (typeof payload === 'object' && payload !== null) {
        const objectPayload = payload as Record<string, unknown>;
        if (typeof objectPayload.message === 'string' && objectPayload.message.trim() !== '') {
            return objectPayload.message.trim();
        }

        if (typeof objectPayload.error === 'object' && objectPayload.error !== null) {
            const errorPayload = objectPayload.error as Record<string, unknown>;
            if (typeof errorPayload.message === 'string' && errorPayload.message.trim() !== '') {
                return errorPayload.message.trim();
            }
        }
    }

    return fallback;
}

export function isConvertSuccessPayload(payload: unknown): payload is ConvertSuccessPayload {
    if (typeof payload !== 'object' || payload === null) {
        return false;
    }

    const objectPayload = payload as Record<string, unknown>;

    return typeof objectPayload.fileName === 'string'
        && objectPayload.fileName.trim() !== ''
        && typeof objectPayload.mimeType === 'string'
        && objectPayload.mimeType.trim() !== ''
        && typeof objectPayload.contentBase64 === 'string'
        && objectPayload.contentBase64.trim() !== '';
}

export function decodeBase64ToBlob(contentBase64: string, mimeType: string): Blob {
    const binary = window.atob(contentBase64);
    const bytes = new Uint8Array(binary.length);
    for (let i = 0; i < binary.length; i += 1) {
        bytes[i] = binary.charCodeAt(i);
    }

    return new Blob([bytes], { type: mimeType });
}
