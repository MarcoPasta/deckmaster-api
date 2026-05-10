import { ref } from "vue";
import type { NarutoCard } from "../../interfaces/interfaces";
import { useApi } from "../useApi";

export function useNarutoBuilder() {
  const { get, post } = useApi();

  const allCards = ref<NarutoCard[]>([]);
  const loading = ref(false);
  const filters = ["All", "Ninja", "Jutsu", "Mission", "Client"];

  async function fetchCards() {
    if (allCards.value.length > 0) return;
    loading.value = true;
    try {
      const data = await get<{ content: NarutoCard[] }>("/naruto/cards");
      allCards.value = data.content;
    } catch (e) {
      console.error("Failed to fetch naruto cards:", e);
    } finally {
      loading.value = false;
    }
  }

  function getCardImage(card: NarutoCard): string {
    if (!card.card_number || !card.rarity) return "";
    return `/assets/cards/en/naruto/${card.card_number}_${card.rarity}.webp`;
  }

  function cardKey(card: NarutoCard): string {
    return `${card.card_number}_${card.rarity}`;
  }

  function filterCards(
    cards: NarutoCard[],
    searchQuery: string,
    activeFilter: string,
  ): NarutoCard[] {
    return cards.filter((card) => {
      const matchesSearch = (card.name ?? "")
        .toLowerCase()
        .includes(searchQuery.toLowerCase());
      const matchesFilter =
        activeFilter === "all" ||
        (card.keyword ?? "").toLowerCase() === activeFilter.toLowerCase();
      return matchesSearch && matchesFilter;
    });
  }

  function parseDeckList(text: string, cardPool: NarutoCard[]): NarutoCard[] {
    const lines = text
      .split("\n")
      .map((l) => l.trim())
      .filter(Boolean);
    const imported: NarutoCard[] = [];

    for (const line of lines) {
      // Skip section headers
      if (/^.+:\s*\d+$/i.test(line) && !line.match(/^\d/)) continue;

      // Parse lines like "4 Naruto Uzumaki 01 C"
      const match = line.match(
        /^(\d+)\s+(.+?)\s+(\d+|M\d+|Legendary)\s+(\w+)$/,
      );
      if (!match) continue;

      const count = parseInt(match[1], 10);
      const cardNumber = match[3];
      const rarity = match[4];

      const card = cardPool.find(
        (c) => c.card_number === cardNumber && c.rarity === rarity,
      );

      if (card) {
        for (let i = 0; i < count; i++) {
          imported.push(card);
        }
      }
    }

    return imported;
  }

  async function saveDeck(
    userId: string,
    name: string,
    cards: NarutoCard[],
  ): Promise<void> {
    await post("/auth/naruto/decks/create", {
      UserId: userId,
      Name: name,
      Entries: cards.map((card) => ({ cardId: card.ID })),
    });
  }

  return {
    allCards,
    loading,
    fetchCards,
    filters,
    getCardImage,
    cardKey,
    filterCards,
    parseDeckList,
    saveDeck,
  };
}
