export interface Response {
	state: boolean;
	message: string;
	data: any;
	extra_data: any;
}

export interface UserSession {
	user: User;
}

export interface User {
	uuid: string;
	firstName: string;
	lastName: string;
}
