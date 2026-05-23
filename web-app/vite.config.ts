import { resolve } from 'path';
import { purgeCss } from 'vite-plugin-tailwind-purgecss';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

const themeAliases = [
	'crimson', 'gold-nouveau', 'hamlindigo', 'modern', 'rocket',
	'sahara', 'seafoam', 'skeleton', 'vintage', 'wintry'
].reduce(
	(acc, name) => {
		acc[`@skeletonlabs/skeleton/themes/${name}`] = resolve(__dirname, 'src/lib/empty-theme.ts');
		return acc;
	},
	{} as Record<string, string>
);

export default defineConfig({
	plugins: [sveltekit(), purgeCss()],
	resolve: {
		alias: themeAliases
	},
	server: {
		host: true,
		port: 3000
	}
});