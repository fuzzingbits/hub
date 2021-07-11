<template>
	<div>
		<TheTaskHeader />
		<div class="card">
			<form @submit.prevent="submit" id="page-form">
				<h2>New Task</h2>
				<p>This is for creating a new task.</p>
				<label>Completed? <input name="completed" type="checkbox"/></label>
				<label>Name <input name="name" type="text" :value="task.name" required/></label>
				<label>Has Due Date? <input v-model="hasDueDate" type="checkbox"/></label>
				<label v-if="hasDueDate">Due Date <input name="dueDate" type="date" :value="task.dueDate" required/></label>
				<label v-if="hasDueDate"
					>Due Date Type
					<label><input name="canBeCompletedEarly" value="true" type="radio" :checked="task.canBeCompletedEarly" required /> Due Before</label>
					<label><input name="canBeCompletedEarly" value="false" type="radio" :checked="!task.canBeCompletedEarly" required /> Due On</label>
				</label>
				<label
					>Notes <textarea name="notes">{{ task.notes }}</textarea></label
				>

				<PosterMessage :poster="formPoster" />
				<label><input type="submit" value="Update Task"/></label>
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
	computed: {
		uuid: function(): string {
			return this.$route.params["uuid"];
		},
		task: function(): types.Task | null {
			for (let i = 0; i < this.$store.state.task.tasks.length; i++) {
				const task = this.$store.state.task.tasks[i];
				if (task.uuid === this.uuid) {
					this.hasDueDate = task.dueDate !== "";
					return task;
				}
			}

			return null;
		},
	},
});
</script>
