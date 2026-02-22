import type { QueueItem, QueueItemStatus } from './types.ts';

type BuildItemRowOptions = {
    isConverting: boolean;
    hasDownload: boolean;
    targets: string[];
};

export function buildItemRow(item: QueueItem, options: BuildItemRowOptions): HTMLLIElement {
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

    const statusBadge = buildStatusBadge(item.status);

    const label = document.createElement('span');
    label.className = 'text-sm text-base-content/70';
    label.textContent = 'Convert to';

    const targetSelect = document.createElement('select');
    targetSelect.className = 'select select-bordered select-sm min-w-28';
    targetSelect.dataset.itemId = item.id;
    targetSelect.setAttribute('data-action', 'change->upload-queue#onRowTargetChange');
    targetSelect.disabled = options.isConverting || item.status === 'PROCESSING';
    targetSelect.add(new Option('to', '', item.target === ''));
    options.targets.forEach((target) => {
        targetSelect.add(new Option(target.toUpperCase(), target, false, target === item.target));
    });

    const removeButton = document.createElement('button');
    removeButton.type = 'button';
    removeButton.className = 'btn btn-ghost btn-sm text-error';
    removeButton.dataset.itemId = item.id;
    removeButton.disabled = options.isConverting || item.status === 'PROCESSING';
    removeButton.setAttribute('data-action', 'click->upload-queue#removeFile');
    removeButton.textContent = 'X';

    controls.append(statusBadge, label, targetSelect);

    if (item.status === 'DONE' && options.hasDownload) {
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

function buildStatusBadge(status: QueueItemStatus): HTMLSpanElement {
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
