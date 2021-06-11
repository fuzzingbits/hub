<template>
	<div>
		<transition name="fade">
			<div v-if="loading" class="loading">
				<i class="fas fa-circle-notch fa-spin"></i>
			</div>
		</transition>
		<div v-if="showSplash" id="splash">
			<PosterMessage :poster="userPoster" />
			<PosterMessage :poster="serverPoster" />
		</div>
	</div>
</template>

<script lang="ts">
import Vue from 'vue';
import HubApi from '~/ui/assets/api';
import Poster from '~/ui/assets/poster';

export default Vue.extend({
	data: function() {
		return {
			userDataLoadedOnce: false,
			serverDataLoadedOnce: false,
			userPoster: new Poster(true),
			serverPoster: new Poster(true),
		};
	},
	computed: {
		loading: function(): boolean {
			// If it's not the first page load
			if (this.userDataLoadedOnce && this.serverDataLoadedOnce) {
				return false;
			}

			// If any are loading
			if (this.userPoster.loading || this.serverPoster.loading) {
				return true;
			}

			return false;
		},
		showSplash: function(): boolean {
			if (!this.userPoster.state || !this.serverPoster.state) {
				return true;
			}

			return false;
		},
	},
	methods: {
		checkForLogin() {
			this.userPoster.reset(true);
			HubApi.userMe()
				.then(response => {
					this.$store.commit('user/setState', response.data);
					this.userPoster.setResponse(response);
					this.userDataLoadedOnce = true;
				})
				.catch(err => {
					this.userPoster.handlerError(err);
				});
		},
		checkServerStatus() {
			this.serverPoster.reset(true);
			HubApi.serverStatus()
				.then(response => {
					this.$store.commit('server/setStatus', response.data);
					this.serverPoster.setResponse(response);
					this.serverDataLoadedOnce = true;

					if (response.data && response.data.setupRequired) {
						this.$router.push('/setup');
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
#splash {
	padding: var(--spacing);
}

.loading {
	align-items: center;
	background: var(--background);
	bottom: 0;
	display: flex;
	font-size: 2.5rem;
	justify-content: center;
	left: 0;
	position: fixed;
	right: 0;
	top: 0;
	z-index: 100;
}
</style>
