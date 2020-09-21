export interface ServerStatus {
	setupRequired: boolean;
}

export interface Task {
	uuid: string;
	userUUID: string;
	name: string;
	notes: string;
	dueDate: string;
	completed: boolean;
	createdAt: string;
	deletedAt: string|null;
	canBeCompletedEarly: boolean;
}

export interface TaskCreateRequest {
	name: string;
	notes: string;
	dueDate: string;
	canBeCompletedEarly: boolean;
}

export interface User {
	uuid: string;
	email: string;
	firstName: string;
	lastName: string;
}

export interface UserContext {
	user: User;
	userSettings: UserSettings;
}

export interface UserCreateRequest {
	firstName: string;
	lastName: string;
	email: string;
	password: string;
}

export interface UserDeleteRequest {
	uuid: string;
}

export interface UserLoginRequest {
	email: string;
	password: string;
}

export interface UserSettings {
	themeColorLight: string;
	themeColorDark: string;
}

export interface UserUpdateRequest {
	uuid: string;
	firstName: string;
	lastName: string;
	email: string;
	themeColorLight: string;
	themeColorDark: string;
}

export interface Response {
	state: boolean;
	message: string;
	data: any;
	extraData: any;
}
