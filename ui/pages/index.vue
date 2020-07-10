<template>
	<div>
		<h1>{{ hello() }}</h1>
	</div>
</template>

<script lang="ts">
import Vue from "vue";
import greeter from "~/ui/assets/greeter";
export default Vue.extend({
	methods: {
		hello(): string {
			return greeter.hello("world");
		},
	},
});
</script>
