import { createRouter, createWebHistory } from "vue-router";
import BuilderView from "../views/BuilderView.vue";
import DecksView from "../views/DecksView.vue";
import HomeView from "../views/HomeView.vue";
import ProfileView from "../views/ProfileView.vue";
import UploadView from "../views/UploadView.vue";

const routes = [
  { path: "/", name: "Home", component: HomeView },
  { path: "/builder/:deckId?", name: "Builder", component: BuilderView },
  {
    path: "/decks",
    name: "Decks",
    component: DecksView,
    meta: { requiresAuth: true },
  },
  { path: "/profile", name: "Profile", component: ProfileView },
  { path: "/upload", name: "Upload", component: UploadView },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
