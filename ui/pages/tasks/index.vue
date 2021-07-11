<template>
	<div>
		<TheTaskHeader />

		<h1>Tasks!</h1>
		<ul>
			<li v-for="task in tasks" :key="task.uuid">
				<nuxt-link :to="`/tasks/view/${task.uuid}`">{{ task.name }}</nuxt-link>
			</li>
		</ul>
	</div>
</template>

<script lang="ts">
import Vue from "vue";
import HubApi from "~/ui/assets/api";
import * as types from "~/ui/assets/types";

export default Vue.extend({
	computed: {
		tasks: function(): types.Task[] | null {
			return this.$store.state.task.tasks;
		},
	},
});
</script>
