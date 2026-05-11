<script setup lang="ts">
import { useUser } from "@clerk/vue";
import { computed, ref } from "vue";
import { useRouter } from "vue-router";
import { useNarutoDecks } from "../composables/naruto/useNarutoDecks";
import { usePokemonDecks } from "../composables/pokemon/usePokemonDecks";
import { useApi } from "../composables/useApi";
import { useTcgState } from "../composables/useTcgState";
import type {
  NarutoCard,
  NarutoDecks,
  PokemonCard,
  PokemonDecks,
} from "../interfaces/interfaces";

type AnyDeck = PokemonDecks | NarutoDecks;
type AnyCard = PokemonCard | NarutoCard;

const router = useRouter();
const { user } = useUser();
const { selectedTcg } = useTcgState();
const { getBlob } = useApi();

const pokemonDecks = usePokemonDecks(() => user.value?.id);
const narutoDecks = useNarutoDecks(() => user.value?.id);

const decks = computed<AnyDeck[]>(() => {
  if (selectedTcg.value === "naruto") return narutoDecks.decks.value;
  return pokemonDecks.decks.value;
});

function getPreviewImage(deck: AnyDeck): string {
  if (selectedTcg.value === "naruto")
    return narutoDecks.getPreviewImage(deck as NarutoDecks);
  return pokemonDecks.getPreviewImage(deck as PokemonDecks);
}

function getCardImage(card: AnyCard): string {
  if (selectedTcg.value === "naruto")
    return narutoDecks.getCardImage(card as NarutoCard);
  return pokemonDecks.getCardImage(card as PokemonCard);
}

function buildDeckText(deck: AnyDeck): string {
  if (selectedTcg.value === "naruto")
    return narutoDecks.buildDeckText(deck as NarutoDecks);
  return pokemonDecks.buildDeckText(deck as PokemonDecks);
}

function deleteDeckLocal(id: number) {
  if (selectedTcg.value === "naruto") return narutoDecks.deleteDeckLocal(id);
  return pokemonDecks.deleteDeckLocal(id);
}

function formatDate(date: Date | string): string {
  return new Date(date).toLocaleDateString("de-DE", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
  });
}

const dialogOpen = ref(false);
const selectedDeck = ref<AnyDeck | null>(null);

function openDeck(deck: AnyDeck) {
  selectedDeck.value = deck;
  dialogOpen.value = true;
}

function editDeck() {
  if (!selectedDeck.value) return;
  sessionStorage.setItem(
    "builder-deck",
    JSON.stringify(selectedDeck.value.entries.map((e) => e.card)),
  );
  sessionStorage.setItem("builder-deck-name", selectedDeck.value.name);
  sessionStorage.setItem("builder-deck-id", String(selectedDeck.value.ID));
  dialogOpen.value = false;
  router.push({ name: "Builder", params: { deckId: selectedDeck.value.ID } });
}

const shareMenuOpen = ref(false);

async function copyAsText() {
  if (!selectedDeck.value) return;
  const text = buildDeckText(selectedDeck.value);
  await navigator.clipboard.writeText(text);
  shareMenuOpen.value = false;
}

async function copyAsImage() {
  if (!selectedDeck.value) return;

  const tcg = selectedTcg.value === "naruto" ? "naruto" : "pokemon";
  const newTab = window.open("", "_blank");

  const blob = await getBlob(
    `/auth/${tcg}/decks/${selectedDeck.value.ID}/image`,
  );

  if (newTab) {
    const url = URL.createObjectURL(blob);
    newTab.location.href = url;
  }
  shareMenuOpen.value = false;
}

function duplicateDeck() {
  if (!selectedDeck.value) return;
  // TODO: call API to duplicate
  dialogOpen.value = false;
}

function deleteDeck() {
  if (!selectedDeck.value) return;
  deleteDeckLocal(selectedDeck.value.ID);
  dialogOpen.value = false;
}
</script>

<template>
  <div class="decks-page">
    <h2 class="page-title">Your Decks</h2>

    <div class="deck-grid">
      <div
        v-for="deck in decks"
        :key="deck.ID"
        class="deck-card"
        @click="openDeck(deck)"
      >
        <div class="deck-thumbnail">
          <img
            v-if="getPreviewImage(deck)"
            :src="getPreviewImage(deck)"
            :alt="deck.name"
            class="deck-preview-img"
          />
        </div>
        <div class="deck-info">
          <span class="deck-name">{{ deck.name }}</span>
          <span class="deck-meta">
            {{ deck.entries?.length ?? 0 }} cards ·
            {{ formatDate(deck.CreatedAt) }}
          </span>
        </div>
      </div>
    </div>

    <!-- Deck actions dialog -->
    <v-dialog v-model="dialogOpen" max-width="720">
      <v-card rounded="xl">
        <v-card-title class="dialog-title">
          {{ selectedDeck?.name }}
        </v-card-title>
        <v-card-text class="dialog-meta">
          {{ selectedDeck?.entries?.length ?? 0 }} cards · Created
          {{ selectedDeck ? formatDate(selectedDeck.CreatedAt) : "" }}
        </v-card-text>
        <div class="dialog-card-grid">
          <div
            v-for="entry in selectedDeck?.entries ?? []"
            :key="entry.ID"
            class="dialog-card-slot"
          >
            <img
              v-if="getCardImage(entry.card)"
              :src="getCardImage(entry.card)"
              :alt="(entry.card as any).Name ?? (entry.card as any).name ?? ''"
              class="dialog-card-img"
            />
          </div>
        </div>
        <v-divider />
        <div class="dialog-actions">
          <v-btn icon variant="tonal" @click="editDeck">
            <v-icon icon="mdi-pencil" />
            <v-tooltip activator="parent" location="bottom">Edit</v-tooltip>
          </v-btn>
          <v-menu v-model="shareMenuOpen" location="bottom">
            <template #activator="{ props }">
              <v-btn icon variant="tonal" v-bind="props">
                <v-icon>mdi-share-variant</v-icon>
                <v-tooltip activator="parent" location="bottom"
                  >Share</v-tooltip
                >
              </v-btn>
            </template>
            <v-card min-width="200" rounded="lg">
              <v-list density="compact">
                <v-list-item
                  prepend-icon="mdi-clipboard-text"
                  @click="copyAsText"
                >
                  <v-list-item-title>Copy to Clipboard</v-list-item-title>
                </v-list-item>
                <v-list-item prepend-icon="mdi-image" @click="copyAsImage">
                  <v-list-item-title>Copy as Image</v-list-item-title>
                </v-list-item>
              </v-list>
            </v-card>
          </v-menu>
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
  overflow: hidden;
}

.deck-preview-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: center 15%;
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
  overflow: hidden;
}

.dialog-preview-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: center 20%;
}

.dialog-card-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
  gap: 8px;
  padding: 12px;
  margin: 8px 16px;
  max-height: 400px;
  overflow-y: auto;
  background: rgba(0, 0, 0, 0.25);
  border: 1px solid var(--border);
  border-radius: 12px;
}

.dialog-card-slot {
  aspect-ratio: 5 / 7;
}

.dialog-card-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 8px;
}

.dialog-actions {
  display: flex;
  gap: 8px;
  padding: 12px 16px;
  justify-content: center;
}
</style>
