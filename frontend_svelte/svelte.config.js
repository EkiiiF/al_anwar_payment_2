import adapter from '@sveltejs/adapter-node';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  preprocess: vitePreprocess(),

  compilerOptions: {
    runes: ({ filename }) =>
      filename.split(/[/\\]/).includes('node_modules') ? undefined : true
  },

  kit: {
    adapter: adapter(),
    csrf: {
      checkOrigin: true // Keamanan: hanya izinkan request dari origin yang sah
    },
    alias: {
      // Clean Architecture imports
      $api:        './src/lib/api',
      $components: './src/lib/components',
      $types:      './src/lib/types',
      $utils:      './src/lib/utils'
    }
  }
};

export default config;