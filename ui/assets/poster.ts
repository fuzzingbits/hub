import * as types from './types';

function isApiResponse(response: types.Response): boolean {
	return true;
}

class Poster {
	public state = true;
	public complete = false;
	public loading = false;
	public message = '';

	public constructor(startLoading = false) {
		this.loading = startLoading;
	}

	public setResponse(response: types.Response): void {
		// Complete the request
		this.loading = false;
		this.complete = true;

		// Copy values out of the response
		this.state = response.state;
		this.message = response.message;
	}

	public handlerError(err: any): void {
		// Complete the request
		this.loading = false;
		this.complete = true;

		// Copy values out of the response
		if (err && err.response && err.response.data) {
			const response = err.response.data as types.Response;
			if (isApiResponse(response)) {
				this.state = response.state;
				this.message = response.message;
				return;
			}
		}

		this.state = false;
		this.message = 'Something went wrong';
	}

	public reset(loading = false) {
		this.complete = false;
		this.state = true;
		this.message = '';
		this.loading = loading;
	}
}

export default Poster;
