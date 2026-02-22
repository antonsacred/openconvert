import { Controller } from '@hotwired/stimulus';

type FormatsBySource = Record<string, string[]>;
type QueueItemStatus = 'READY' | 'PROCESSING' | 'DONE' | 'ERROR';
type QueueItem = {
    id: string;
    fileName: string;
    source: string;
    target: string;
    status: QueueItemStatus;
    errorMessage: string;
    outputFileName: string;
};
type PersistedQueueItem = {
    id: string;
    fileName: string;
    source: string;
    target: string;
};
type DownloadEntry = {
    objectUrl: string;
    fileName: string;
    mimeType: string;
};
type ConvertSuccessPayload = {
    fileName: string;
    mimeType: string;
    contentBase64: string;
};

const UPLOAD_QUEUE_STORAGE_KEY = 'openconvert.upload_queue_state';
const ORIGINAL_FILE_STORE = new Map<string, File>();
const DOWNLOAD_STORE = new Map<string, DownloadEntry>();

export default class extends Controller<HTMLElement> {
    static targets = ['fileInput', 'fileList', 'emptyState', 'queuePanel', 'error', 'errorMessage', 'convertButton'] as const;
    static values = {
        formatsBySource: Object,
        sourcePageTemplate: String,
        selectedFrom: String,
        convertUrl: String,
    } as const;

    declare readonly fileInputTarget: HTMLInputElement;
    declare readonly fileListTarget: HTMLUListElement;
    declare readonly emptyStateTarget: HTMLElement;
    declare readonly queuePanelTarget: HTMLElement;
    declare readonly hasQueuePanelTarget: boolean;
    declare readonly errorTarget: HTMLElement;
    declare readonly errorMessageTarget: HTMLElement;
    declare readonly hasErrorTarget: boolean;
    declare readonly hasErrorMessageTarget: boolean;
    declare readonly convertButtonTarget: HTMLButtonElement;
    declare readonly hasConvertButtonTarget: boolean;
    declare readonly formatsBySourceValue: FormatsBySource;
    declare readonly sourcePageTemplateValue: string;
    declare readonly selectedFromValue: string;
    declare readonly hasSelectedFromValue: boolean;
    declare readonly convertUrlValue: string;
    declare readonly hasConvertUrlValue: boolean;

    private items: QueueItem[] = [];
    private isConverting = false;

    connect(): void {
        this.clearError();
        this.items = this.loadPersistedItems();
        this.pruneItemsWithoutSourceData();
        this.reconcileItemsWithCurrentSelection();
        this.hydrateDoneStateFromDownloadStore();
        this.syncTransientStoresWithItems();
        this.render();
    }

    openFilePicker(): void {
        this.clearError();
        this.fileInputTarget.click();
    }

    onFilesSelected(): void {
        const selectedFiles = Array.from(this.fileInputTarget.files ?? []);
        if (selectedFiles.length === 0) {
            return;
        }

        this.clearError();

        const supportedFiles: Array<{ file: File; fileName: string; source: string }> = [];
        const unsupportedFiles: string[] = [];
        selectedFiles.forEach((file) => {
            const detectedSource = this.detectSourceFromFileName(file.name);
            if (detectedSource === null) {
                unsupportedFiles.push(file.name);

                return;
            }

            supportedFiles.push({
                file,
                fileName: file.name,
                source: detectedSource,
            });
        });

        if (supportedFiles.length === 0) {
            this.showError(this.buildUnsupportedFilesMessage(unsupportedFiles));
            this.fileInputTarget.value = '';

            return;
        }

        const isFirstUpload = this.items.length === 0;
        const firstUploadedSource = supportedFiles[0].source;
        const pageTarget = this.currentPageTarget();

        const queueItemsToAdd: QueueItem[] = [];
        supportedFiles.forEach((supportedFile) => {
            const defaultTarget = this.resolvePreferredTarget(supportedFile.source, pageTarget);
            const queueItem = this.createQueueItem(supportedFile.fileName, supportedFile.source, defaultTarget);
            queueItemsToAdd.push(queueItem);
            ORIGINAL_FILE_STORE.set(queueItem.id, supportedFile.file);
        });

        this.items = [...this.items, ...queueItemsToAdd];
        this.persistItems();
        this.render();

        if (isFirstUpload) {
            const currentSource = this.currentSource();
            if (currentSource !== firstUploadedSource) {
                this.visit(this.buildSourcePageUrl(firstUploadedSource));
            }
        }

        if (unsupportedFiles.length > 0) {
            this.showError(this.buildUnsupportedFilesMessage(unsupportedFiles));
        }

        this.fileInputTarget.value = '';
    }

    removeFile(event: Event): void {
        const button = event.currentTarget;
        if (!(button instanceof HTMLButtonElement)) {
            return;
        }

        const itemId = button.dataset.itemId ?? '';
        if (itemId === '') {
            return;
        }

        this.removeQueueItem(itemId);
    }

    onRowTargetChange(event: Event): void {
        const select = event.currentTarget;
        if (!(select instanceof HTMLSelectElement)) {
            return;
        }

        const itemId = select.dataset.itemId ?? '';
        if (itemId === '') {
            return;
        }

        const itemIndex = this.items.findIndex((item) => item.id === itemId);
        if (itemIndex < 0) {
            return;
        }

        const item = this.items[itemIndex];
        const normalizedTarget = this.resolvePreferredTarget(item.source, this.normalizeFormat(select.value));
        this.items[itemIndex] = {
            ...item,
            target: normalizedTarget,
            status: 'READY',
            errorMessage: '',
            outputFileName: '',
        };
        this.revokeDownload(item.id);

        this.persistItems();
        this.render();
    }

    async convertQueue(): Promise<void> {
        if (this.isConverting) {
            return;
        }

        const removedItemsCount = this.pruneItemsWithoutSourceData();
        if (removedItemsCount > 0) {
            this.render();
        }

        if (this.items.length === 0) {
            return;
        }

        const itemsToConvert = this.items.filter((item) => item.status !== 'DONE');
        if (itemsToConvert.length === 0) {
            return;
        }

        if (!this.hasConvertUrlValue || this.convertUrlValue.trim() === '') {
            this.showError('Conversion endpoint is not configured.');

            return;
        }

        if (itemsToConvert.some((item) => item.target === '')) {
            this.showError('Choose target format for every file before converting.');

            return;
        }

        this.clearError();
        this.isConverting = true;
        this.render();

        const queueSnapshot = [...itemsToConvert];
        for (const item of queueSnapshot) {
            await this.convertSingleItem(item.id);
        }

        this.isConverting = false;
        this.render();
    }

    downloadFile(event: Event): void {
        const button = event.currentTarget;
        if (!(button instanceof HTMLButtonElement)) {
            return;
        }

        const itemId = button.dataset.itemId ?? '';
        if (itemId === '') {
            return;
        }

        const download = DOWNLOAD_STORE.get(itemId);
        if (download === undefined) {
            return;
        }

        const link = document.createElement('a');
        link.href = download.objectUrl;
        link.download = download.fileName;
        link.rel = 'noopener';
        link.style.display = 'none';
        document.body.appendChild(link);
        link.click();
        link.remove();
    }

    private async convertSingleItem(itemId: string): Promise<void> {
        const itemIndex = this.items.findIndex((item) => item.id === itemId);
        if (itemIndex < 0) {
            return;
        }

        const item = this.items[itemIndex];
        if (item.target === '') {
            this.setItemError(itemId, 'Target format is required before conversion.');

            return;
        }

        const sourceFile = ORIGINAL_FILE_STORE.get(itemId);
        if (sourceFile === undefined) {
            this.removeQueueItem(itemId);

            return;
        }

        this.items[itemIndex] = {
            ...item,
            status: 'PROCESSING',
            errorMessage: '',
            outputFileName: '',
        };
        this.revokeDownload(itemId);
        this.render();

        try {
            const formData = new FormData();
            formData.append('from', item.source);
            formData.append('to', item.target);
            formData.append('file', sourceFile, item.fileName);

            const response = await fetch(this.convertUrlValue, {
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

            if (!response.ok) {
                this.setItemError(itemId, this.extractErrorMessage(payload, `Failed to convert ${item.fileName}.`));

                return;
            }

            if (!this.isConvertSuccessPayload(payload)) {
                this.setItemError(itemId, `Invalid conversion response for ${item.fileName}.`);

                return;
            }

            const outputBlob = this.decodeBase64ToBlob(payload.contentBase64, payload.mimeType);
            this.setDownload(itemId, payload.fileName, payload.mimeType, outputBlob);

            const doneIndex = this.items.findIndex((queueItem) => queueItem.id === itemId);
            if (doneIndex >= 0) {
                this.items[doneIndex] = {
                    ...this.items[doneIndex],
                    status: 'DONE',
                    errorMessage: '',
                    outputFileName: payload.fileName,
                };
            }
            this.render();
        } catch {
            this.setItemError(itemId, `Failed to convert ${item.fileName}.`);
        }
    }

    private render(): void {
        this.fileListTarget.innerHTML = '';

        if (this.items.length === 0) {
            this.emptyStateTarget.classList.remove('hidden');
            if (this.hasQueuePanelTarget) {
                this.queuePanelTarget.classList.add('hidden');
            }
            if (this.hasConvertButtonTarget) {
                this.convertButtonTarget.disabled = true;
                this.convertButtonTarget.textContent = 'Convert';
            }

            return;
        }

        this.emptyStateTarget.classList.add('hidden');
        if (this.hasQueuePanelTarget) {
            this.queuePanelTarget.classList.remove('hidden');
        }

        const fragment = document.createDocumentFragment();
        this.items.forEach((item) => {
            fragment.appendChild(this.buildItemRow(item));
        });
        this.fileListTarget.appendChild(fragment);

        if (this.hasConvertButtonTarget) {
            const itemsToConvert = this.items.filter((item) => item.status !== 'DONE');
            const canConvert = !this.isConverting
                && itemsToConvert.length > 0
                && itemsToConvert.every((item) => item.target !== '');
            this.convertButtonTarget.disabled = !canConvert;

            if (this.isConverting) {
                this.convertButtonTarget.innerHTML = '<span class="loading loading-spinner loading-sm"></span>Converting...';
            } else {
                this.convertButtonTarget.textContent = 'Convert';
            }
        }
    }

    private buildItemRow(item: QueueItem): HTMLLIElement {
        const row = document.createElement('li');
        row.className = 'flex flex-wrap items-center gap-4 p-4 md:flex-nowrap';

        const fileInfo = document.createElement('div');
        fileInfo.className = 'min-w-0 flex-1';
        const fileName = document.createElement('p');
        fileName.className = 'truncate text-lg font-medium';
        fileName.textContent = item.fileName;
        const sourceLabel = document.createElement('p');
        sourceLabel.className = 'text-sm text-base-content/60';
        sourceLabel.textContent = `Source: ${item.source.toUpperCase()}`;
        fileInfo.append(fileName, sourceLabel);

        if (item.errorMessage !== '') {
            const errorLabel = document.createElement('p');
            errorLabel.className = 'mt-1 text-sm text-error';
            errorLabel.textContent = item.errorMessage;
            fileInfo.append(errorLabel);
        } else if (item.status === 'DONE' && item.outputFileName !== '') {
            const doneLabel = document.createElement('p');
            doneLabel.className = 'mt-1 text-sm text-success';
            doneLabel.textContent = `Output: ${item.outputFileName}`;
            fileInfo.append(doneLabel);
        }

        const controls = document.createElement('div');
        controls.className = 'flex shrink-0 items-center gap-2 whitespace-nowrap';

        const statusBadge = this.buildStatusBadge(item.status);

        const label = document.createElement('span');
        label.className = 'text-sm text-base-content/70';
        label.textContent = 'Convert to';

        const targetSelect = document.createElement('select');
        targetSelect.className = 'select select-bordered select-sm min-w-28';
        targetSelect.dataset.itemId = item.id;
        targetSelect.setAttribute('data-action', 'change->upload-queue#onRowTargetChange');
        targetSelect.disabled = this.isConverting || item.status === 'PROCESSING';
        targetSelect.add(new Option('to', '', item.target === ''));
        this.resolveTargetsForSource(item.source).forEach((target) => {
            const normalizedTarget = this.normalizeFormat(target);
            if (normalizedTarget === '') {
                return;
            }

            targetSelect.add(new Option(normalizedTarget.toUpperCase(), normalizedTarget, false, normalizedTarget === item.target));
        });

        const removeButton = document.createElement('button');
        removeButton.type = 'button';
        removeButton.className = 'btn btn-ghost btn-sm text-error';
        removeButton.dataset.itemId = item.id;
        removeButton.disabled = this.isConverting || item.status === 'PROCESSING';
        removeButton.setAttribute('data-action', 'click->upload-queue#removeFile');
        removeButton.textContent = 'X';

        controls.append(statusBadge, label, targetSelect);

        if (item.status === 'DONE' && DOWNLOAD_STORE.has(item.id)) {
            const downloadButton = document.createElement('button');
            downloadButton.type = 'button';
            downloadButton.className = 'btn btn-primary btn-sm';
            downloadButton.dataset.itemId = item.id;
            downloadButton.setAttribute('data-action', 'click->upload-queue#downloadFile');
            downloadButton.textContent = 'Download';
            controls.append(downloadButton);
        }

        controls.append(removeButton);
        row.append(fileInfo, controls);

        return row;
    }

    private buildStatusBadge(status: QueueItemStatus): HTMLSpanElement {
        const badge = document.createElement('span');
        badge.classList.add('badge', 'badge-sm');

        if (status === 'PROCESSING') {
            badge.classList.add('badge-info', 'gap-1');
            const spinner = document.createElement('span');
            spinner.className = 'loading loading-spinner loading-xs';
            badge.append(spinner, document.createTextNode('PROCESSING'));

            return badge;
        }

        if (status === 'DONE') {
            badge.classList.add('badge-success');
            badge.textContent = 'DONE';

            return badge;
        }

        if (status === 'ERROR') {
            badge.classList.add('badge-error');
            badge.textContent = 'ERROR';

            return badge;
        }

        badge.classList.add('badge-ghost');
        badge.textContent = 'READY';

        return badge;
    }

    private createQueueItem(fileName: string, source: string, target: string): QueueItem {
        return {
            id: this.generateId(),
            fileName,
            source,
            target,
            status: 'READY',
            errorMessage: '',
            outputFileName: '',
        };
    }

    private generateId(): string {
        if (typeof window.crypto?.randomUUID === 'function') {
            return window.crypto.randomUUID();
        }

        return `${Date.now()}-${Math.random().toString(16).slice(2)}`;
    }

    private detectSourceFromFileName(fileName: string): string | null {
        const normalizedFileName = fileName.trim().toLowerCase();
        const extensionSeparatorIndex = normalizedFileName.lastIndexOf('.');
        if (extensionSeparatorIndex < 0 || extensionSeparatorIndex === normalizedFileName.length - 1) {
            return null;
        }

        const rawExtension = normalizedFileName.slice(extensionSeparatorIndex + 1);
        const extensionAliases: Record<string, string> = {
            jpeg: 'jpg',
            tif: 'tiff',
        };
        const normalizedExtension = extensionAliases[rawExtension] ?? rawExtension;

        return this.isSupportedSource(normalizedExtension) ? normalizedExtension : null;
    }

    private isSupportedSource(format: string): boolean {
        return Object.prototype.hasOwnProperty.call(this.formatsBySourceValue, format);
    }

    private resolveTargetsForSource(source: string): string[] {
        if (source === '' || typeof this.formatsBySourceValue !== 'object' || this.formatsBySourceValue === null) {
            return [];
        }

        const targets = this.formatsBySourceValue[source];

        return Array.isArray(targets) ? targets : [];
    }

    private resolvePreferredTarget(source: string, candidate: string): string {
        const normalizedCandidate = this.normalizeFormat(candidate);
        if (normalizedCandidate === '') {
            return '';
        }

        const targets = this.resolveTargetsForSource(source).map((target) => this.normalizeFormat(target));

        return targets.includes(normalizedCandidate) ? normalizedCandidate : '';
    }

    private currentSource(): string {
        const fromSelect = document.querySelector<HTMLSelectElement>('[data-conversion-selector-target="fromSelect"]');
        if (fromSelect !== null) {
            return this.normalizeFormat(fromSelect.value);
        }

        if (this.hasSelectedFromValue) {
            return this.normalizeFormat(this.selectedFromValue);
        }

        return '';
    }

    private currentPageTarget(): string {
        const toSelect = document.querySelector<HTMLSelectElement>('[data-conversion-selector-target="toSelect"]');
        if (toSelect !== null) {
            return this.normalizeFormat(toSelect.value);
        }

        return '';
    }

    private reconcileItemsWithCurrentSelection(): void {
        let hasChanges = false;

        this.items = this.items.map((item) => {
            const normalizedTarget = this.resolvePreferredTarget(item.source, item.target);
            if (normalizedTarget === item.target) {
                return item;
            }

            hasChanges = true;

            return {
                ...item,
                target: normalizedTarget,
                status: 'READY',
                errorMessage: '',
                outputFileName: '',
            };
        });

        if (hasChanges) {
            this.persistItems();
        }
    }

    private hydrateDoneStateFromDownloadStore(): void {
        this.items = this.items.map((item) => {
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

    private syncTransientStoresWithItems(): void {
        const activeItemIds = new Set(this.items.map((item) => item.id));

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

    private pruneItemsWithoutSourceData(): number {
        if (this.items.length === 0) {
            return 0;
        }

        const removedItemIds = this.items
            .filter((item) => !ORIGINAL_FILE_STORE.has(item.id))
            .map((item) => item.id);

        if (removedItemIds.length === 0) {
            return 0;
        }

        this.items = this.items.filter((item) => ORIGINAL_FILE_STORE.has(item.id));
        removedItemIds.forEach((itemId) => {
            this.revokeDownload(itemId);
        });
        this.persistItems();

        return removedItemIds.length;
    }

    private loadPersistedItems(): QueueItem[] {
        const rawState = this.readPersistedState();
        if (rawState === null || rawState === '') {
            return [];
        }

        let decodedState: unknown = null;
        try {
            decodedState = JSON.parse(rawState);
        } catch {
            this.clearPersistedState();

            return [];
        }

        if (!Array.isArray(decodedState)) {
            this.clearPersistedState();

            return [];
        }

        const parsedItems: QueueItem[] = [];
        decodedState.forEach((value) => {
            if (!this.isPersistedQueueItem(value)) {
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

    private persistItems(): void {
        if (this.items.length === 0) {
            this.clearPersistedState();

            return;
        }

        const persistedItems: PersistedQueueItem[] = this.items.map((item) => ({
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

    private readPersistedState(): string | null {
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

    private clearPersistedState(): void {
        try {
            window.sessionStorage.removeItem(UPLOAD_QUEUE_STORAGE_KEY);
        } catch {
        }
        try {
            window.localStorage.removeItem(UPLOAD_QUEUE_STORAGE_KEY);
        } catch {
        }
    }

    private isPersistedQueueItem(value: unknown): value is PersistedQueueItem {
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

    private showError(message: string): void {
        if (!this.hasErrorTarget) {
            return;
        }

        if (this.hasErrorMessageTarget) {
            this.errorMessageTarget.textContent = message;
        } else {
            this.errorTarget.textContent = message;
        }
        this.errorTarget.classList.remove('hidden');
    }

    private clearError(): void {
        if (!this.hasErrorTarget) {
            return;
        }

        if (this.hasErrorMessageTarget) {
            this.errorMessageTarget.textContent = '';
        } else {
            this.errorTarget.textContent = '';
        }
        this.errorTarget.classList.add('hidden');
    }

    private setItemError(itemId: string, message: string): void {
        const itemIndex = this.items.findIndex((item) => item.id === itemId);
        if (itemIndex < 0) {
            return;
        }

        this.items[itemIndex] = {
            ...this.items[itemIndex],
            status: 'ERROR',
            errorMessage: message,
            outputFileName: '',
        };
        this.revokeDownload(itemId);
        this.render();
    }

    private removeQueueItem(itemId: string): void {
        const nextItems = this.items.filter((item) => item.id !== itemId);
        if (nextItems.length === this.items.length) {
            return;
        }

        this.items = nextItems;
        ORIGINAL_FILE_STORE.delete(itemId);
        this.revokeDownload(itemId);
        this.persistItems();
        this.render();
    }

    private setDownload(itemId: string, fileName: string, mimeType: string, blob: Blob): void {
        this.revokeDownload(itemId);
        const objectUrl = URL.createObjectURL(blob);
        DOWNLOAD_STORE.set(itemId, {
            objectUrl,
            fileName,
            mimeType,
        });
    }

    private revokeDownload(itemId: string): void {
        const download = DOWNLOAD_STORE.get(itemId);
        if (download === undefined) {
            return;
        }

        URL.revokeObjectURL(download.objectUrl);
        DOWNLOAD_STORE.delete(itemId);
    }

    private extractErrorMessage(payload: unknown, fallback: string): string {
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

    private isConvertSuccessPayload(payload: unknown): payload is ConvertSuccessPayload {
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

    private decodeBase64ToBlob(contentBase64: string, mimeType: string): Blob {
        const binary = window.atob(contentBase64);
        const bytes = new Uint8Array(binary.length);
        for (let i = 0; i < binary.length; i += 1) {
            bytes[i] = binary.charCodeAt(i);
        }

        return new Blob([bytes], { type: mimeType });
    }

    private buildUnsupportedFilesMessage(unsupportedFiles: string[]): string {
        if (unsupportedFiles.length === 0) {
            return 'This file format is not supported yet.';
        }

        return `Unsupported file format for: ${unsupportedFiles.join(', ')}`;
    }

    private normalizeFormat(value: unknown): string {
        return String(value ?? '').trim().toLowerCase();
    }

    private buildSourcePageUrl(source: string): string {
        return this.sourcePageTemplateValue.replace('sourceplaceholder', encodeURIComponent(source));
    }

    private visit(path: string): void {
        const destination = new URL(path, window.location.origin);
        const current = new URL(window.location.href);
        const destinationPath = `${destination.pathname}${destination.search}`;
        const currentPath = `${current.pathname}${current.search}`;
        if (destinationPath === currentPath) {
            return;
        }

        const turbo = (window as Window & { Turbo?: { visit: (location: string) => void } }).Turbo;
        if (turbo !== undefined) {
            turbo.visit(destinationPath);

            return;
        }

        window.location.assign(destination.toString());
    }
}
