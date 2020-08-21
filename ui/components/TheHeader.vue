<template>
	<div id="header">
		<nuxt-link v-if="setupRequired" to="/setup">Setup</nuxt-link>
		<nuxt-link v-if="!setupRequired" to="/">Home Page</nuxt-link>
		<nuxt-link v-if="!setupRequired" to="/about">About</nuxt-link>
		<nuxt-link v-if="!session && !setupRequired" to="/login">Login</nuxt-link>
		<nuxt-link v-if="session && !setupRequired" to="/profile">Profile</nuxt-link>
		<a v-if="session && !setupRequired" @click.prevent="logout" href="#">Logout</a>
	</div>
</template>

<script lang="ts">
import Vue from "vue";
import * as types from "~/ui/assets/types";

export default Vue.extend({
	computed: {
		session: function(): types.UserContext | null {
			return this.$store.state.user.session;
		},
		serverStatus: function(): types.ServerStatus | null {
			return this.$store.state.server.status;
		},
		setupRequired: function(): boolean {
			const serverStatus = this.serverStatus;
			if (!serverStatus) {
				return false;
			}

			return serverStatus.setupRequired;
		},
	},
});
</script>

<style>
#header a {
	padding: var(--spacing-half) var(--spacing);
}

#header {
	background: var(--background-accent);
	padding: var(--spacing-double);
}
</style>
