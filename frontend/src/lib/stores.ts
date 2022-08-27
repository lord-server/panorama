import { writable, type Writable } from 'svelte/store';

export interface WorldCursorPosition {
    x: number,
    z: number,
}

export const worldPositionUnderCursor: Writable<WorldCursorPosition | null> = writable(null);
