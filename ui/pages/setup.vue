<template>
	<form @submit.prevent="submit" id="setup-form">
		<div>
			<label>First Name</label>
			<input name="firstName" />
		</div>
		<div>
			<label>Last Name</label>
			<input name="lastName" />
		</div>
		<div>
			<label>Email</label>
			<input name="email" />
		</div>
		<div>
			<label>Username</label>
			<input name="username" />
		</div>
		<div>
			<label>Password</label>
			<input name="password" />
		</div>
		<div>
			<input type="submit" />
		</div>
	</form>
</template>

<script lang="ts">
import Vue from "vue";
import HubApi from "~/ui/assets/api";

export default Vue.extend({
	methods: {
		submit(): void {
			const form = document.querySelector("#setup-form") as HTMLFormElement;
			const formData = new FormData(form);

			HubApi.serverSetup({
				firstName: formData.get("firstName") as string,
				lastName: formData.get("lastName") as string,
				email: formData.get("email") as string,
				username: formData.get("username") as string,
				password: formData.get("password") as string,
			}).then(response => {
				this.$store.commit("user/setState", response.data);
			});
		},
	},
	mounted() {
		HubApi.serverStatus().then(response => {
			if (!response.data) {
				return;
			}

			if (!response.data.setupRequired) {
				this.$router.push("/");
			}
		});
		console.log("This is the setup page.");
	},
});
</script>
