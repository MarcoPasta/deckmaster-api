<script setup lang="ts">
import { useSession, useUser } from "@clerk/vue";
import { ref } from "vue";

const { user } = useUser();
const userData = ref<any>(null);

const { session } = useSession();
const token = await session.value?.getToken();
console.log("Fetching user data with token:", token);
const res = await fetch(
  `http://0.0.0.0:8080/api/auth/user?id=${user.value?.id}`,
  {
    method: "GET",
    headers: {
      Authorization: `Bearer ${token}`,
    },
  },
);

// console.log(import.meta.env.VITE_CLERK_SECRET_KEY);

setTimeout(async () => {
  console.log("Simulating slow network...");
  userData.value = await res.json();
}, 2000);
</script>

<template>
  <div>{{ userData }}</div>
</template>
