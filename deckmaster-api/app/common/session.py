from collections.abc import AsyncGenerator
from contextlib import asynccontextmanager
from sqlite3 import IntegrityError

from fastapi.params import Depends
from sqlalchemy.ext.asyncio import AsyncSession as SqlAlchemyAsyncSession
from sqlalchemy.ext.asyncio import async_sessionmaker, create_async_engine
from typing_extensions import Annotated

from common.database import get_db_string
from common.errors import integrity_error
from common.exceptions.exceptions import ConflictException

engine = create_async_engine(get_db_string())

session_factory = async_sessionmaker(bind=engine, class_=SqlAlchemyAsyncSession, autocommit=False, autoflush=False, expire_on_commit=False)

@asynccontextmanager
async def session_scope() -> AsyncGenerator[SqlAlchemyAsyncSession]:
    db : SqlAlchemyAsyncSession = session_factory()
    try:
        yield db
        await db.commit()
    except Exception:
        await db.rollback()
        raise
    finally:
        await db.close()


async def get_db() -> AsyncGenerator[SqlAlchemyAsyncSession]:
    async with session_scope() as db:
        try:
            yield db
        except IntegrityError as e:
            raise ConflictException(integrity_error()) from e
        

AsyncSession = Annotated[SqlAlchemyAsyncSession, Depends(get_db)]