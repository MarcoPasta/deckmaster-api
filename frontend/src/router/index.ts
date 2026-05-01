import { createRouter, createWebHistory } from "vue-router";
import BuilderView from "../views/BuilderView.vue";
import DecksView from "../views/DecksView.vue";
import HomeView from "../views/HomeView.vue";
import ProfileView from "../views/ProfileView.vue";

const routes = [
  { path: "/", name: "Home", component: HomeView },
  { path: "/builder", name: "Builder", component: BuilderView },
  { path: "/decks", name: "Decks", component: DecksView },
  { path: "/profile", name: "Profile", component: ProfileView },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
