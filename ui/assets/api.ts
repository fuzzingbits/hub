import axios, { AxiosResponse } from "axios";
import time from "./time";
import * as types from "./types";

interface GenericResponse<T> extends types.Response {
	data: T;
}

const client = axios.create();

client.interceptors.response.use(
	(response: AxiosResponse<any>): AxiosResponse<any> => {
		const contentType = response.headers["content-type"] as string;
		if (!contentType.includes("application/json")) {
			throw "not a json response";
		}

		return response;
	}
);

client.interceptors.request.use(function(config) {
	config.timeout = time.SECOND * 6;
	return config;
});

class HubAPI {
	// ---- Auto Generated Functions BEGIN ---- //

	public async serverStatus(): Promise<GenericResponse<types.ServerStatus | null>> {
		const response = await client.get("/api/server/status");
		return response.data;
	}

	public async serverSetup(payload: types.UserCreateRequest): Promise<GenericResponse<types.UserContext | null>> {
		const response = await client.post("/api/server/setup", payload);
		return response.data;
	}

	public async userNew(payload: types.UserCreateRequest): Promise<GenericResponse<types.UserContext | null>> {
		const response = await client.post("/api/user/new", payload);
		return response.data;
	}

	public async userLogin(payload: types.UserLoginRequest): Promise<GenericResponse<types.UserContext | null>> {
		const response = await client.post("/api/user/login", payload);
		return response.data;
	}

	public async userLogout(): Promise<GenericResponse<null>> {
		const response = await client.get("/api/user/logout");
		return response.data;
	}

	public async userMe(): Promise<GenericResponse<types.UserContext | null>> {
		const response = await client.get("/api/user/me");
		return response.data;
	}

	public async userList(): Promise<GenericResponse<types.User[] | null>> {
		const response = await client.get("/api/user/list");
		return response.data;
	}

	public async userDelete(): Promise<GenericResponse<null>> {
		const response = await client.get("/api/user/delete");
		return response.data;
	}

	public async userUpdate(payload: types.UserUpdateRequest): Promise<GenericResponse<types.UserContext | null>> {
		const response = await client.post("/api/user/update", payload);
		return response.data;
	}

	// ---- Auto Generated Functions END ---- //
}

export default new HubAPI();
