<template>
	<div>
		<TheSplash />
		<div id="page">
			<TheHeader />
			<div id="body">
				<nuxt />
			</div>
			<TheFooter />
		</div>
	</div>
</template>

<script lang="ts">
import Vue from 'vue';
import HubApi from '~/ui/assets/api';
import Poster from '~/ui/assets/poster';
import * as types from '~/ui/assets/types';

export default Vue.extend({});
</script>

<style>
@import url("../assets/prism/prism.css");

.fade-enter-active,
.fade-leave-active {
	transition: opacity 0.5s;
}

.fade-enter,
.fade-leave-to {
	opacity: 0;
}

a[href^="http"]:after {
	content: "\f35d";
	font-family: "Font Awesome 5 Free";
	font-weight: 900;
}

a.nuxt-link-exact-active {
	text-decoration: underline;
}

#body {
	min-height: 200px;
	margin: var(--spacing);
}

.alert {
	border: 1px solid green;
	color: green;
	position: relative;
	padding: var(--spacing) var(--spacing-double);
	margin-bottom: var(--spacing-half);
}

.alert:last-child {
	margin-bottom: 0;
}

.alert.alert-danger {
	color: tomato;
	border: 1px solid tomato;
}

.alert .close {
	position: absolute;
	top: var(--spacing-half);
	right: var(--spacing);
	cursor: pointer;
}

.card {
	background: var(--background-accent);
	padding: var(--spacing-double);
	max-width: 512px;
	margin: auto;
}
</style>
