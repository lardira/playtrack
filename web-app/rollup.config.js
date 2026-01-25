import { spawn } from 'child_process';
import svelte from 'rollup-plugin-svelte';
import commonjs from '@rollup/plugin-commonjs';
import terser from '@rollup/plugin-terser';
import resolve from '@rollup/plugin-node-resolve';
import alias from '@rollup/plugin-alias';
import livereload from 'rollup-plugin-livereload';
import css from 'rollup-plugin-css-only';
import postcss from 'rollup-plugin-postcss';
import sveltePreprocess from 'svelte-preprocess';
import typescriptPlugin from '@rollup/plugin-typescript';
import path from 'path';
import { fileURLToPath } from 'url';
import { createRequire } from 'module';
import tailwindcss from '@tailwindcss/postcss';
import autoprefixer from 'autoprefixer';

const require = createRequire(import.meta.url);
const typescript = require('typescript');

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const production = !process.env.ROLLUP_WATCH;

function serve() {
	let server;

	function toExit() {
		if (server) server.kill(0);
	}

	return {
		writeBundle() {
			if (server) return;
			server = spawn('npm', ['run', 'start', '--', '--dev'], {
				stdio: ['ignore', 'inherit', 'inherit'],
				shell: true
			});

			process.on('SIGTERM', toExit);
			process.on('exit', toExit);
		}
	};
}

export default {
	input: 'src/main.js',
	output: {
		sourcemap: true,
		format: 'iife',
		name: 'app',
		file: 'public/build/bundle.js'
	},
	plugins: [
		alias({
			entries: [
				{ find: '$lib', replacement: path.resolve(__dirname, 'src/lib') }
			]
		}),
		typescriptPlugin({
			tsconfig: './tsconfig.json',
			sourceMap: !production,
			inlineSources: !production,
			rootDir: './src',
			exclude: ['node_modules/**'],
		}),
		svelte({
			preprocess: sveltePreprocess({
				typescript: {
					tsconfigFile: './tsconfig.json',
					compilerOptions: {
						module: typescript.ModuleKind.ESNext,
						target: typescript.ScriptTarget.ES2020,
					},
				},
				postcss: {
					plugins: [
						tailwindcss,
						autoprefixer,
					],
				},
			}),
			compilerOptions: {
				// enable run-time checks when not in production
				dev: !production
			}
		}),
		// we'll extract any component CSS out into
		// a separate file - better for performance
		css({ output: 'bundle.css' }),
		// Process global CSS files with PostCSS
		postcss({
			extract: path.resolve(__dirname, 'public/build/app.css'),
			minify: production,
			include: ['**/*.css', '**/app.css'],
		}),

		// If you have external dependencies installed from
		// npm, you'll most likely need these plugins. In
		// some cases you'll need additional configuration -
		// consult the documentation for details:
		// https://github.com/rollup/plugins/tree/master/packages/commonjs
		resolve({
			browser: true,
			dedupe: ['svelte'],
			exportConditions: ['svelte']
		}),
		commonjs(),

		// In dev mode, call `npm run start` once
		// the bundle has been generated
		!production && serve(),

		// Watch the `public` directory and refresh the
		// browser on changes when not in production
		!production && livereload('public'),

		// If we're building for production (npm run build
		// instead of npm run dev), minify
		production && terser()
	],
	watch: {
		clearScreen: false
	}
};
