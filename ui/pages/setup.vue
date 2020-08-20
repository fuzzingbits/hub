<template>
	<div class="card">
		<form @submit.prevent="submit" id="setup-form">
			<h2>Server Setup</h2>
			<p>This is for first time setups only.</p>
			<label>First Name <input name="firstName"/></label>
			<label>Last Name <input name="lastName"/></label>
			<label>Email <input name="email" type="email"/></label>
			<label>Username <input name="username"/></label>
			<label>Password <input name="password" type="password"/></label>
			<label><input type="submit"/></label>
		</form>
	</div>
</template>

<script lang="ts">
import Vue from "vue";
import HubApi from "~/ui/assets/api";
import * as types from "~/ui/assets/types";

export default Vue.extend({
	computed: {
		serverStatus: function(): types.ServerStatus | null {
			return this.$store.state.server.status;
		},
	},
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
				const serverStatus = this.serverStatus;
				this.$store.commit("user/setState", response.data);
				if (response.data && serverStatus) {
					this.$router.push("/");
					serverStatus.setupRequired = false;
				}
				this.$store.commit("server/setStatus", serverStatus);
			});
		},
	},
});
</script>
