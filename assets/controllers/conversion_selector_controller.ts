import { Controller } from '@hotwired/stimulus';

type FormatsBySource = Record<string, string[]>;

export default class extends Controller<HTMLFormElement> {
    static targets = ['fromSelect', 'toSelect'] as const;
    static values = {
        formatsBySource: Object,
        selectedTo: String,
    } as const;

    declare readonly fromSelectTarget: HTMLSelectElement;
    declare readonly toSelectTarget: HTMLSelectElement;
    declare readonly formatsBySourceValue: FormatsBySource;
    declare readonly selectedToValue: string;
    declare readonly hasSelectedToValue: boolean;

    connect(): void {
        const preferredTarget = this.hasSelectedToValue ? this.normalizeFormat(this.selectedToValue) : '';
        this.rebuildTargetOptions(preferredTarget);
    }

    onFromChange(): void {
        this.rebuildTargetOptions('');
    }

    rebuildTargetOptions(preferredTarget: string): void {
        const toSelect = this.toSelectTarget;
        const fromFormat = this.normalizeFormat(this.fromSelectTarget.value);
        const targets = this.resolveTargetsForSource(fromFormat);

        toSelect.innerHTML = '';
        toSelect.add(new Option('to', '', preferredTarget === ''));

        if (fromFormat === '' || targets.length === 0) {
            toSelect.disabled = true;
            toSelect.value = '';

            return;
        }

        let preferredTargetExists = false;

        targets.forEach((target) => {
            const normalizedTarget = this.normalizeFormat(target);
            if (normalizedTarget === '') {
                return;
            }

            const isSelected = normalizedTarget === preferredTarget;
            if (isSelected) {
                preferredTargetExists = true;
            }

            toSelect.add(new Option(normalizedTarget, normalizedTarget, false, isSelected));
        });

        toSelect.disabled = false;
        if (!preferredTargetExists) {
            toSelect.value = '';
        }
    }

    resolveTargetsForSource(source: string): string[] {
        if (source === '' || typeof this.formatsBySourceValue !== 'object' || this.formatsBySourceValue === null) {
            return [];
        }

        const bySource = this.formatsBySourceValue;
        const targets = bySource[source];

        return Array.isArray(targets) ? targets : [];
    }

    normalizeFormat(value: unknown): string {
        return String(value ?? '').trim().toLowerCase();
    }
}
