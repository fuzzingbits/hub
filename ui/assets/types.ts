export interface Response {
	state: boolean;
	message: string;
	data: any;
	extra_data: any;
}

export interface CreateUserRequest {
	email: string;
	username: string;
	password: string;
}

export interface ServerStatus {
	setupRequired: boolean;
}

export interface UserSession {
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
