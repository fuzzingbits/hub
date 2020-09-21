<template>
	<div>
		<TheTaskHeader />
		<div class="card">
			<form @submit.prevent="submit" id="page-form">
				<h2>New Task</h2>
				<p>This is for creating a new task.</p>
				<label>Name <input name="name" type="text" required/></label>
				<label>Has Due Date? <input v-model="hasDueDate" type="checkbox"/></label>
				<label v-if="hasDueDate">Due Date <input name="dueDate" type="date" required/></label>
				<label v-if="hasDueDate"
					>Due Date Type
					<label><input name="canBeCompletedEarly" value="true" type="radio" required /> Due Before</label>
					<label><input name="canBeCompletedEarly" value="false" type="radio" checked required /> Due On</label>
				</label>
				<label>Notes <textarea name="notes"></textarea></label>

				<PosterMessage :poster="formPoster" />
				<label><input type="submit" value="Create New Task"/></label>
			</form>
		</div>
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
			hasDueDate: false,
			formPoster: new Poster(),
		};
	},
	methods: {
		submit(): void {
			this.formPoster.reset(true);

			const form = document.querySelector("#page-form") as HTMLFormElement;
			const formData = new FormData(form);
			const payload: types.TaskCreateRequest = {
				name: (formData.get("name") as string) || "",
				notes: (formData.get("note") as string) || "",
				dueDate: (formData.get("dueDate") as string) || "",
				canBeCompletedEarly: formData.get("canBeCompletedEarly") === "true",
			};
			console.log(payload);

			// HubApi.userLogin({
			// 	email: formData.get("email") as string,
			// 	password: formData.get("password") as string,
			// })
			// 	.then(response => {
			// 		this.formPoster.setResponse(response);
			// 		if (!response.state) {
			// 			return;
			// 		}

			// 		this.$store.commit("user/setState", response.data);
			// 		this.$router.push("/tasks");
			// 	})
			// 	.catch(err => {
			// 		this.formPoster.handlerError(err);
			// 	});
		},
	},
});
</script>
