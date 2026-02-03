import { MessageObject } from 'svelte-i18n/types/runtime/types';

declare global {
  namespace svelteI18n {
    interface InterpolationValues {
      [key: string]: string | number | boolean | Date | undefined | null;
    }
  }
}