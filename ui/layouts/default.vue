<template>
	<div class="card">
		<nuxt-link to="/">Home Page</nuxt-link>
		<nuxt />
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
