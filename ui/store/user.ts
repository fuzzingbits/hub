import * as types from "~/ui/assets/types";

export const state = () => ({
	session: null as types.UserContext | null,
});

export const mutations = {
	setState(state: any, target: types.UserContext | null) {
		state.session = target;

		if (target && target.context.userSettings.themeColor) {
			document.documentElement.style.setProperty("--primary", target.context.userSettings.themeColor);
		}
	},
};
