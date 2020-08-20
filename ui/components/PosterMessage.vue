<template>
	<div
		v-if="poster.complete && poster.message"
		class="alert"
		:class="{
			'alert-danger': !poster.state,
		}"
	>
		{{ poster.message }}
	</div>
</template>

<script lang="ts">
import Vue from "vue";
import Poster from "~/ui/assets/poster";

export default Vue.extend({
	props: {
		poster: {
			required: true,
			type: Poster,
		},
	},
});
</script>
