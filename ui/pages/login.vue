<template>
	<form @submit.prevent="submit" id="setup-form">
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

			HubApi.userLogin({
				username: formData.get("username") as string,
				password: formData.get("password") as string,
			}).then(response => {
				this.$store.commit("user/setState", response.data);
			});
		},
	},
	mounted() {
		HubApi.serverStatus().then(response => {
			if (response.data && response.data.setupRequired) {
				this.$router.push("/setup");
			}
		});
		console.log("This is the login page.");
	},
});
</script>
