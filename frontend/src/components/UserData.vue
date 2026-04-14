<script setup lang="ts">
import { useSession, useUser } from "@clerk/vue";
import { ref } from "vue";

const { user } = useUser();
const userData = ref<any>(null);

const { session } = useSession();
const token = await session.value?.getToken();
console.log("Fetching user data with token:", token);
// const res = await fetch(
//   `http://0.0.0.0:8080/api/auth/user?id=${user.value?.id}`,
//   {
//     method: "GET",
//     headers: {
//       Authorization: `Bearer ${token}`,
//     },
//   },
// );

const deckres = await fetch(`http://0.0.0.0:8080/api/auth/decks/create`, {
  method: "POST",
  headers: {
    Authorization: `Bearer ${token}`,
  },
  body: JSON.stringify({
    userId: "marcopasta",
    name: "First Deck",
    cards: [
      {
        Id: "BT24-001",
        Name: "Gigimon",
        CardType: "Digi-Egg",
        Color: ["Red"],
        Cost: 0,
        Dp: 0,
        DvCost: [""],
        Form: "In-Training",
        Attribute: "",
        Species: ["Lesser/LIBERATOR"],
        Text: [
          "[Inherited Effect]",
          "[Your Turn] [Once Per Turn] When your opponent's security stack is removed from, you may delete 1 of their Digimon with 3000 DP or less.",
        ],
        Set: "TIME STRANGER [BT-24]",
      },
      {
        Id: "BT24-002",
        Name: "Bukamon",
        CardType: "Digi-Egg",
        Color: ["Blue"],
        Cost: 0,
        Dp: 0,
        DvCost: [""],
        Form: "In-Training",
        Attribute: "",
        Species: ["Lesser/Iliad/TS"],
        Text: [
          "[Inherited Effect]",
          "[End of Your Turn] [Once Per Turn] By paying 1 cost, this blue Digimon with the [TS] trait unsuspends.",
        ],
        Set: "TIME STRANGER [BT-24]",
      },
    ],
  }),
});

// console.log(import.meta.env.VITE_CLERK_SECRET_KEY);

setTimeout(async () => {
  console.log("Simulating slow network...");
  userData.value = await deckres.json();
}, 2000);
</script>

<template>
  <div>{{ userData }}</div>
</template>
