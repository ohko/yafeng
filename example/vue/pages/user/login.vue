<template>
	<a-card style="width: 500px;margin: 10% auto;">
		<a-form :model="form" name="basic" :label-col="{ span: 8 }" :wrapper-col="{ span: 16 }" autocomplete="off" @finish="onFinish" @finishFailed="onFinishFailed">
			<a-form-item label="账号" name="Account" :rules="[{ required: true, message: '请输入账号' }]">
				<a-input v-model:value="form.Account" />
			</a-form-item>

			<a-form-item label="密码" name="Password" :rules="[{ required: true, message: '请输入密码' }]">
				<a-input-password v-model:value="form.Password" />
			</a-form-item>

			<a-form-item name="remember" :wrapper-col="{ offset: 8, span: 16 }">
				<a-checkbox v-model:checked="form.Remember">记住密码</a-checkbox>
			</a-form-item>

			<a-form-item :wrapper-col="{ offset: 8, span: 16 }">
				<a-button type="primary" html-type="submit">登陆</a-button>
			</a-form-item>
		</a-form>
	</a-card>
</template>

<script setup>
	import { ref, onMounted } from 'vue';
	import { useRoute, useRouter } from 'vue-router';
	import { useAuthStore } from '@/store/auth'

	const authStore = useAuthStore()
	const router = useRouter()
	const form = ref({ Account: authStore.account, Password: authStore.password, Remember: true })

	const onFinish = async (values) => {
		const info = await post("/api/user/login", { data: form.value })
		if (info) {
			if (form.value.Remember) authStore.save(form.value.Account, form.value.Password, info.Token)
			message.success("success")
			router.push({ path: '/' })
		}
	};

	const onFinishFailed = (errorInfo) => {
		message.error("found error");
	};
</script>