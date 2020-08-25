<template>
	<div class="card">
		<form @submit.prevent="submit" id="page-form">
			<h2>Server Setup</h2>
			<p>This is for first time setups only.</p>
			<label>First Name <input name="firstName" required/></label>
			<label>Last Name <input name="lastName" required/></label>
			<label>Email <input name="email" type="email" required/></label>
			<label>Password <input name="password" type="password" required/></label>
			<PosterMessage :poster="formPoster" />
			<label><input type="submit" value="Complete Setup"/></label>
		</form>
	</div>
</template>

<script lang="ts">
import Vue from "vue";
import HubApi from "~/ui/assets/api";
import * as types from "~/ui/assets/types";
import Poster from "~/ui/assets/poster";

export default Vue.extend({
	data: function() {
		return {
			formPoster: new Poster(),
		};
	},
	computed: {
		serverStatus: function(): types.ServerStatus | null {
			return this.$store.state.server.status;
		},
	},
	methods: {
		submit(): void {
			this.formPoster.reset(true);

			const form = document.querySelector("#page-form") as HTMLFormElement;
			const formData = new FormData(form);

			HubApi.serverSetup({
				firstName: formData.get("firstName") as string,
				lastName: formData.get("lastName") as string,
				email: formData.get("email") as string,
				password: formData.get("password") as string,
			})
				.then(response => {
					this.formPoster.setResponse(response);
					if (!response.state) {
						return;
					}

					// Login the new user
					this.$store.commit("user/setState", response.data);

					// Update the server status
					let serverStatus: types.ServerStatus = {
						setupRequired: false,
					};
					this.$store.commit("server/setStatus", serverStatus);

					// Redirect to the home page
					this.$router.push("/");
				})
				.catch(err => {
					this.formPoster.handlerError(err);
				});
		},
	},
});
</script>
