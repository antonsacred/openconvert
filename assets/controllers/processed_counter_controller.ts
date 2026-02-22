import { Controller } from '@hotwired/stimulus';

export default class extends Controller<HTMLElement> {
    static targets = ['count'] as const;
    static values = {
        start: Number,
    } as const;

    declare readonly countTarget: HTMLElement;
    declare readonly hasCountTarget: boolean;
    declare readonly startValue: number;
    declare readonly hasStartValue: boolean;

    private currentValue = 0;
    private timer: number | null = null;

    connect(): void {
        this.currentValue = this.hasStartValue
            ? Math.floor(this.startValue)
            : Math.floor(Date.now() / 1000 / 2);
        this.render();

        this.timer = window.setInterval(() => {
            this.currentValue += 1;
            this.render();
        }, 2000);
    }

    disconnect(): void {
        if (this.timer !== null) {
            window.clearInterval(this.timer);
            this.timer = null;
        }
    }

    private render(): void {
        if (!this.hasCountTarget) {
            return;
        }

        this.countTarget.textContent = new Intl.NumberFormat('en-US').format(this.currentValue);
    }
}
