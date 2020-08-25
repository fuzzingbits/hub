<template>
	<div class="card">
		<form @submit.prevent="submit" id="page-form">
			<h2>User Login</h2>
			<p>This is for users to login.</p>
			<label>Email <input name="email" type="email" required/></label>
			<label>Password <input name="password" type="password" required/></label>
			<PosterMessage :poster="formPoster" />
			<label><input type="submit" value="Login"/></label>
		</form>
	</div>
</template>

<script lang="ts">
import Vue from "vue";
import HubApi from "~/ui/assets/api";
import Poster from "~/ui/assets/poster";

export default Vue.extend({
	data: function() {
		return {
			formPoster: new Poster(),
		};
	},
	methods: {
		submit(): void {
			this.formPoster.reset(true);

			const form = document.querySelector("#page-form") as HTMLFormElement;
			const formData = new FormData(form);

			HubApi.userLogin({
				email: formData.get("email") as string,
				password: formData.get("password") as string,
			})
				.then(response => {
					this.formPoster.setResponse(response);
					if (!response.state) {
						return;
					}

					// Login the user
					this.$store.commit("user/setState", response.data);

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
