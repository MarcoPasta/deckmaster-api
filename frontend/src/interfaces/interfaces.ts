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
