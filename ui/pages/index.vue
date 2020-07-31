<template>
	<div>
		<h1>Hello, {{ fullName }}!</h1>
	</div>
</template>

<script lang="ts">
import Vue from "vue";
import HubApi from "~/ui/assets/api";
import * as types from "~/ui/assets/types";
export default Vue.extend({
	computed: {
		fullName: function(): string {
			if (this.session === null) {
				return "world";
			}

			return `${this.session.user.firstName} ${this.session.user.lastName}`;
		},
		session: function(): types.UserContext | null {
			return this.$store.state.user.session;
		},
	},
	mounted() {
		HubApi.serverStatus().then(response => {
			if (response.data && response.data.setupRequired) {
				this.$router.push("/setup");
			}
		});
	},
});
</script>
