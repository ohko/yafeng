<template>
	<HelloWorld msg="Hello world!"></HelloWorld>
	<a-card>
		Token: <RouterLink to="/login">{{authStore.token}}</RouterLink>
		<hr>
		Account: {{userInfo.Account}}
	</a-card>
</template>

<script setup>
	import { ref, onMounted } from 'vue';
	import { useAuthStore } from '@/store/auth'

	const authStore = useAuthStore()
	const userInfo = ref({})

	onMounted(async _ => {
		const info = await post("/api/user/info", {})
		if (info) {
			userInfo.value = info
		}
	})
</script>