<template>
	<div class="card">
		<nuxt-link to="/">Home Page</nuxt-link>
		<nuxt-link to="/about">About</nuxt-link>
		<nuxt />
		<p>
			<a href="https://github.com/fuzzingbits/hub" target="_blank" rel="noreferrer">
				Source Code
			</a>
		</p>
	</div>
</template>

<style>
@import url("../assets/prism/prism.css");

body {
	margin: var(--spacing);
}

.card {
	max-width: 768px;
	background: var(--background-accent);
	margin: var(--spacing);
	padding: var(--spacing-double);
	margin: auto;
}
</style>
