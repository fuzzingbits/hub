<template>
	<div class="card">
		<form v-if="session" @submit.prevent="submit" id="page-form">
			<h2>Update Profile</h2>
			<p>This is for users to update their profile.</p>
			<input type="hidden" name="uuid" :value="session.user.uuid" />
			<label>First Name <input :value="session.user.firstName" name="firstName" type="text" required/></label>
			<label>Last Name <input :value="session.user.lastName" name="lastName" type="text" required/></label>
			<label>Email <input :value="session.user.email" name="email" type="email" required/></label>
			<label>Light Theme Color <input :value="session.userSettings.themeColorLight" name="themeColorLight" type="color" @input="changeColor" required/></label>
			<label>Dark Theme Color <input :value="session.userSettings.themeColorDark" name="themeColorDark" type="color" @input="changeColor" required/></label>
			<PosterMessage :poster="formPoster" />
			<label><input type="submit" value="Submit Update"/></label>
		</form>
	</div>
</template>

<script lang="ts">
import Vue from 'vue';
import HubApi from '~/ui/assets/api';
import Poster from '~/ui/assets/poster';
import * as types from '~/ui/assets/types';

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
			const form = document.querySelector('#page-form') as HTMLFormElement;
			const formData = new FormData(form);
			document.documentElement.style.setProperty('--primary-light', formData.get('themeColorLight') as string);
			document.documentElement.style.setProperty('--primary-dark', formData.get('themeColorDark') as string);
		},
		submit(): void {
			this.formPoster.reset(true);

			const form = document.querySelector('#page-form') as HTMLFormElement;
			const formData = new FormData(form);

			HubApi.userUpdate({
				uuid: formData.get('uuid') as string,
				firstName: formData.get('firstName') as string,
				lastName: formData.get('lastName') as string,
				email: formData.get('email') as string,
				themeColorLight: formData.get('themeColorLight') as string,
				themeColorDark: formData.get('themeColorDark') as string,
			})
				.then(response => {
					this.formPoster.setResponse(response);
					if (!response.state) {
						return;
					}

					// Login the user
					this.$store.commit('user/setState', response.data);
				})
				.catch(err => {
					this.formPoster.handlerError(err);
				});
		},
	},
});
</script>
