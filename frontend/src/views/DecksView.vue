<script setup lang="ts">
import { ref } from "vue";

interface Deck {
  id: number;
  name: string;
  cardCount: number;
  updatedAt: string;
}

// TODO: replace with API call
const decks = ref<Deck[]>([
  { id: 1, name: "Mono Red Aggro", cardCount: 60, updatedAt: "2026-04-28" },
  { id: 2, name: "Blue Control", cardCount: 60, updatedAt: "2026-04-27" },
  { id: 3, name: "Green Ramp", cardCount: 60, updatedAt: "2026-04-25" },
  { id: 4, name: "Black Midrange", cardCount: 60, updatedAt: "2026-04-24" },
  { id: 5, name: "White Weenie", cardCount: 60, updatedAt: "2026-04-22" },
  { id: 6, name: "Izzet Spells", cardCount: 60, updatedAt: "2026-04-20" },
]);

const dialogOpen = ref(false);
const selectedDeck = ref<Deck | null>(null);

function openDeck(deck: Deck) {
  selectedDeck.value = deck;
  dialogOpen.value = true;
}

function editDeck() {
  // TODO: navigate to builder with deck loaded
  console.log("Edit deck:", selectedDeck.value?.id);
  dialogOpen.value = false;
}

function shareDeck() {
  // TODO: call share API / copy link
  console.log("Share deck:", selectedDeck.value?.id);
  dialogOpen.value = false;
}

function duplicateDeck() {
  if (!selectedDeck.value) return;
  const copy: Deck = {
    ...selectedDeck.value,
    id: Date.now(),
    name: `${selectedDeck.value.name} (Copy)`,
    updatedAt: new Date().toISOString().slice(0, 10),
  };
  decks.value.push(copy);
  dialogOpen.value = false;
}

function deleteDeck() {
  if (!selectedDeck.value) return;
  decks.value = decks.value.filter((d) => d.id !== selectedDeck.value!.id);
  dialogOpen.value = false;
}
</script>

<template>
  <div class="decks-page">
    <h2 class="page-title">Your Decks</h2>

    <div class="deck-grid">
      <div
        v-for="deck in decks"
        :key="deck.id"
        class="deck-card"
        @click="openDeck(deck)"
      >
        <div class="deck-thumbnail"></div>
        <div class="deck-info">
          <span class="deck-name">{{ deck.name }}</span>
          <span class="deck-meta"
            >{{ deck.cardCount }} cards · {{ deck.updatedAt }}</span
          >
        </div>
      </div>
    </div>

    <!-- Deck actions dialog -->
    <v-dialog v-model="dialogOpen" max-width="360">
      <v-card rounded="xl">
        <v-card-title class="dialog-title">
          {{ selectedDeck?.name }}
        </v-card-title>
        <v-card-text class="dialog-meta">
          {{ selectedDeck?.cardCount }} cards · Updated
          {{ selectedDeck?.updatedAt }}
        </v-card-text>
        <div class="dialog-thumbnail"></div>
        <v-divider />
        <div class="dialog-actions">
          <v-btn icon variant="tonal" @click="editDeck">
            <v-icon>mdi-pencil</v-icon>
            <v-tooltip activator="parent" location="bottom">Edit</v-tooltip>
          </v-btn>
          <v-btn icon variant="tonal" @click="shareDeck">
            <v-icon>mdi-share-variant</v-icon>
            <v-tooltip activator="parent" location="bottom">Share</v-tooltip>
          </v-btn>
          <v-btn icon variant="tonal" @click="duplicateDeck">
            <v-icon>mdi-content-copy</v-icon>
            <v-tooltip activator="parent" location="bottom"
              >Duplicate</v-tooltip
            >
          </v-btn>
          <v-btn icon variant="tonal" color="error" @click="deleteDeck">
            <v-icon>mdi-delete</v-icon>
            <v-tooltip activator="parent" location="bottom">Delete</v-tooltip>
          </v-btn>
        </div>
        <v-divider />
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="dialogOpen = false">Cancel</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<style scoped>
.decks-page {
  padding: 16px 24px;
}

.page-title {
  font-size: 22px;
  margin: 0 0 16px;
}

.deck-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 14px;
}

.deck-card {
  background: var(--code-bg);
  border: 1px solid var(--border);
  border-radius: 14px;
  overflow: hidden;
  box-shadow: var(--shadow);
  cursor: pointer;
  transition:
    border-color 0.2s,
    transform 0.15s;
}

.deck-card:hover {
  border-color: var(--accent-border);
  transform: translateY(-2px);
}

.deck-thumbnail {
  width: 100%;
  aspect-ratio: 16 / 10;
  background: var(--accent-bg);
  border-bottom: 1px solid var(--border);
}

.deck-info {
  padding: 12px 14px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.deck-name {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-h);
}

.deck-meta {
  font-size: 12px;
  color: var(--text);
}

.dialog-title {
  font-size: 18px;
  font-weight: 600;
  padding-bottom: 0;
}

.dialog-meta {
  font-size: 13px;
  padding-top: 4px;
}

.dialog-thumbnail {
  width: calc(100% - 32px);
  aspect-ratio: 16 / 10;
  background: var(--accent-bg);
  border: 1px solid var(--border);
  border-radius: 10px;
  margin: 12px 16px;
}

.dialog-actions {
  display: flex;
  gap: 8px;
  padding: 12px 16px;
  justify-content: center;
}
</style>
