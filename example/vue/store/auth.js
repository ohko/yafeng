import { defineStore } from 'pinia'

export const useAuthStore = defineStore("auth", {
	state: () => ({ xx: '' }),
	getters: {
		xxx: (state) => state.xx,
		account: (state) => localStorage.getItem("account") || "",
		password: (state) => localStorage.getItem("password") || "",
		token: (state) => localStorage.getItem("token") || "",
	},
	actions: {
		save(account, password, token) {
			localStorage.setItem("account", account)
			localStorage.setItem("password", password)
			localStorage.setItem("token", token)
		},
		logout() {
			localStorage.setItem("account", "")
			localStorage.setItem("password", "")
			localStorage.setItem("token", "")
		}
	}
})