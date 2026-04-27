from typing import Any, Dict

from common.session import get_db
from fastapi import Depends, FastAPI
from modules.digimon_cards.models import DigimonCards
from sqlalchemy import select

app = FastAPI()



@app.get("/")
def read_root() -> Dict[str, str]:
    return {"Hello": "World"}


@app.get("/items/{item_id}")
def read_item(item_id: int, q: str = "") -> Dict[str, int | str]:
    return {"item_id": item_id, "q": q}

@app.get("/test")
async def read_test(db: Any = Depends(get_db)):
    statement = select(DigimonCards)
    result = await db.execute(statement)
    cards = result.scalars().all()
    return {"cards": cards}