import { SvelteComponentTyped } from 'svelte';

export type IconComponent = typeof SvelteComponentTyped<any, any, any>;

export interface MenuItem {
  name: string;
  path: string;
  icon: IconComponent;
}
