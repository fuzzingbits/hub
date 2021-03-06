import { NuxtConfig } from '@nuxt/types';

const config: NuxtConfig = {
	buildDir: 'var/.nuxt',
	ssr: false,
	generate: {
		dir: 'resources/dist',
		fallback: '200.html',
		exclude: [/^.*$/],
	},
	dir: {
		layouts: 'ui/layouts',
		middleware: 'ui/middleware',
		pages: 'ui/pages',
		static: 'ui/static',
		store: 'ui/store',
	},
	build: {
		extractCSS: true,
		publicPath: '/nuxt-build/',
	},
	components: [
		'~/ui/components',
	],
	buildModules: ['@nuxt/typescript-build'],
	head: {
		title: 'Hub',
		htmlAttrs: {
			lang: 'en',
		},
		link: [
			{
				rel: 'icon',
				href: '/favicon.svg',
				type: 'image/svg+xml',
			},
			{
				rel: 'apple-touch-icon',
				href: '/apple-touch-icon.png',
			},
		],
		meta: [
			{ charset: 'utf-8' },
			{ name: 'viewport', content: 'width=device-width, initial-scale=1' },
			{
				name: 'description',
				content: 'A productivity center for just you or a small group.',
			},
		],
	},
};

export default config;
