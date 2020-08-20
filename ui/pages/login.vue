<template>
	<div class="card">
		<form @submit.prevent="submit" id="setup-form">
			<h2>User Login</h2>
			<p>This is for users to login.</p>
			<label>Username <input name="username"/></label>
			<label>Password <input name="password" type="password"/></label>
			<label><input type="submit"/></label>
		</form>
	</div>
</template>

<script lang="ts">
import Vue from "vue";
import HubApi from "~/ui/assets/api";

export default Vue.extend({
	methods: {
		submit(): void {
			const form = document.querySelector("#setup-form") as HTMLFormElement;
			const formData = new FormData(form);

			HubApi.userLogin({
				username: formData.get("username") as string,
				password: formData.get("password") as string,
			}).then(response => {
				this.$store.commit("user/setState", response.data);
				if (response.data) {
					this.$router.push("/");
				}
			});
		},
	},
});
</script>
