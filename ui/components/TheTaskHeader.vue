<template>
	<div id="task-header">
		<nuxt-link to="/tasks">List</nuxt-link>
		<nuxt-link to="/tasks/new">New</nuxt-link>
	</div>
</template>

<script lang="ts">
import Vue from "vue";
import * as types from "~/ui/assets/types";
import md5 from "md5";
import HubApi from "~/ui/assets/api";

export default Vue.extend({});
</script>

<style></style>
