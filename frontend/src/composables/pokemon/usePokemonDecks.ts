import { ref, watch } from "vue";
import type { PokemonCard, PokemonDecks } from "../../interfaces/interfaces";
import { useApi } from "../useApi";

export function usePokemonDecks(userId: () => string | undefined) {
  const { get } = useApi();
  const decks = ref<PokemonDecks[]>([]);

  async function fetchDecks(id: string) {
    try {
      const data = await get<{ content: PokemonDecks[] }>(
        "/auth/pokemon/decks/user",
        { id },
      );
      decks.value = data.content;
    } catch (e) {
      console.error("Failed to fetch decks:", e);
    }
  }

  watch(
    userId,
    (id) => {
      if (id) fetchDecks(id);
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
