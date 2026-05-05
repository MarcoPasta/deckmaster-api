<script setup lang="ts">
import { useUser } from "@clerk/vue";
import { ref, watch } from "vue";
import { useRouter } from "vue-router";
import { useApi } from "../composables/useApi";
import type { PokemonCard, PokemonDecks } from "../interfaces/interfaces";

const router = useRouter();
const { user } = useUser();
const { get } = useApi();

const decks = ref<PokemonDecks[]>([]);

async function fetchDecks(userId: string) {
  try {
    const data = await get<{ content: PokemonDecks[] }>(
      "/auth/pokemon/decks/user",
      { id: userId },
    );
    decks.value = data.content;
  } catch (e) {
    console.error("Failed to fetch decks:", e);
  }
}

watch(
  () => user.value?.id,
  (userId) => {
    if (userId) fetchDecks(userId);
  },
  { immediate: true },
);

function getPreviewImage(deck: PokemonDecks): string {
  const entry = deck.entries?.[0];
  if (!entry?.card?.CardNumber || !entry?.card?.Code) return "";
  const number = entry.card.CardNumber.split("/")[0];
  return `/assets/cards/de/${entry.card.Code}/${entry.card.Code}_${number}.png`;
}

function getCardImage(card: PokemonCard): string {
  if (!card.CardNumber || !card.Code) return "";
  const number = card.CardNumber.split("/")[0];
  return `/assets/cards/de/${card.Code}/${card.Code}_${number}.png`;
}

function formatDate(date: Date | string): string {
  return new Date(date).toLocaleDateString("de-DE", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
  });
}

const dialogOpen = ref(false);
const selectedDeck = ref<PokemonDecks | null>(null);

function openDeck(deck: PokemonDecks) {
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

function buildDeckText(deck: PokemonDecks): string {
  const entries = deck.entries ?? [];
  const groups: Record<string, { card: PokemonCard; count: number }[]> = {
    Pokémon: [],
    Trainer: [],
    Energy: [],
  };

  const cardCounts = new Map<number, { card: PokemonCard; count: number }>();
  for (const entry of entries) {
    const existing = cardCounts.get(entry.card.ID);
    if (existing) {
      existing.count++;
    } else {
      cardCounts.set(entry.card.ID, { card: entry.card, count: 1 });
    }
  }

  for (const { card, count } of cardCounts.values()) {
    const type = card.CardType?.toLowerCase();
    if (type === "pokemon" || type === "pokémon") {
      groups["Pokémon"].push({ card, count });
    } else if (type === "trainer") {
      groups["Trainer"].push({ card, count });
    } else if (type === "energy" || type === "energie") {
      groups["Energy"].push({ card, count });
    } else {
      groups["Trainer"].push({ card, count });
    }
  }

  const lines: string[] = [];
  for (const [label, cards] of Object.entries(groups)) {
    if (cards.length === 0) continue;
    const total = cards.reduce((sum, c) => sum + c.count, 0);
    lines.push(`${label}: ${total}`);
    for (const { card, count } of cards) {
      const number = card.CardNumber?.split("/")[0] ?? "";
      const code = (card.Code ?? "").toUpperCase();
      const cardType = (card.CardType ?? "").toLowerCase();
      const displayName =
        cardType === "energie" || cardType === "energy"
          ? `${card.EnergyType}-Energie`
          : card.Name;
      lines.push(`${count} ${displayName} ${code} ${number}`);
    }
    lines.push("");
  }

  return lines.join("\n").trim();
}

async function copyAsText() {
  if (!selectedDeck.value) return;
  const text = buildDeckText(selectedDeck.value);
  await navigator.clipboard.writeText(text);
  shareMenuOpen.value = false;
}

async function copyAsImage() {
  if (!selectedDeck.value) return;
  const deck = selectedDeck.value;
  const entries = deck.entries ?? [];

  // Group cards by ID
  const cardGroups: { card: PokemonCard; count: number }[] = [];
  const groupMap = new Map<number, { card: PokemonCard; count: number }>();
  for (const entry of entries) {
    const existing = groupMap.get(entry.card.ID);
    if (existing) {
      existing.count++;
    } else {
      const group = { card: entry.card, count: 1 };
      groupMap.set(entry.card.ID, group);
      cardGroups.push(group);
    }
  }

  const cardWidth = 120;
  const cardHeight = Math.round(cardWidth * (7 / 5));
  const cols = 6;
  const gap = 16;
  const stackOffset = 6;
  const padding = 20;
  const rows = Math.ceil(cardGroups.length / cols);
  const maxStack = Math.max(...cardGroups.map((g) => g.count), 1);
  const extraOffset = (maxStack - 1) * stackOffset;
  const canvasWidth =
    padding * 2 + cols * cardWidth + (cols - 1) * gap + extraOffset;
  const canvasHeight =
    padding * 2 + rows * (cardHeight + extraOffset) + (rows - 1) * gap;

  const canvas = document.createElement("canvas");
  canvas.width = canvasWidth;
  canvas.height = canvasHeight;
  const ctx = canvas.getContext("2d")!;

  // Draw tiled pattern background
  const patternImg = await new Promise<HTMLImageElement>((resolve, reject) => {
    const img = new Image();
    img.onload = () => resolve(img);
    img.onerror = reject;
    img.src = "/assets/pattern-bg-dark.svg";
  });
  const pattern = ctx.createPattern(patternImg, "repeat");
  if (pattern) {
    ctx.fillStyle = pattern;
    ctx.fillRect(0, 0, canvasWidth, canvasHeight);
  } else {
    ctx.fillStyle = "#16171d";
    ctx.fillRect(0, 0, canvasWidth, canvasHeight);
  }

  // Load unique card images
  const loadImage = (src: string): Promise<HTMLImageElement> => {
    return new Promise((resolve, reject) => {
      const img = new Image();
      img.crossOrigin = "anonymous";
      img.onload = () => resolve(img);
      img.onerror = reject;
      img.src = src;
    });
  };

  const images = await Promise.all(
    cardGroups.map(({ card }) => {
      const src = getCardImage(card);
      return src ? loadImage(src).catch(() => null) : Promise.resolve(null);
    }),
  );

  for (let i = 0; i < cardGroups.length; i++) {
    const img = images[i];
    if (!img) continue;
    const { count } = cardGroups[i];
    const col = i % cols;
    const row = Math.floor(i / cols);
    const baseX = padding + col * (cardWidth + gap);
    const baseY = padding + row * (cardHeight + extraOffset + gap);

    // Draw stacked copies
    for (let n = 0; n < count; n++) {
      const x = baseX + n * stackOffset;
      const y = baseY + n * stackOffset;
      ctx.drawImage(img, x, y, cardWidth, cardHeight);
    }

    // Draw count badge if more than 1
    if (count > 1) {
      const badgeX = baseX + (count - 1) * stackOffset + cardWidth - 8;
      const badgeY = baseY + (count - 1) * stackOffset + cardHeight - 8;
      const text = String(count);
      ctx.font = "bold 14px sans-serif";
      const textWidth = ctx.measureText(text).width;
      const badgeWidth = Math.max(22, textWidth + 10);
      const badgeHeight = 22;

      ctx.fillStyle = "#7c4dff";
      ctx.beginPath();
      ctx.roundRect(
        badgeX - badgeWidth,
        badgeY - badgeHeight,
        badgeWidth,
        badgeHeight,
        11,
      );
      ctx.fill();

      ctx.fillStyle = "#ffffff";
      ctx.textAlign = "center";
      ctx.textBaseline = "middle";
      ctx.fillText(text, badgeX - badgeWidth / 2, badgeY - badgeHeight / 2);
    }
  }

  canvas.toBlob(async (blob) => {
    if (!blob) return;
    await navigator.clipboard.write([new ClipboardItem({ "image/png": blob })]);
    shareMenuOpen.value = false;
  }, "image/png");
}

function duplicateDeck() {
  if (!selectedDeck.value) return;
  // TODO: call API to duplicate
  dialogOpen.value = false;
}

function deleteDeck() {
  if (!selectedDeck.value) return;
  decks.value = decks.value.filter((d) => d.ID !== selectedDeck.value!.ID);
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
              :alt="entry.card.Name"
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
