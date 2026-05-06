<script setup lang="ts">
import { watch } from "vue";
import NavBar from "./components/NavBar.vue";
import { useTcgState } from "./composables/useTcgState";

const { selectedTcg } = useTcgState();

watch(
  selectedTcg,
  (tcg) => {
    document.body.setAttribute("data-tcg", tcg);
  },
  { immediate: true },
);
</script>

<template>
  <NavBar />
  <router-view />
</template>
