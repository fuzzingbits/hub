<template>
	<div>
		<transition name="fade">
			<div v-if="showSplash" id="splash">
				<div v-if="userPoster.loading || serverPoster.loading" class="loading">
					<i class="fas fa-circle-notch fa-spin"></i>
				</div>
				<PosterMessage :poster="userPoster" />
				<PosterMessage :poster="serverPoster" />
			</div>
		</transition>
		<div v-if="!showSplash" id="page">
			<div id="header">
				<nuxt-link v-if="setupRequired" to="/setup">Setup</nuxt-link>
				<nuxt-link v-if="!setupRequired" to="/">Home Page</nuxt-link>
				<nuxt-link v-if="!setupRequired" to="/about">About</nuxt-link>
				<nuxt-link v-if="!session && !setupRequired" to="/login">Login</nuxt-link>
				<nuxt-link v-if="session && !setupRequired" to="/profile">Profile</nuxt-link>
				<a v-if="session && !setupRequired" @click.prevent="logout" href="#">Logout</a>
			</div>
			<div id="body">
				<nuxt />
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
import Poster from "~/ui/assets/poster";
import * as types from "~/ui/assets/types";
export default Vue.extend({
	data: function() {
		return {
			userPoster: new Poster(),
			serverPoster: new Poster(),
		};
	},
	computed: {
		showSplash: function(): boolean {
			if (this.userPoster.loading || this.serverPoster.loading) {
				return true;
			}

			if (!this.userPoster.state || !this.serverPoster.state) {
				return true;
			}

			return false;
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
		logout() {
			HubApi.userLogout();
			this.$store.commit("user/setState", null);
		},
		checkForLogin() {
			this.userPoster.reset(true);
			HubApi.userMe()
				.then(response => {
					this.$store.commit("user/setState", response.data);
					this.userPoster.setResponse(response);
				})
				.catch(err => {
					this.userPoster.handlerError(err);
				});
		},
		checkServerStatus() {
			HubApi.serverStatus()
				.then(response => {
					this.$store.commit("server/setStatus", response.data);
					this.serverPoster.setResponse(response);

					if (response.data && response.data.setupRequired) {
						this.$router.push("/setup");
					}
				})
				.catch(err => {
					this.serverPoster.handlerError(err);
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
