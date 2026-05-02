<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useApi } from "../composables/useApi";
import type { PokemonCard } from "../interfaces/interfaces";

const { get } = useApi();
const searchQuery = ref("");
const activeFilter = ref("all");

const allCards = ref<PokemonCard[]>([]);
const chosenCards = ref<PokemonCard[]>([]);

onMounted(async () => {
  try {
    const data = await get<{ content: PokemonCard[] }>("/pokemon/cards");
    allCards.value = data.content;
  } catch (e) {
    console.error("Failed to fetch cards:", e);
  }
});

const filters = ["All", "Pokémon", "Trainer", "Energy"];

function getCardImage(card: PokemonCard): string {
  if (!card.CardNumber || !card.Code) return "";
  const number = card.CardNumber.split("/")[0];
  return `/assets/cards/de/${card.Code}/${card.Code}_${number}.png`;
}

const filteredCards = computed(() => {
  return allCards.value.filter((card) => {
    const matchesSearch = (card.Name ?? "")
      .toLowerCase()
      .includes(searchQuery.value.toLowerCase());
    const matchesFilter =
      activeFilter.value === "all" ||
      (card.CardType ?? "").toLowerCase() === activeFilter.value.toLowerCase();
    return matchesSearch && matchesFilter;
  });
});

function addCard(card: PokemonCard) {
  chosenCards.value.push(card);
}

function removeCard(cardId: number) {
  const idx = chosenCards.value.findIndex((c) => c.ID === cardId);
  if (idx !== -1) chosenCards.value.splice(idx, 1);
}

interface CardGroup {
  card: PokemonCard;
  count: number;
}

const groupedChosenCards = computed<CardGroup[]>(() => {
  const map = new Map<number, CardGroup>();
  for (const card of chosenCards.value) {
    const existing = map.get(card.ID);
    if (existing) {
      existing.count++;
    } else {
      map.set(card.ID, { card, count: 1 });
    }
  }
  return [...map.values()];
});
</script>

<template>
  <div class="builder">
    <!-- Left panel: card pool -->
    <div class="builder-panel">
      <div class="toolbar">
        <div class="search-wrapper">
          <v-text-field
            v-model="searchQuery"
            variant="outlined"
            density="compact"
            placeholder="Search cards..."
            prepend-inner-icon="mdi-magnify"
            hide-details
            class="search-input"
          />
        </div>
        <div class="filters">
          <v-chip
            v-for="filter in filters"
            :key="filter"
            :variant="
              activeFilter === filter.toLowerCase() ? 'elevated' : 'outlined'
            "
            :color="
              activeFilter === filter.toLowerCase()
                ? 'deep-purple-accent-2'
                : undefined
            "
            size="small"
            @click="activeFilter = filter.toLowerCase()"
          >
            {{ filter }}
          </v-chip>
        </div>
      </div>

      <div class="card-grid">
        <div
          v-for="card in filteredCards"
          :key="card.ID"
          class="card-slot"
          @click="addCard(card)"
        >
          <img :src="getCardImage(card)" :alt="card.Name" class="card-image" />
        </div>
      </div>
    </div>

    <!-- Right panel: chosen cards -->
    <div class="builder-panel">
      <div class="panel-header">
        <span class="panel-title">Deck ({{ chosenCards.length }})</span>
      </div>

      <div class="deck-grid">
        <div
          v-for="group in groupedChosenCards"
          :key="group.card.ID"
          class="card-stack"
          @click="removeCard(group.card.ID)"
        >
          <img
            v-for="n in group.count"
            :key="n"
            :src="getCardImage(group.card)"
            :alt="group.card.Name"
            class="card-image card-chosen"
            :style="{
              position: 'absolute',
              top: `${(n - 1) * 10}px`,
              left: `${(n - 1) * 10}px`,
              zIndex: n,
            }"
          />
          <span v-if="group.count > 1" class="card-count">{{
            group.count
          }}</span>
        </div>
      </div>

      <div v-if="chosenCards.length === 0" class="empty-state">
        Click cards on the left to add them
      </div>
    </div>
  </div>
</template>

<style scoped>
.builder {
  display: flex;
  gap: 12px;
  padding: 12px 16px;
  height: calc(100vh - 64px);
  box-sizing: border-box;
}

.builder-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 12px;
  background: var(--code-bg);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 14px 16px;
  box-shadow: var(--shadow);
  overflow: hidden;
}

.toolbar {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.search-wrapper {
  max-width: 320px;
}

.search-input {
  font-family: var(--sans);
}

.filters {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.panel-header {
  display: flex;
  align-items: center;
}

.panel-title {
  font-family: var(--heading);
  font-size: 16px;
  font-weight: 600;
  color: var(--text-h);
}

.card-grid {
  flex: 1;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
  grid-auto-rows: min-content;
  gap: 10px;
  overflow-y: auto;
  padding-right: 4px;
}

.card-slot {
  aspect-ratio: 5 / 7;
  cursor: pointer;
}

.card-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 10px;
  transition:
    transform 0.2s,
    box-shadow 0.2s;
}

.card-image:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.card-stack .card-image {
  position: absolute;
  width: 100%;
  height: 100%;
}

.deck-grid {
  flex: 1;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
  grid-auto-rows: min-content;
  gap: 40px 40px;
  overflow-y: auto;
  padding-right: 4px;
}

.card-stack {
  position: relative;
  aspect-ratio: 5 / 7;
  cursor: pointer;
}

.card-count {
  position: absolute;
  bottom: 4px;
  right: 4px;
  background: var(--accent-border, #7c4dff);
  color: #fff;
  font-size: 12px;
  font-weight: 700;
  min-width: 22px;
  height: 22px;
  border-radius: 11px;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  padding: 0 6px;
}

.empty-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  color: var(--text);
  opacity: 0.6;
}
</style>
