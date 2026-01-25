import App from './App.svelte';
import '../public/app.css';
// import '@skeletonlabs/skeleton/themes/wintry.css';
// import '@skeletonlabs/skeleton/styles/all.css';

// Инициализация темы Skeleton
if (typeof document !== 'undefined') {
	document.documentElement.setAttribute('data-theme', 'wintry');
}

const app = new App({
	target: document.body
});

export default app;