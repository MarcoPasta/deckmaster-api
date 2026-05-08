<script setup lang="ts">
import { useUser } from "@clerk/vue";
import { computed, onMounted, onUnmounted, ref, watch } from "vue";
import { useNarutoBuilder } from "../composables/naruto/useNarutoBuilder";
import { usePokemonBuilder } from "../composables/pokemon/usePokemonBuilder";
import { useTcgState } from "../composables/useTcgState";
import type { NarutoCard, PokemonCard } from "../interfaces/interfaces";

type AnyCard = PokemonCard | NarutoCard;

const STORAGE_KEY = "builder-deck";
const STORAGE_NAME_KEY = "builder-deck-name";
const STORAGE_ID_KEY = "builder-deck-id";

const { user } = useUser();
const { selectedTcg } = useTcgState();

const pokemon = usePokemonBuilder();
const naruto = useNarutoBuilder();

const filters = computed(() => {
  if (selectedTcg.value === "naruto") return naruto.filters;
  return pokemon.filters;
});

function getCardImage(card: AnyCard): string {
  if (selectedTcg.value === "naruto")
    return naruto.getCardImage(card as NarutoCard);
  return pokemon.getCardImage(card as PokemonCard);
}

function cardKey(card: AnyCard): string {
  if (selectedTcg.value === "naruto") return naruto.cardKey(card as NarutoCard);
  return pokemon.cardKey(card as PokemonCard);
}

const searchQuery = ref("");
const activeFilter = ref("all");
const deckName = ref(sessionStorage.getItem(STORAGE_NAME_KEY) || "");
const editingDeckId = ref<number | null>(
  sessionStorage.getItem(STORAGE_ID_KEY)
    ? Number(sessionStorage.getItem(STORAGE_ID_KEY))
    : null,
);
const saving = ref(false);

const tooltipImage = ref("");
const tooltipPos = ref({ x: 0, y: 0 });
const tooltipVisible = ref(false);

function showTooltip(card: AnyCard, event: MouseEvent) {
  tooltipImage.value = getCardImage(card);
  tooltipPos.value = { x: event.clientX + 16, y: event.clientY + 16 };
  tooltipVisible.value = true;
}

function moveTooltip(event: MouseEvent) {
  tooltipPos.value = { x: event.clientX + 16, y: event.clientY + 16 };
}

function hideTooltip() {
  tooltipVisible.value = false;
}

const mobileThumbnail = ref<string | null>(null);
let pressTimer: ReturnType<typeof setTimeout> | null = null;

function onTouchStart(card: AnyCard) {
  if (!isMobile.value) return;
  pressTimer = setTimeout(() => {
    mobileThumbnail.value = getCardImage(card);
  }, 300);
}

function onTouchEnd() {
  if (pressTimer) {
    clearTimeout(pressTimer);
    pressTimer = null;
  }
  mobileThumbnail.value = null;
}

function onTouchMove() {
  if (pressTimer) {
    clearTimeout(pressTimer);
    pressTimer = null;
  }
  mobileThumbnail.value = null;
}

const chosenCards = ref<AnyCard[]>(
  JSON.parse(sessionStorage.getItem(STORAGE_KEY) || "[]"),
);

// Clear deck when TCG switches
watch(selectedTcg, () => {
  chosenCards.value = [];
  activeFilter.value = "all";
  searchQuery.value = "";
});

watch(
  chosenCards,
  (val) => {
    sessionStorage.setItem(STORAGE_KEY, JSON.stringify(val));
  },
  { deep: true },
);

watch(deckName, (val) => {
  sessionStorage.setItem(STORAGE_NAME_KEY, val);
});

const filteredCards = computed(() => {
  if (selectedTcg.value === "naruto") {
    return naruto.filterCards(
      naruto.allCards.value,
      searchQuery.value,
      activeFilter.value,
    );
  }
  return pokemon.filterCards(
    pokemon.allCards.value,
    searchQuery.value,
    activeFilter.value,
  );
});

function addCard(card: AnyCard) {
  chosenCards.value.push(card);
}

function removeCard(card: AnyCard) {
  const key = cardKey(card);
  const idx = chosenCards.value.findIndex((c) => cardKey(c) === key);
  if (idx !== -1) chosenCards.value.splice(idx, 1);
  hideTooltip();
}

function clearDeck() {
  chosenCards.value = [];
}

interface CardGroup {
  card: AnyCard;
  key: string;
  count: number;
}

const groupedChosenCards = computed<CardGroup[]>(() => {
  const map = new Map<string, CardGroup>();
  for (const card of chosenCards.value) {
    const k = cardKey(card);
    const existing = map.get(k);
    if (existing) {
      existing.count++;
    } else {
      map.set(k, { card, key: k, count: 1 });
    }
  }
  return [...map.values()];
});

async function saveDeck() {
  if (!deckName.value.trim() || chosenCards.value.length === 0) return;
  saving.value = true;
  try {
    if (selectedTcg.value === "naruto") {
      await naruto.saveDeck(
        user.value!.id,
        deckName.value.trim(),
        chosenCards.value as NarutoCard[],
      );
    } else {
      await pokemon.saveDeck(
        user.value!.id,
        deckName.value.trim(),
        chosenCards.value as PokemonCard[],
      );
    }
    chosenCards.value = [];
    deckName.value = "";
    editingDeckId.value = null;
    sessionStorage.removeItem(STORAGE_KEY);
    sessionStorage.removeItem(STORAGE_NAME_KEY);
    sessionStorage.removeItem(STORAGE_ID_KEY);
  } catch (e) {
    console.error("Failed to save deck:", e);
  } finally {
    saving.value = false;
  }
}

async function importDeck() {
  try {
    const text = await navigator.clipboard.readText();
    let imported: AnyCard[];
    if (selectedTcg.value === "naruto") {
      imported = naruto.parseDeckList(text, naruto.allCards.value);
    } else {
      imported = pokemon.parseDeckList(text, pokemon.allCards.value);
    }
    if (imported.length > 0) {
      chosenCards.value = imported;
    }
  } catch (e) {
    console.error("Failed to import deck:", e);
  }
}

function cardName(card: AnyCard): string {
  if ("Name" in card) return card.Name ?? "";
  if ("name" in card) return card.name ?? "";
  return "";
}

// Responsive detection
const isNarrow = ref(window.innerWidth < 1280);
const isMobile = ref(window.innerWidth <= 480);
const drawerOpen = ref(false);

function onResize() {
  isNarrow.value = window.innerWidth < 1280;
  isMobile.value = window.innerWidth <= 480;
}
onMounted(() => window.addEventListener("resize", onResize));
onUnmounted(() => window.removeEventListener("resize", onResize));

// Draggable divider
const splitPercent = ref(40);
const dragging = ref(false);
const builderEl = ref<HTMLElement | null>(null);

function startDrag() {
  dragging.value = true;
  document.addEventListener("mousemove", onDrag);
  document.addEventListener("mouseup", stopDrag);
}

function onDrag(e: MouseEvent) {
  if (!dragging.value || !builderEl.value) return;
  const rect = builderEl.value.getBoundingClientRect();
  const isVertical = window.innerWidth < 1280;
  let percent: number;
  if (isVertical) {
    percent = ((e.clientY - rect.top) / rect.height) * 100;
  } else {
    percent = ((e.clientX - rect.left) / rect.width) * 100;
  }
  splitPercent.value = Math.max(20, Math.min(80, percent));
}

function stopDrag() {
  dragging.value = false;
  document.removeEventListener("mousemove", onDrag);
  document.removeEventListener("mouseup", stopDrag);
}
</script>

<template>
  <div
    ref="builderEl"
    class="builder"
    :class="{ dragging }"
    :style="{
      '--split': splitPercent + '%',
    }"
  >
    <!-- Deck panel (top on narrow, right on wide) -->
    <div
      class="builder-panel deck-panel"
      :class="{ 'drawer-open': drawerOpen }"
    >
      <div
        v-if="isMobile"
        class="drawer-handle"
        @click="drawerOpen = !drawerOpen"
      >
        <div class="drawer-handle-bar" />
      </div>
      <div class="panel-header" @click="isMobile && (drawerOpen = !drawerOpen)">
        <span class="panel-title">Deck ({{ chosenCards.length }})</span>
      </div>

      <div class="deck-actions">
        <v-text-field
          v-model="deckName"
          variant="outlined"
          density="compact"
          placeholder="Deck name..."
          hide-details
          class="deck-name-input"
        />
        <v-btn
          color="deep-purple-accent-2"
          :disabled="!deckName.trim() || chosenCards.length === 0 || saving"
          :loading="saving"
          @click="saveDeck"
        >
          Save Deck
        </v-btn>
        <v-btn
          variant="outlined"
          color="deep-purple-accent-2"
          @click="importDeck"
        >
          <v-icon start>mdi-import</v-icon>
          Import
        </v-btn>
        <v-btn
          variant="outlined"
          color="red"
          :disabled="chosenCards.length === 0"
          @click="clearDeck"
        >
          <v-icon start>mdi-delete-sweep</v-icon>
          Clear
        </v-btn>
      </div>

      <div class="deck-grid">
        <div
          v-for="group in groupedChosenCards"
          :key="group.key"
          class="card-stack"
          @click="removeCard(group.card)"
          @mouseenter="showTooltip(group.card, $event)"
          @mousemove="moveTooltip($event)"
          @mouseleave="hideTooltip"
          @touchstart="onTouchStart(group.card)"
          @touchend="onTouchEnd"
          @touchcancel="onTouchEnd"
          @touchmove="onTouchMove"
        >
          <!-- Narrow: single image per card -->
          <template v-if="isNarrow">
            <img
              :src="getCardImage(group.card)"
              :alt="cardName(group.card)"
              class="card-image card-chosen"
            />
          </template>
          <!-- Wide: stacked images -->
          <template v-else>
            <img
              v-for="n in group.count"
              :key="n"
              :src="getCardImage(group.card)"
              :alt="cardName(group.card)"
              class="card-image card-chosen"
              :style="{
                position: 'absolute',
                bottom: `${(n - 4) * 0}px`,
                right: `${(n - 5) * 5}px`,
                zIndex: n,
              }"
            />
          </template>
          <span v-if="group.count > 1" class="card-count">{{
            group.count
          }}</span>
        </div>
      </div>

      <div v-if="chosenCards.length === 0" class="empty-state">
        Click cards to add them to your deck
      </div>
    </div>

    <!-- Draggable divider -->
    <div v-if="!isMobile" class="divider" @mousedown.prevent="startDrag">
      <div class="divider-handle" />
    </div>

    <!-- Card pool panel (bottom on narrow, left on wide) -->
    <div class="builder-panel pool-panel">
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
          @mouseenter="showTooltip(card, $event)"
          @mousemove="moveTooltip($event)"
          @mouseleave="hideTooltip"
          @touchstart="onTouchStart(card)"
          @touchend="onTouchEnd"
          @touchcancel="onTouchEnd"
          @touchmove="onTouchMove"
        >
          <img
            :src="getCardImage(card)"
            :alt="cardName(card)"
            loading="lazy"
            class="card-image"
          />
        </div>
      </div>
    </div>

    <!-- Card tooltip (desktop) -->
    <div
      v-if="tooltipVisible && tooltipImage"
      class="card-tooltip"
      :style="{ top: tooltipPos.y + 'px', left: tooltipPos.x + 'px' }"
    >
      <img :src="tooltipImage" class="card-tooltip-img" />
    </div>

    <!-- Mobile thumbnail overlay -->
    <div
      v-if="mobileThumbnail"
      class="mobile-thumbnail-overlay"
      @touchend="onTouchEnd"
    >
      <img :src="mobileThumbnail" class="mobile-thumbnail-img" />
    </div>
  </div>
</template>

<style scoped>
.builder {
  display: flex;
  gap: 0;
  padding: 12px 16px;
  height: calc(100vh - 64px);
  box-sizing: border-box;
}

.builder.dragging {
  user-select: none;
  cursor: col-resize;
}

/* Wide layout: pool left, deck right */
.pool-panel {
  flex: 1;
  order: 1;
}

.divider {
  order: 2;
  width: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: col-resize;
  flex-shrink: 0;
}

.divider-handle {
  width: 4px;
  height: 48px;
  background: var(--border);
  border-radius: 2px;
  transition: background 0.2s;
}

.divider:hover .divider-handle {
  background: var(--accent);
}

.deck-panel {
  flex: 1;
  order: 3;
}

.builder-panel {
  display: flex;
  flex-direction: column;
  gap: 12px;
  background: var(--code-bg);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 14px 12px;
  box-shadow: var(--shadow);
  overflow: hidden;
}

.toolbar {
  display: flex;
  flex-direction: column;
  gap: 12px;
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

.deck-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

.deck-name-input {
  max-width: 240px;
  font-family: var(--sans);
}

.card-grid {
  flex: 1;
  display: grid;
  grid-template-columns: repeat(6, minmax(90px, 1fr));
  grid-auto-rows: min-content;
  gap: 12px;
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
  border-radius: 4px;
  transition: box-shadow 0.2s;
}

.card-image:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.deck-grid .card-image:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.card-stack:has(.card-image:hover) {
  z-index: 10;
}

.card-stack .card-image {
  position: absolute;
  width: 100%;
  height: 100%;
  object-fit: contain;
  object-position: top left;
}

.deck-grid {
  flex: 1;
  min-height: 0;
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  grid-auto-rows: min-content;
  gap: 12px 16px;
  overflow-y: auto;
  padding-right: 30px;
  padding-bottom: 20px;
}

.card-stack {
  position: relative;
  aspect-ratio: 5 / 7;
  min-height: 0;
  cursor: pointer;
}

.card-count {
  position: absolute;
  bottom: 15%;
  right: 10%;
  background: #7c4dff;
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
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.4);
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

.card-tooltip {
  position: fixed;
  z-index: 9999;
  pointer-events: none;
  width: 300px;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
}

.card-tooltip-img {
  width: 100%;
  height: auto;
  display: block;
  border-radius: 8px;
}

/* Narrow layout: deck on top, pool on bottom, vertical split */
@media (max-width: 1279px) {
  .builder {
    flex-direction: column;
  }

  .builder.dragging {
    cursor: row-resize;
  }

  .deck-panel {
    order: 1;
    flex: none;
    height: var(--split);
  }

  .divider {
    order: 2;
    width: 100%;
    height: 12px;
    cursor: row-resize;
  }

  .divider-handle {
    width: 48px;
    height: 4px;
  }

  .pool-panel {
    order: 3;
    flex: 1;
    min-height: 0;
  }

  .card-grid {
    grid-template-columns: repeat(8, minmax(50px, 1fr));
    gap: 8px;
  }

  .deck-grid {
    grid-template-columns: repeat(6, 1fr);
    grid-auto-rows: min-content;
    overflow-y: auto;
  }

  .card-stack {
    aspect-ratio: 5 / 7;
  }

  .card-stack .card-image {
    position: relative;
    object-fit: cover;
    object-position: center;
  }
}

/* Mobile layout: collapsible deck drawer */
@media (max-width: 480px) {
  .builder {
    flex-direction: column;
    padding: 8px;
    padding-bottom: 0;
    height: calc(100vh - 76px);
    position: relative;
    overflow: hidden;
  }

  .divider {
    display: none;
  }

  .pool-panel {
    order: 1;
    flex: 1;
    min-height: 0;
    border-radius: 12px;
    padding: 8px 4px;
    margin-bottom: 48px;
    overflow: hidden;
  }

  .deck-panel {
    position: absolute;
    bottom: env(safe-area-inset-bottom, 20px);
    left: 0;
    right: 0;
    height: 60%;
    max-height: 60%;
    border-radius: 16px 16px 0 0;
    padding: 10px 8px;
    transform: translateY(calc(100% - 44px));
    transition: transform 0.3s ease;
    z-index: 100;
    box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.3);
  }

  .deck-panel.drawer-open {
    transform: translateY(0);
  }

  .card-grid {
    grid-template-columns: repeat(4, minmax(0, 1fr));
    gap: 4px;
    overflow-y: auto;
    overflow-x: hidden;
    padding: 0;
  }

  .deck-grid {
    grid-template-columns: repeat(4, 1fr);
    gap: 6px 8px;
    padding-right: 4px;
  }

  .deck-actions {
    flex-wrap: wrap;
    gap: 8px;
  }

  .deck-name-input {
    max-width: 100%;
    flex: 1 1 100%;
  }

  .search-wrapper {
    max-width: 100%;
  }

  .card-tooltip {
    display: none;
  }

  .mobile-thumbnail-overlay {
    position: fixed;
    inset: 0;
    z-index: 9999;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(0, 0, 0, 0.6);
  }

  .mobile-thumbnail-img {
    width: 70%;
    max-width: 280px;
    border-radius: 10px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
  }

  .drawer-handle {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 18px 0 4px;
    cursor: pointer;
  }

  .drawer-handle-bar {
    width: 36px;
    height: 4px;
    background: var(--border);
    border-radius: 2px;
  }
}
</style>
