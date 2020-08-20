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
	username: string;
	password: string;
}

export interface UpdateUserRequest {
	uuid: string;
	firstName: string;
	lastName: string;
	email: string;
	username: string;
	themeColor: string;
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
	username: string;
	email: string;
	firstName: string;
	lastName: string;
}

export interface UserSettings {
	themeColor: string;
}

export interface CreateUserRequest {
	firstName: string;
	lastName: string;
	email: string;
	username: string;
	password: string;
}

export interface DeleteUserRequest {
	uuid: string;
}

export interface UserLoginRequest {
	username: string;
	password: string;
}
