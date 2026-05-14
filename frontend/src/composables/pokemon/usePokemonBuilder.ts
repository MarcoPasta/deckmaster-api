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

  async function parseDeckList(text: string): Promise<PokemonCard[]> {
    const cardPool: PokemonCard[] = [];
    try {
      const parsedList = await get<{ content: PokemonCard[] }>(
        "/pokemon/cards/names",
        {
          decklist: text,
        },
      );
      cardPool.push(...parsedList.content);
    } catch (e) {
      console.error("Failed to fetch cards:", e);
    } finally {
      loading.value = false;
    }

    // const lines = text
    //   .split("\n")
    //   .map((l) => l.trim())
    //   .filter(Boolean);
    // const imported: PokemonCard[] = [];
    // const queue: string[] = [];

    // for (const line of lines) {
    //   // console.log(`Processing line: "${line}"`);

    //   if (/^(Pokémon|Pokemon|Trainer|Energy):\s*\d+$/i.test(line)) continue;

    //   let count: number;
    //   let name: string;
    //   let code: string;
    //   let number: string;

    //   const match = line.match(/^(\d+)\s+(.+?)\s+([A-Z]{2,})\s+(\d+)$/);
    //   const matchEnergy = line.match(/^(\d+)\s+(.+?)\s+(\d+)$/);
    //   if (match) {
    //     count = parseInt(match[1], 10);
    //     name = match[2];
    //     code = match[3].toLowerCase();
    //     number = match[4];
    //     // console.log("##### matches card: ", name, code, number);
    //   } else if (matchEnergy) {
    //     count = parseInt(matchEnergy[1], 10);
    //     name = matchEnergy[2];
    //     number = matchEnergy[3];
    //     code = "mee";
    //     // console.log("##### matches energy: ", name);
    //   } else {
    //     console.warn(`Line did not match expected format: "${line}"`);
    //     continue;
    //   }

    //   let card: PokemonCard | undefined;

    //   // advanced check for energy cards to pattern match with any energy card
    //   if (name.startsWith("Basic") && name.endsWith("Energy")) {
    //     // console.log(`Found energy card line: "${name}"`);
    //     const match = name.match(/^Basic\s+{(.+?)}.+|$/i);
    //     if (match) {
    //       ({ code, number } = findEnergyCard(match));
    //     }
    //     card = cardPool.find(
    //       (c) => c.Code === code && c.CardNumber?.split("/")[0] === number,
    //     );
    //   } else {
    //     // search for any other non energy card
    //     card = cardPool.find(
    //       (c) => c.Code === code && c.CardNumber?.split("/")[0] === number,
    //     );
    //     //  if the card is not found, search for a substitue in the database with the same name and type
    //     if (!card) {
    //       console.warn(`No card found for line: "${line}", adding to queue"`);
    //       queue.push(name);
    //     }

    //     console.log("Cards that could not be found: ", queue);

    //     if (card) {
    //       for (let i = 0; i < count; i++) {
    //         imported.push(card);
    //       }
    //     }
    //   }
    // }
    // const missing = await fetchMissingCards(queue);
    // imported.push(...missing);
    return cardPool;
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
