import Axios from "axios";
import { message } from 'ant-design-vue';
import { createApp, h } from 'vue';
import { Spin } from 'ant-design-vue';

let spinInstance;

function createSpin() {
	const div = document.createElement('div');
	document.body.appendChild(div);

	const app = createApp({
		render() {
			return h(Spin, {
				spinning: true,
				// size: 'large',
				tip: ' 加载中...',
				style: {
					position: 'fixed',
					top: '0',
					left: '0',
					width: '100%',
					height: '100%',
					display: 'flex',
					alignItems: 'center',
					justifyContent: 'center',
					// backgroundColor: 'rgba(255, 255, 255, 0.1)',
					zIndex: 1000,
				},
			});
		},
	});

	spinInstance = app.mount(div);
}

function showSpin() {
	if (!spinInstance) {
		createSpin();
	}
	spinInstance.$el.style.display = 'flex';
}

function hideSpin() {
	if (spinInstance) {
		spinInstance.$el.style.display = 'none';
	}
}

const defaultConfig = {
	withCredentials: true,
	timeout: 10000,
	headers: {
		Accept: "application/json, text/plain, */*",
		"Content-Type": "application/json",
		// "X-Requested-With": "XMLHttpRequest"
	}
};

const axiosInstance = Axios.create(defaultConfig)
let token = localStorage.getItem("token") || ''

axiosInstance.interceptors.request.use(function (config) {
	showSpin()
	if (token) config.headers["token"] = token
	return config;
}, function (error) {
	hideSpin();
	return Promise.reject(error);
});

axiosInstance.interceptors.response.use(function (response) {
	hideSpin();
	try {
		// 识别平台登陆成功，获取token
		if ('data' in response && 'data' in response.data && 'Token' in response.data.data && response.data.data.Token != "") {
			token = response.data.data.Token
			localStorage.setItem("token", token)
		}
	} catch (e) { }
	return response;
}, function (error) {
	hideSpin();
	return Promise.reject(error);
});

const request = (method, url, param) => {
	return new Promise((resolve, reject) => {
		axiosInstance
			.request({ method, url, ...param })
			.then((x) => {
				// {no:[0|1],data:any}
				if ('no' in x.data) {
					if (x.data.no != 0) {
						message.error(x.data.data)
						return reject(x.data.data);
					}
					return resolve(x.data.data);
				}
				resolve(x.data);
			})
			.catch(error => {
				message.error(error.toString())
				reject(error);
			});
	})
}

const post = (url, params) => request('post', url, params)
const get = (url, params) => request('get', url, params)

export { post, get };