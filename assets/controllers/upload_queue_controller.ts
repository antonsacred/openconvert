import { Controller } from '@hotwired/stimulus';
import { decodeBase64ToBlob, extractErrorMessage, isConvertSuccessPayload, submitConversionRequest } from './upload_queue/api.ts';
import {
    deleteOriginalFile,
    getDownload,
    getOriginalFile,
    hasDownload,
    hasOriginalFile,
    hydrateDoneStateFromDownloads,
    persistQueueItems,
    readPersistedQueueItems,
    revokeDownload,
    setDownload,
    setOriginalFile,
    syncTransientStoresWithItems,
} from './upload_queue/storage.ts';
import type { FormatsBySource, QueueItem } from './upload_queue/types.ts';
import { buildItemRow } from './upload_queue/view.ts';

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
        this.items = readPersistedQueueItems();
        this.pruneItemsWithoutSourceData();
        this.reconcileItemsWithCurrentSelection();
        this.items = hydrateDoneStateFromDownloads(this.items);
        syncTransientStoresWithItems(this.items);
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
            setOriginalFile(queueItem.id, supportedFile.file);
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
        revokeDownload(item.id);

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

        const download = getDownload(itemId);
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

        const sourceFile = getOriginalFile(itemId);
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
        revokeDownload(itemId);
        this.render();

        try {
            const response = await submitConversionRequest(this.convertUrlValue, {
                from: item.source,
                to: item.target,
                fileName: item.fileName,
                file: sourceFile,
            });

            if (!response.ok) {
                this.setItemError(itemId, extractErrorMessage(response.payload, `Failed to convert ${item.fileName}.`));

                return;
            }

            if (!isConvertSuccessPayload(response.payload)) {
                this.setItemError(itemId, `Invalid conversion response for ${item.fileName}.`);

                return;
            }

            const outputBlob = decodeBase64ToBlob(response.payload.contentBase64, response.payload.mimeType);
            setDownload(itemId, response.payload.fileName, response.payload.mimeType, outputBlob);

            const doneIndex = this.items.findIndex((queueItem) => queueItem.id === itemId);
            if (doneIndex >= 0) {
                this.items[doneIndex] = {
                    ...this.items[doneIndex],
                    status: 'DONE',
                    errorMessage: '',
                    outputFileName: response.payload.fileName,
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
            const targets = this.resolveTargetsForSource(item.source)
                .map((target) => this.normalizeFormat(target))
                .filter((target) => target !== '');

            fragment.appendChild(buildItemRow(item, {
                isConverting: this.isConverting,
                hasDownload: hasDownload(item.id),
                targets,
            }));
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

    private pruneItemsWithoutSourceData(): number {
        if (this.items.length === 0) {
            return 0;
        }

        const removedItemIds = this.items
            .filter((item) => !hasOriginalFile(item.id))
            .map((item) => item.id);

        if (removedItemIds.length === 0) {
            return 0;
        }

        this.items = this.items.filter((item) => hasOriginalFile(item.id));
        removedItemIds.forEach((itemId) => {
            revokeDownload(itemId);
        });
        this.persistItems();

        return removedItemIds.length;
    }

    private persistItems(): void {
        persistQueueItems(this.items);
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
        revokeDownload(itemId);
        this.render();
    }

    private removeQueueItem(itemId: string): void {
        const nextItems = this.items.filter((item) => item.id !== itemId);
        if (nextItems.length === this.items.length) {
            return;
        }

        this.items = nextItems;
        deleteOriginalFile(itemId);
        revokeDownload(itemId);
        this.persistItems();
        this.render();
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
