<template>
	<div id="header">
		<div id="hamburger" @click="hideMobileMenu = !hideMobileMenu">
			<i class="fas fa-bars"></i>
		</div>
		<div id="activity">
			<nuxt-link :to="profilePictureHref">
				<img :src="profilePictureURL" alt="Profile Picture" />
			</nuxt-link>
		</div>
		<div id="brand">
			<nuxt-link to="/">
				<svg width="512px" height="512px" viewBox="0 0 512 512" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
					<title>Hub Logo</title>
					<path
						d="M300.241565,11.7367846 L446.076299,95.934508 C473.181243,111.583554 489.878576,140.504184 489.878576,171.802277 L489.878576,340.197723 C489.878576,371.495816 473.181243,400.416446 446.076299,416.065492 L300.241565,500.263215 C273.136622,515.912262 239.741954,515.912262 212.637011,500.263215 L66.8022766,416.065492 C39.6973335,400.416446 23,371.495816 23,340.197723 L23,171.802277 C23,140.504184 39.6973335,111.583554 66.8022766,95.934508 L212.637011,11.7367846 C239.741954,-3.91226155 273.136622,-3.91226155 300.241565,11.7367846 Z"
						id="brand-icon"
					></path>
				</svg>
			</nuxt-link>
		</div>
		<div id="menu" :class="{ hide: hideMobileMenu }">
			<nuxt-link to="/">Home Page</nuxt-link>
			<nuxt-link to="/about">About</nuxt-link>
			<nuxt-link v-if="setupRequired" to="/setup">Setup</nuxt-link>
			<nuxt-link v-if="!session && !setupRequired" to="/login">Login</nuxt-link>
			<nuxt-link v-if="session && !setupRequired" to="/profile">Profile</nuxt-link>
			<nuxt-link v-if="session && !setupRequired" to="/tasks">Tasks</nuxt-link>
			<a v-if="session && !setupRequired" @click.prevent="logout" href="#">Logout</a>
		</div>
	</div>
</template>

<script lang="ts">
import Vue from "vue";
import * as types from "~/ui/assets/types";
import md5 from "md5";
import HubApi from "~/ui/assets/api";

export default Vue.extend({
	data: function() {
		return {
			hideMobileMenu: true,
		};
	},
	methods: {
		logout() {
			HubApi.userLogout();
			this.$store.commit("user/setState", null);
			this.$router.push("/");
		},
	},
	computed: {
		profilePictureHref: function(): string {
			if (this.setupRequired) {
				return "/setup";
			}

			if (!this.session) {
				return "/login";
			}

			return "/profile";
		},
		profilePictureURL: function(): string {
			let email = "";
			if (this.session) {
				email = this.session.user.email;
				email.toLowerCase();
				email.trim();
			}

			const hash = md5(email);

			return `https://www.gravatar.com/avatar/${hash}.jpg?d=mp`;
		},
		session: function(): types.UserContext | null {
			return this.$store.state.user.session;
		},
		serverStatus: function(): types.ServerStatus | null {
			return this.$store.state.server.status;
		},
		setupRequired: function(): boolean {
			const serverStatus = this.serverStatus;
			if (!serverStatus) {
				return false;
			}

			return serverStatus.setupRequired;
		},
	},
	watch: {
		$route() {
			this.hideMobileMenu = true;
		},
	},
});
</script>

<style>
#brand-icon {
	fill: var(--primary);
}

#header {
	--image-size: 32px;
	--activity-width: calc(var(--image-size) + (var(--spacing) * 2));
	--brand-width: calc(var(--image-size) + (var(--spacing) * 2));
	--header-height: 48px;
	background: var(--background-accent);
	position: relative;
}

#header a {
	max-height: var(--header-height);
	display: inline-block;
}

#brand,
#activity {
	height: var(--header-height);
}

#brand svg {
	width: var(--image-size);
	height: var(--image-size);
	margin: calc((var(--header-height) - var(--image-size)) / 2) var(--spacing);
}

#activity img {
	width: var(--image-size);
	height: var(--image-size);
	border-radius: var(--image-size);
	margin: calc((var(--header-height) - var(--image-size)) / 2) var(--spacing);
}

@media (min-width: 768px) {
	#header {
		height: var(--header-height);
	}

	#brand,
	#menu,
	#activity {
		position: absolute;
		top: 0;
		bottom: 0;
	}

	#brand {
		left: 0;
		width: var(--brand-width);
	}

	#menu {
		right: var(--activity-width);
		left: var(--brand-width);
	}

	#menu a {
		line-height: var(--header-height);
		padding: 0 var(--spacing);
		float: left;
	}

	#activity {
		right: 0;
		width: var(--activity-width);
	}

	#hamburger {
		display: none;
	}
}

@media (max-width: 767px) {
	#brand {
		text-align: center;
	}

	#activity {
		float: right;
		text-align: center;
		width: var(--activity-width);
	}

	#hamburger {
		cursor: pointer;
		display: block;
		float: left;
		text-align: center;
		line-height: var(--header-height);
		width: var(--activity-width);
	}

	#menu {
		border-top: 1px solid var(--background);
	}

	#menu.hide {
		display: none;
	}

	#menu a {
		border-top: 1px solid var(--background);
		display: block;
		padding: var(--spacing) var(--spacing);
	}
}
</style>
