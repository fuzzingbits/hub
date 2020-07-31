<template>
	<div class="card">
		<nuxt-link to="/">Home Page</nuxt-link>
		<nuxt-link to="/about">About</nuxt-link>
		<nuxt-link v-if="!session" to="/login">Login</nuxt-link>
		<nuxt />
		<p>
			<a href="https://github.com/fuzzingbits/hub" target="_blank" rel="noreferrer">
				Source Code
			</a>
		</p>
	</div>
</template>

<script lang="ts">
import Vue from "vue";
import HubApi from "~/ui/assets/api";
import * as types from "~/ui/assets/types";
export default Vue.extend({
	computed: {
		session: function(): types.UserContext | null {
			return this.$store.state.user.session;
		},
	},
	methods: {
		checkForLogin() {
			HubApi.userMe()
				.then(response => {
					this.$store.commit("user/setState", response.data);
				})
				.catch(err => {
					console.error("ajax error: " + err);
				});
		},
	},
	mounted() {
		this.checkForLogin();
	},
});
</script>

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
