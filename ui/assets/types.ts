export interface Response {
	state: boolean;
	message: string;
	data: any;
	extraData: any;
}

export interface CreateUserRequest {
	firstName: string;
	lastName: string;
	email: string;
	password: string;
}

export interface UpdateUserRequest {
	uuid: string;
	firstName: string;
	lastName: string;
	email: string;
	themeColorLight: string;
	themeColorDark: string;
}

export interface ServerStatus {
	setupRequired: boolean;
}

export interface UserContext {
	user: User;
	userSettings: UserSettings;
}

export interface User {
	uuid: string;
	email: string;
	firstName: string;
	lastName: string;
}

export interface UserSettings {
	themeColorLight: string;
	themeColorDark: string;
}

export interface CreateUserRequest {
	firstName: string;
	lastName: string;
	email: string;
	password: string;
}

export interface DeleteUserRequest {
	uuid: string;
}

export interface UserLoginRequest {
	email: string;
	password: string;
}

export interface Planner {
	userUUID: string;
	date: string;
	updated: string;
	created: string;
	priorities: string[]|null;
	accomplishments: string[]|null;
	tasksToday: PlannerTask[]|null;
	tasksTomorrow: PlannerTask[]|null;
	schedule: PlannerEvent[]|null;
}

export interface PlannerTask {
	value: string;
	completed: boolean;
}

export interface PlannerEvent {
	value: string;
	end: string;
	start: string;
	color: string;
}
