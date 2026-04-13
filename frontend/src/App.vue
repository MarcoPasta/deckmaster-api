<script setup lang="ts">
import { Show, SignInButton, SignUpButton, UserButton } from "@clerk/vue";
import UserData from "./components/UserData.vue";
</script>

<template>
  <header>
    <Show when="signed-out">
      <SignInButton />
      <SignUpButton />
    </Show>
    <Show when="signed-in">
      <UserButton />
    </Show>
    <Show when="signed-in">
      <Suspense fallback="Loading user data...">
        <UserData />
      </Suspense>
    </Show>
  </header>
</template>
