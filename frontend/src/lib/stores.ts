import { writable, type Writable } from 'svelte/store';

export interface WorldCursorPosition {
    x: number | null,
    z: number | null,
}

export const worldPositionUnderCursor: Writable<WorldCursorPosition> = writable({x: null, z: null});
