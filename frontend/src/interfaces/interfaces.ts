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
  ID: number;
  CreatedAt: Date;
  UpdatedAt: Date;
  DeletedAt: Date | null;
  DeckId: number;
  CardType: string;
  Phase: string;
  Name: string;
  Kp: number;
  EnergyType: string;
  Ability: string;
  Attacks: Attack[];
  AttacksEN: Attack[] | null;
  Weakness: string;
  Resistance: string;
  RetreatCost: number;
  CardNumber: string;
  Code: string;
  Set: string;
}

export interface Attack {
  cost: string[];
  name: string;
  effect: string;
  damage: string;
}

export interface PokemonDecks {
  ID: number;
  CreatedAt: Date;
  UpdatedAt: Date;
  DeletedAt: Date | null;
  UserId: string;
  Name: string;
  Cards: PokemonCard[];
}
