import * as types from "~/ui/assets/types";

export const state = () => ({
	status: null as types.ServerStatus | null,
});

export const mutations = {
	setStatus(state: any, target: types.ServerStatus | null) {
		state.status = target;
	},
};
