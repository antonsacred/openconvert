import { Controller } from '@hotwired/stimulus';

type FormatsBySource = Record<string, string[]>;
type QueueItem = {
    id: string;
    fileName: string;
    source: string;
    target: string;
};

const UPLOAD_QUEUE_STORAGE_KEY = 'openconvert.upload_queue_state';

export default class extends Controller<HTMLElement> {
    static targets = ['fileInput', 'fileList', 'emptyState', 'queuePanel', 'error', 'errorMessage', 'convertButton'] as const;
    static values = {
        formatsBySource: Object,
        sourcePageTemplate: String,
        selectedFrom: String,
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

    private items: QueueItem[] = [];

    connect(): void {
        this.clearError();
        this.items = this.loadPersistedItems();
        this.reconcileItemsWithCurrentSelection();
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

        const supportedFiles: Array<{ fileName: string; source: string }> = [];
        const unsupportedFiles: string[] = [];
        selectedFiles.forEach((file) => {
            const detectedSource = this.detectSourceFromFileName(file.name);
            if (detectedSource === null) {
                unsupportedFiles.push(file.name);

                return;
            }

            supportedFiles.push({
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
        const queueItemsToAdd = supportedFiles.map((file) => {
            const defaultTarget = this.resolvePreferredTarget(file.source, pageTarget);

            return this.createQueueItem(file.fileName, file.source, defaultTarget);
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

        this.items = this.items.filter((item) => item.id !== itemId);
        this.persistItems();
        this.render();
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
        };

        this.persistItems();
        this.render();
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
            this.convertButtonTarget.disabled = !this.items.every((item) => item.target !== '');
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

        const controls = document.createElement('div');
        controls.className = 'flex shrink-0 items-center gap-2 whitespace-nowrap';

        const label = document.createElement('span');
        label.className = 'text-sm text-base-content/70';
        label.textContent = 'Convert to';

        const targetSelect = document.createElement('select');
        targetSelect.className = 'select select-bordered select-sm min-w-28';
        targetSelect.dataset.itemId = item.id;
        targetSelect.setAttribute('data-action', 'change->upload-queue#onRowTargetChange');
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
        removeButton.setAttribute('data-action', 'click->upload-queue#removeFile');
        removeButton.textContent = 'X';

        controls.append(label, targetSelect, removeButton);
        row.append(fileInfo, controls);

        return row;
    }

    private createQueueItem(fileName: string, source: string, target: string): QueueItem {
        return {
            id: this.generateId(),
            fileName,
            source,
            target,
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
            const normalizedExistingTarget = this.resolvePreferredTarget(item.source, item.target);
            const nextTarget = normalizedExistingTarget;

            if (nextTarget === item.target) {
                return item;
            }

            hasChanges = true;

            return {
                ...item,
                target: nextTarget,
            };
        });

        if (hasChanges) {
            this.persistItems();
        }
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
        decodedState.forEach((item) => {
            if (!this.isQueueItem(item)) {
                return;
            }

            parsedItems.push({
                id: item.id,
                fileName: item.fileName,
                source: item.source,
                target: item.target,
            });
        });

        return parsedItems;
    }

    private persistItems(): void {
        if (this.items.length === 0) {
            this.clearPersistedState();

            return;
        }

        const encoded = JSON.stringify(this.items);
        try {
            window.sessionStorage.setItem(UPLOAD_QUEUE_STORAGE_KEY, encoded);
        } catch {
            try {
                window.localStorage.setItem(UPLOAD_QUEUE_STORAGE_KEY, encoded);
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

    private isQueueItem(value: unknown): value is QueueItem {
        if (typeof value !== 'object' || value === null) {
            return false;
        }

        const item = value as Record<string, unknown>;

        return typeof item.id === 'string'
            && item.id.trim() !== ''
            && typeof item.fileName === 'string'
            && item.fileName.trim() !== ''
            && typeof item.source === 'string'
            && item.source.trim() !== ''
            && typeof item.target === 'string';
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
