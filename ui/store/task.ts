import * as types from "~/ui/assets/types";

export const state = () => ({
	tasks: [
		{
			uuid: "234-432-432-432",
			name: "Laundry",
			dueDate: "2020-05-22",
			canBeCompletedEarly: true,
			completed: false,
			deletedAt: null,
			notes: "Just a simple note.",
			userUUID: "bla-bla-bla",
		},
	] as types.Task[] | null,
});

export const mutations = {
	setTasks(state: any, target: types.Task[] | null) {
		state.tasks = target;
	},
};
