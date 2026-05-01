import { clerkPlugin } from "@clerk/vue";
import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import "./style.css";
// Vuetify
import { createVuetify } from "vuetify";
import * as components from "vuetify/components";
import * as directives from "vuetify/directives";
import "vuetify/styles";

const PUBLISHABLE_KEY = import.meta.env.VITE_CLERK_PUBLISHABLE_KEY;

if (!PUBLISHABLE_KEY) {
  throw new Error("Add your Clerk Publishable Key to the .env file");
}

const vuetify = createVuetify({
  components,
  directives,
});

const app = createApp(App);
app.use(router);
app.use(clerkPlugin, { publishableKey: PUBLISHABLE_KEY });
app.use(vuetify);
app.mount("#app");
