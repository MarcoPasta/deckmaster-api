import { ref } from "vue";
import type { PokemonCard } from "../../interfaces/interfaces";
import { useApi } from "../useApi";

export function usePokemonBuilder() {
  const { get, post } = useApi();

  const allCards = ref<PokemonCard[]>([]);
  const loading = ref(false);
  const filters = ["All", "Pokémon", "Trainer", "Energy"];

  async function fetchCards() {
    if (allCards.value.length > 0) return;
    loading.value = true;
    try {
      const data = await get<{ content: PokemonCard[] }>("/pokemon/cards");
      allCards.value = data.content;
    } catch (e) {
      console.error("Failed to fetch cards:", e);
    } finally {
      loading.value = false;
    }
  }

  function getCardImage(card: PokemonCard): string {
    if (!card.CardNumber || !card.Code) return "";
    const number = card.CardNumber.split("/")[0];
    return `/assets/cards/de/${card.Code}/${card.Code}_${number}.png`;
  }

  function cardKey(card: PokemonCard): string {
    return `${card.Code}_${card.CardNumber}`;
  }

  function filterCards(
    cards: PokemonCard[],
    searchQuery: string,
    activeFilter: string,
  ): PokemonCard[] {
    return cards.filter((card) => {
      const matchesSearch = (card.Name ?? "")
        .toLowerCase()
        .includes(searchQuery.toLowerCase());
      let matchesFilter = activeFilter === "all";
      if (!matchesFilter) {
        const type = (card.CardType ?? "").toLowerCase();
        if (activeFilter === "energy") {
          matchesFilter = type === "energie" || type === "energy";
        } else {
          matchesFilter = type === activeFilter.toLowerCase();
        }
      }
      return matchesSearch && matchesFilter;
    });
  }

  function parseDeckList(text: string, cardPool: PokemonCard[]): PokemonCard[] {
    const lines = text
      .split("\n")
      .map((l) => l.trim())
      .filter(Boolean);
    const imported: PokemonCard[] = [];

    for (const line of lines) {
      if (/^(Pokémon|Pokemon|Trainer|Energy):\s*\d+$/i.test(line)) continue;

      const match = line.match(/^(\d+)\s+(.+?)\s+([A-Z]{2,})\s+(\d+)$/);
      if (!match) continue;

      const count = parseInt(match[1], 10);
      const name = match[2];
      const code = match[3].toLowerCase();
      const number = match[4];

      let card: PokemonCard | undefined;

      if (name.endsWith("-Energie")) {
        const energyType = name.replace("-Energie", "");
        card = cardPool.find(
          (c) =>
            c.Code === code &&
            c.CardNumber?.split("/")[0] === number &&
            c.EnergyType?.toLowerCase() === energyType.toLowerCase(),
        );
      } else {
        card = cardPool.find(
          (c) => c.Code === code && c.CardNumber?.split("/")[0] === number,
        );
      }

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
    cards: PokemonCard[],
  ): Promise<void> {
    await post("/auth/pokemon/decks/create", {
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
