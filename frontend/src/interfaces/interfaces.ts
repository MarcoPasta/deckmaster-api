export interface DigimonCards {
  id: string;
  deckId: number;
  name: string;
  cardType: string;
  color: string[];
  cost: number;
  dp: number;
  dvCost: string;
  form: string;
  attribute: string;
  species: string[];
  text: string[];
  set: string;
  createdAt: Date;
  updatedAt: Date;
  deletedAt: Date | null;
}

export interface PokemonCard {
  id: number;
  createdAt: Date;
  updatedAt: Date;
  deletedAt: Date | null;
  deckId: number; // Foreign key to PokemonDecks
  cardType: string;
  phase: string;
  name: string;
  kp: number;
  energyType: string;
  ability: string;
  attacks: Attack[];
  attacksEN: Attack[];
  weakness: string;
  resistance: string;
  retreatCost: number;
  cardNumber: string;
  set: string;
}

export interface Attack {
  cost: string[];
  name: string;
  effect: string;
  damage: string;
}

export interface PokemonDecks {
  id: number;
  createdAt: Date;
  updatedAt: Date;
  deletedAt: Date | null;
  userId: string;
  name: string;
  cards: PokemonCard[];
}
