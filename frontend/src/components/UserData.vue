<script setup lang="ts">
import { useSession } from "@clerk/vue";
import { ref } from "vue";
import { useApi } from "../composables/useApi";
import type { DigimonCards } from "../interfaces/interfaces";

const userData = ref<any>(null);
const { get } = useApi();

const { session } = useSession();
const token = await session.value?.getToken();
console.log("Fetching user data with token:", token);

const cards = await get<DigimonCards[]>("/api/digimon/cards");
</script>

<template>
  <div>{{ userData }}</div>
  <div>{{ cards }}</div>
</template>
