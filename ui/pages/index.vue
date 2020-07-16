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
	data: function () {
		return {
			userSession: null as types.UserSession | null,
		};
	},
	computed: {
		fullName: function (): string {
			if (this.userSession === null) {
				return "world";
			}

			return `${this.userSession.user.firstName} ${this.userSession.user.lastName}`;
		},
	},
	mounted() {
		HubApi.getMe()
			.then((response) => {
				this.userSession = response.data;
			})
			.catch((err) => {
				console.error("ajax error: " + err);
			})
			.finally(() => {
				// console.log("completed");
			});
	},
});
</script>
