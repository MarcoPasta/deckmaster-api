<script setup lang="ts">
import { Show, SignInButton, UserButton } from "@clerk/vue";

const navItems = [
  { label: "Builder", path: "/builder" },
  { label: "Decks", path: "/decks" },
  { label: "Profile", path: "/profile" },
  { label: "", path: "/" },
];
</script>

<template>
  <nav class="navbar">
    <router-link
      v-for="item in navItems"
      :key="item.path"
      :to="item.path"
      class="nav-item"
      active-class="nav-item--active"
    >
      {{ item.label }}
    </router-link>
    <div class="nav-item nav-user">
      <Show when="signed-in">
        <UserButton />
      </Show>
      <Show when="signed-out">
        <SignInButton />
      </Show>
    </div>
  </nav>
</template>

<style scoped>
.navbar {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px;
  background: var(--code-bg);
  border: 1px solid var(--border);
  border-radius: 12px;
  box-shadow: var(--shadow);
  width: fit-content;
  margin: 8px auto 0;
}

.nav-item {
  padding: 8px 20px;
  border-radius: 8px;
  font-family: var(--sans);
  font-size: 15px;
  font-weight: 500;
  color: var(--text);
  text-decoration: none;
  transition:
    background 0.2s,
    color 0.2s;
  min-width: 40px;
  text-align: center;
}

.nav-item:hover {
  background: var(--accent-bg);
  color: var(--text-h);
}

.nav-item--active {
  background: var(--accent);
  color: #fff;
}

.nav-item--active:hover {
  background: var(--accent);
  color: #fff;
}

.nav-user {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4px 12px;
}
</style>
