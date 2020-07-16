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
	// TODO: create actual session management and stop doing this terrible thing
	config.headers.UUID = "313efbe9-173b-4a1b-9a5b-7b69d95a66b9";

	return config;
});

class HubAPI {
	public async getMe(): Promise<GenericResponse<types.UserSession>> {
		const response = await client.get("/api/users/me");
		return response.data;
	}
}

export default new HubAPI();
