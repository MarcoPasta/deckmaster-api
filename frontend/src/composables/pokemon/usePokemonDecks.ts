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

  async function copyDeckAsImage(deck: PokemonDecks): Promise<Blob | null> {
    const entries = deck.entries ?? [];

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
    try {
      const patternImg = await new Promise<HTMLImageElement>(
        (resolve, reject) => {
          const img = new Image();
          img.onload = () => resolve(img);
          img.onerror = reject;
          img.src = "/assets/pattern-bg-dark.svg";
        },
      );
      const pattern = ctx.createPattern(patternImg, "repeat");
      if (pattern) {
        ctx.fillStyle = pattern;
        ctx.fillRect(0, 0, canvasWidth, canvasHeight);
      } else {
        ctx.fillStyle = "#16171d";
        ctx.fillRect(0, 0, canvasWidth, canvasHeight);
      }
    } catch {
      ctx.fillStyle = "#16171d";
      ctx.fillRect(0, 0, canvasWidth, canvasHeight);
    }

    // Load card images
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

      for (let n = 0; n < count; n++) {
        const x = baseX + n * stackOffset;
        const y = baseY + n * stackOffset;
        ctx.drawImage(img, x, y, cardWidth, cardHeight);
      }

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

    return new Promise((resolve) => {
      canvas.toBlob((blob) => resolve(blob), "image/png");
    });
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
    copyDeckAsImage,
    deleteDeckLocal,
  };
}
