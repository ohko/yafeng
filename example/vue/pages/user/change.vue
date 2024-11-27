<template>
	<a-card style="width: 500px;margin: 10% auto;">
		<a-form :model="form" name="basic" :label-col="{ span: 8 }" :wrapper-col="{ span: 16 }" autocomplete="off" @finish="onFinish" @finishFailed="onFinishFailed">
			<a-form-item label="密码" name="Password" :rules="[{ required: true, message: '请输入密码' }]">
				<a-input-password v-model:value="form.Password" />
			</a-form-item>
			<a-form-item label="确认密码" name="Password2" :rules="[{ required: true, message: '请确认密码' }]">
				<a-input-password v-model:value="form.Password2" />
			</a-form-item>

			<a-form-item :wrapper-col="{ offset: 8, span: 16 }">
				<a-button type="primary" html-type="submit">提交</a-button>
			</a-form-item>
		</a-form>
	</a-card>
</template>

<script setup>
	import { ref, onMounted } from 'vue';
	import { useRoute, useRouter } from 'vue-router';
	import { useTokenStore } from '@/store/token'

	const tokenStore = useTokenStore()
	const router = useRouter()
	const form = ref({ Password: "", Password2: "" })

	const onFinish = async (values) => {
		if (form.value.Password != form.value.Password2) return message.error("两次输入的密码不一致")
		const info = await post("/api/user/change", { data: form.value })
		if (info) {
			message.success("success")
		}
	};

	const onFinishFailed = (errorInfo) => {
		message.error("found error");
	};
</script>