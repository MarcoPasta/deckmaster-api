import datetime
from typing import Optional

from common.database import Base
from sqlalchemy import TIMESTAMP, BigInteger, SmallInteger, Text
from sqlalchemy.orm import Mapped, mapped_column


class DigimonCards(Base):
    __tablename__ = "digimon_cards"
    id: Mapped[str] = mapped_column(Text, primary_key=True)
    name: Mapped[Optional[str]] = mapped_column(Text)
    card_type: Mapped[Optional[str]] = mapped_column(Text)
    color: Mapped[Optional[str]] = mapped_column(Text)
    cost: Mapped[Optional[int]] = mapped_column(SmallInteger)
    dp: Mapped[Optional[int]] = mapped_column(SmallInteger)
    dv_cost: Mapped[Optional[str]] = mapped_column(Text)
    form: Mapped[Optional[str]] = mapped_column(Text)
    attribute: Mapped[Optional[str]] = mapped_column(Text)
    species: Mapped[Optional[str]] = mapped_column(Text)
    text: Mapped[Optional[str]] = mapped_column(Text)
    set: Mapped[Optional[str]] = mapped_column(Text)
    created_at: Mapped[Optional[datetime.datetime]] = mapped_column(TIMESTAMP(timezone=True))
    updated_at: Mapped[Optional[datetime.datetime]] = mapped_column(TIMESTAMP(timezone=True))
    deleted_at: Mapped[Optional[datetime.datetime]] = mapped_column(TIMESTAMP(timezone=True))
    deck_id: Mapped[Optional[int]] = mapped_column(BigInteger)

    def __repr__(self) -> str:
        return f"<DigimonCards(id={self.id}, name={self.name})>"
    



