<template>
	<div class="card">
		<form v-if="session" @submit.prevent="submit" id="profile-form">
			<h2>Update Profile</h2>
			<p>This is for users to update their profile.</p>
			<input type="hidden" name="uuid" :value="session.user.uuid" />
			<label>First Name <input :value="session.user.firstName" name="firstName" required/></label>
			<label>Last Name <input :value="session.user.lastName" name="lastName" required/></label>
			<label>Email <input :value="session.user.email" name="email" required/></label>
			<label>Theme Color <input :value="session.userSettings.themeColor" name="themeColor" type="color" @input="changeColor" required/></label>
			<PosterMessage :poster="formPoster" />
			<label><input type="submit" value="Submit Update"/></label>
		</form>
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
			formPoster: new Poster(),
		};
	},
	computed: {
		session: function(): types.UserContext | null {
			return this.$store.state.user.session;
		},
	},
	methods: {
		changeColor(e: any) {
			document.documentElement.style.setProperty("--primary", e.target.value);
		},
		submit(): void {
			this.formPoster.reset();

			const form = document.querySelector("#profile-form") as HTMLFormElement;
			const formData = new FormData(form);

			HubApi.userUpdate({
				uuid: formData.get("uuid") as string,
				firstName: formData.get("firstName") as string,
				lastName: formData.get("lastName") as string,
				email: formData.get("email") as string,
				themeColor: formData.get("themeColor") as string,
			})
				.then(response => {
					this.formPoster.setResponse(response);
					if (!response.state) {
						return;
					}

					// Login the user
					this.$store.commit("user/setState", response.data);
				})
				.catch(err => {
					this.formPoster.handlerError(err);
				});
		},
	},
});
</script>
