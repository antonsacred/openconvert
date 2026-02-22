import { zipSync } from 'fflate';

type ZipEntry = {
    fileName: string;
    blob: Blob;
};

const ZIP_MIME_TYPE = 'application/zip';

export async function buildZipBlob(entries: ZipEntry[]): Promise<Blob> {
    if (entries.length === 0) {
        throw new Error('No files available for ZIP archive.');
    }

    const usedNames = new Set<string>();
    const files: Record<string, Uint8Array> = {};

    for (let i = 0; i < entries.length; i += 1) {
        const entry = entries[i];
        const sanitizedName = sanitizeFileName(entry.fileName, i);
        const uniqueName = createUniqueName(sanitizedName, usedNames);
        usedNames.add(uniqueName);

        const arrayBuffer = await entry.blob.arrayBuffer();
        files[uniqueName] = new Uint8Array(arrayBuffer);
    }

    const compressed = zipSync(files, { level: 6 });

    return new Blob([compressed], { type: ZIP_MIME_TYPE });
}

export function buildArchiveFileName(now: Date): string {
    const year = now.getFullYear();
    const month = leftPad(now.getMonth() + 1);
    const day = leftPad(now.getDate());
    const hours = leftPad(now.getHours());
    const minutes = leftPad(now.getMinutes());

    return `openconvert-${year}${month}${day}-${hours}${minutes}.zip`;
}

function sanitizeFileName(fileName: string, index: number): string {
    const normalized = fileName
        .trim()
        .replace(/[\\/:*?"<>|\u0000-\u001f]/g, '_')
        .replace(/\s+/g, ' ');

    if (normalized !== '') {
        return normalized;
    }

    return `file-${index + 1}`;
}

function createUniqueName(fileName: string, usedNames: Set<string>): string {
    if (!usedNames.has(fileName)) {
        return fileName;
    }

    const extensionIndex = fileName.lastIndexOf('.');
    const hasExtension = extensionIndex > 0 && extensionIndex < fileName.length - 1;
    const baseName = hasExtension ? fileName.slice(0, extensionIndex) : fileName;
    const extension = hasExtension ? fileName.slice(extensionIndex) : '';

    let suffix = 1;
    while (true) {
        const candidate = `${baseName}(${suffix})${extension}`;
        if (!usedNames.has(candidate)) {
            return candidate;
        }

        suffix += 1;
    }
}

function leftPad(value: number): string {
    return String(value).padStart(2, '0');
}
