import axios, { AxiosResponse } from "axios";
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

client.interceptors.request.use(function (config) {
	return config;
});

class HubAPI {

	public async serverStatus(): Promise<GenericResponse<types.ServerStatus | null>> {
		const response = await client.get("/api/server/status");
		return response.data;
	}

	public async serverSetup(payload: types.CreateUserRequest): Promise<GenericResponse<types.UserContext | null>> {
		const response = await client.post("/api/server/setup", payload);
		return response.data;
	}

	public async userLogin(payload: types.UserLoginRequest): Promise<GenericResponse<types.UserContext | null>> {
		const response = await client.post("/api/user/login", payload);
		return response.data;
	}

	public async userMe(): Promise<GenericResponse<types.UserContext | null>> {
		const response = await client.get("/api/user/me");
		return response.data;
	}

}

export default new HubAPI();
