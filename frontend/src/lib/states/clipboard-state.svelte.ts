export class ClipboardState {
    list = $state<string[]>([]);
    maxLength = $state(10);

    add(item: string): void {
        const value = item.trim();
        if (!value) return;

        // Use a new array to ensure reactivity triggers correctly
        let newList = this.list.filter((i) => i !== value);
        newList = [value, ...newList];

        if (newList.length > this.maxLength) {
            newList = newList.slice(0, this.maxLength);
        }
        this.list = newList;
    }

    setMaxLength(len: number): void {
        const n = Number.isFinite(len) ? Math.floor(len) : 0;
        this.maxLength = Math.max(0, n);
        if (this.list.length > this.maxLength) {
            this.list = this.list.slice(0, this.maxLength);
        }
    }
}

export const clipboardState = new ClipboardState();
