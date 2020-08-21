import * as types from "~/ui/assets/types";

export const state = () => ({
	session: null as types.UserContext | null,
});

export const mutations = {
	setState(state: any, target: types.UserContext | null) {
		state.session = target;

		let faviconURL = "/favicon.svg";

		if (state.session) {
			let color = encodeURIComponent(state.session.userSettings.themeColorLight);
			const isDarkMode = window.matchMedia("(prefers-color-scheme: dark)").matches;
			if (isDarkMode) {
				color = encodeURIComponent(state.session.userSettings.themeColorDark);
			}
			faviconURL += `?color=${color}&time=${new Date().getTime()}`;
			console.log(color);
			console.log(faviconURL);
			document.documentElement.style.setProperty("--primary-light", state.session.userSettings.themeColorLight);
			document.documentElement.style.setProperty("--primary-dark", state.session.userSettings.themeColorDark);
		} else {
			document.documentElement.style.removeProperty("--primary-light");
			document.documentElement.style.removeProperty("--primary-dark");
		}

		(document.querySelector("link[rel='icon']") as HTMLLinkElement).href = faviconURL;
	},
};
