<template>
	<HelloWorld msg="Hello world!"></HelloWorld>
	<a-card>
		Token: <RouterLink to="/login">{{tokenStore.token}}</RouterLink>
		<hr>
		Account: {{userInfo.Account}}
	</a-card>
</template>

<script setup>
	import { ref, onMounted } from 'vue';
	import { useTokenStore } from '@/store/token'
	const tokenStore = useTokenStore()

	const userInfo = ref({})

	onMounted(async _ => {
		const info = await post("/api/user/info", {})
		if (info) {
			userInfo.value = info
		}
	})
</script>