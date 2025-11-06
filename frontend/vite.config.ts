import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
    watch: {
      usePolling: true, // Enable polling for file changes
      interval: 100,    // Check for changes every 100ms
    },
    host: true,          // Allow access from outside the container
    strictPort: true,    // Exit if the port is already in use
  },
});
