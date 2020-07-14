export interface Response {
	state: boolean;
	message: string;
	data: any;
	extra_data: any;
}

export interface UserSession {
	user: User;
	userSettings: UserSettings;
}

export interface User {
	uuid: string;
	firstName: string;
	lastName: string;
}

export interface UserSettings {
	themeColor: string;
}
