import { Controller } from '@hotwired/stimulus';

type FormatsBySource = Record<string, string[]>;

export default class extends Controller<HTMLFormElement> {
    static targets = ['fromSelect', 'toSelect'] as const;
    static values = {
        formatsBySource: Object,
        homePageUrl: String,
        sourcePageTemplate: String,
        pairPageTemplate: String,
        selectedTo: String,
    } as const;

    declare readonly fromSelectTarget: HTMLSelectElement;
    declare readonly toSelectTarget: HTMLSelectElement;
    declare readonly formatsBySourceValue: FormatsBySource;
    declare readonly homePageUrlValue: string;
    declare readonly sourcePageTemplateValue: string;
    declare readonly pairPageTemplateValue: string;
    declare readonly selectedToValue: string;
    declare readonly hasSelectedToValue: boolean;

    connect(): void {
        const preferredTarget = this.hasSelectedToValue ? this.normalizeFormat(this.selectedToValue) : '';
        this.rebuildTargetOptions(preferredTarget);
    }

    onFromChange(): void {
        const preferredTarget = this.normalizeFormat(this.toSelectTarget.value);
        this.rebuildTargetOptions(preferredTarget);

        const source = this.normalizeFormat(this.fromSelectTarget.value);
        if (source === '') {
            this.visit(this.homePageUrlValue);

            return;
        }

        const target = this.normalizeFormat(this.toSelectTarget.value);
        if (target !== '') {
            this.visit(this.buildPairPageUrl(source, target));

            return;
        }

        this.visit(this.buildSourcePageUrl(source));
    }

    onToChange(): void {
        const source = this.normalizeFormat(this.fromSelectTarget.value);
        const target = this.normalizeFormat(this.toSelectTarget.value);
        if (source === '') {
            this.visit(this.homePageUrlValue);

            return;
        }

        if (target === '') {
            this.visit(this.buildSourcePageUrl(source));

            return;
        }

        const targets = this.resolveTargetsForSource(source).map((item) => this.normalizeFormat(item));
        if (!targets.includes(target)) {
            return;
        }

        this.visit(this.buildPairPageUrl(source, target));
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

    private buildSourcePageUrl(source: string): string {
        return this.sourcePageTemplateValue.replace('sourceplaceholder', encodeURIComponent(source));
    }

    private buildPairPageUrl(source: string, target: string): string {
        return this.pairPageTemplateValue
            .replace('sourceplaceholder', encodeURIComponent(source))
            .replace('targetplaceholder', encodeURIComponent(target));
    }

    private visit(path: string): void {
        const turbo = (window as Window & { Turbo?: { visit: (location: string) => void } }).Turbo;
        if (turbo !== undefined) {
            turbo.visit(path);

            return;
        }

        window.location.assign(path);
    }
}
