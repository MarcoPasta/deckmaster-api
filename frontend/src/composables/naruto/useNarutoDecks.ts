import { ref, watch } from "vue";
import type { NarutoCard, NarutoDecks } from "../../interfaces/interfaces";
import { useApi } from "../useApi";

export function useNarutoDecks(userId: () => string | undefined) {
  const { get } = useApi();
  const decks = ref<NarutoDecks[]>([]);

  async function fetchDecks(id: string) {
    try {
      const data = await get<{ content: NarutoDecks[] }>(
        "/auth/naruto/decks/user",
        { id },
      );
      decks.value = data.content;
    } catch (e) {
      console.error("Failed to fetch naruto decks:", e);
    }
  }

  watch(
    userId,
    (id) => {
      if (id) fetchDecks(id);
    },
    { immediate: true },
  );

  function getPreviewImage(deck: NarutoDecks): string {
    const entry = deck.entries?.[0];
    if (!entry?.card?.card_number || !entry?.card?.rarity) return "";
    return `/assets/cards/en/naruto/${entry.card.card_number}_${entry.card.rarity}.webp`;
  }

  function getCardImage(card: NarutoCard): string {
    if (!card.card_number || !card.rarity) return "";
    return `/assets/cards/en/naruto/${card.card_number}_${card.rarity}.webp`;
  }

  function buildDeckText(deck: NarutoDecks): string {
    const entries = deck.entries ?? [];

    const cardCounts = new Map<number, { card: NarutoCard; count: number }>();
    for (const entry of entries) {
      const existing = cardCounts.get(entry.card.ID);
      if (existing) {
        existing.count++;
      } else {
        cardCounts.set(entry.card.ID, { card: entry.card, count: 1 });
      }
    }

    const lines: string[] = [];
    const total = entries.length;
    lines.push(`Naruto Mythos: ${total}`);
    for (const { card, count } of cardCounts.values()) {
      lines.push(`${count} ${card.name} ${card.card_number} ${card.rarity}`);
    }

    return lines.join("\n").trim();
  }

  function deleteDeckLocal(deckId: number) {
    decks.value = decks.value.filter((d) => d.ID !== deckId);
  }

  return {
    decks,
    fetchDecks,
    getPreviewImage,
    getCardImage,
    buildDeckText,
    deleteDeckLocal,
  };
}
