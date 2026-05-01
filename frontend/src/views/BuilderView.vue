<script setup lang="ts">
import { computed, ref } from "vue";

const searchQuery = ref("");
const activeFilter = ref("all");
const filters = ["All", "Creatures", "Spells", "Artifacts", "Lands"];

interface Card {
  id: number;
  name: string;
  type: string;
}

const allCards: Card[] = Array.from({ length: 24 }, (_, i) => ({
  id: i + 1,
  name: `Card ${i + 1}`,
  type: ["creatures", "spells", "artifacts", "lands"][i % 4],
}));

const chosenCards = ref<Card[]>([]);

const filteredCards = computed(() => {
  return allCards.filter((card) => {
    const matchesSearch = card.name
      .toLowerCase()
      .includes(searchQuery.value.toLowerCase());
    const matchesFilter =
      activeFilter.value === "all" || card.type === activeFilter.value;
    return matchesSearch && matchesFilter;
  });
});

function addCard(card: Card) {
  chosenCards.value.push({ ...card, id: Date.now() + Math.random() });
}

function removeCard(index: number) {
  chosenCards.value.splice(index, 1);
}
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
          :key="card.id"
          class="card-slot"
          @click="addCard(card)"
        >
          <div class="card-placeholder">
            <span class="card-label">{{ card.name }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Right panel: chosen cards -->
    <div class="builder-panel">
      <div class="panel-header">
        <span class="panel-title">Deck ({{ chosenCards.length }})</span>
      </div>

      <div class="card-grid">
        <div
          v-for="(card, index) in chosenCards"
          :key="card.id"
          class="card-slot"
          @click="removeCard(index)"
        >
          <div class="card-placeholder card-chosen">
            <span class="card-label">{{ card.name }}</span>
          </div>
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
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  grid-auto-rows: min-content;
  gap: 10px;
  overflow-y: auto;
  padding-right: 4px;
}

.card-slot {
  aspect-ratio: 5 / 7;
  cursor: pointer;
}

.card-placeholder {
  width: 100%;
  height: 100%;
  border: 2px dashed var(--border);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition:
    border-color 0.2s,
    background 0.2s;
}

.card-placeholder:hover {
  border-color: var(--accent-border);
  background: var(--accent-bg);
}

.card-chosen {
  border-style: solid;
  border-color: var(--accent-border);
  background: var(--accent-bg);
}

.card-chosen:hover {
  border-color: #e53935;
  background: rgba(229, 57, 53, 0.1);
}

.card-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--text);
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
