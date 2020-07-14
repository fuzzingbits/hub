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

class HubAPI {
	public async getMe(): Promise<GenericResponse<types.UserSession>> {
		const response = await client.get("/api/test");
		return response.data;
	}
}

export default new HubAPI();
