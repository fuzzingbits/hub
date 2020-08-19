<template>
	<div>
		<div v-if="showSplash" id="splash">
			<i class="fas fa-circle-notch fa-spin"></i>
		</div>
		<div v-if="!showSplash" id="page">
			<div id="header">
				<nuxt-link v-if="setupRequired" to="/setup">Setup</nuxt-link>
				<nuxt-link v-if="!setupRequired" to="/">Home Page</nuxt-link>
				<nuxt-link v-if="!setupRequired" to="/about">About</nuxt-link>
				<nuxt-link v-if="!session && !setupRequired" to="/login">Login</nuxt-link>
			</div>
			<div id="body">
				<div class="card">
					<nuxt />
				</div>
			</div>
			<div id="footer">
				<p>
					<a href="https://github.com/fuzzingbits/hub" target="_blank" rel="noreferrer">
						Source Code
					</a>
				</p>
			</div>
		</div>
	</div>
</template>

<script lang="ts">
import Vue from "vue";
import HubApi from "~/ui/assets/api";
import * as types from "~/ui/assets/types";
export default Vue.extend({
	data: function() {
		return {
			loadingUser: true,
			loadingServer: true,
		};
	},
	computed: {
		showSplash: function(): boolean {
			return this.loading;
		},
		loading: function(): boolean {
			return this.loadingUser || this.loadingServer;
		},
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
	methods: {
		checkForLogin() {
			HubApi.userMe()
				.then(response => {
					this.$store.commit("user/setState", response.data);
					this.loadingUser = false;
				})
				.catch(err => {
					console.error("ajax error: " + err);
				})
				.finally(() => {
					this.loadingUser = false;
				});
		},
		checkServerStatus() {
			HubApi.serverStatus()
				.then(response => {
					this.$store.commit("server/setStatus", response.data);

					if (response.data && response.data.setupRequired) {
						this.$router.push("/setup");
					}
					this.loadingServer = false;
				})
				.catch(err => {
					console.error("ajax error: " + err);
				});
		},
	},
	mounted() {
		this.checkForLogin();
		this.checkServerStatus();
	},
});
</script>

<style>
@import url("../assets/prism/prism.css");
</style>
