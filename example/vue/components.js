import { defineAsyncComponent } from 'vue';

const components = import.meta.glob('./components/**/*.vue');

export default {
	install(app) {
		for (const path in components) {
			const componentName = path
				.replace('./components/', '')
				.replace(/\//g, '-')
				.replace('.vue', '');
			app.component(componentName, defineAsyncComponent(components[path]));
		}
	},
};
