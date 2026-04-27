import os

from sqlalchemy.orm import DeclarativeBase


class Base(DeclarativeBase):
    pass




def get_db_string() -> str:
    host = os.getenv("DB_HOST")
    port = os.getenv("DB_PORT")
    user = os.getenv("DB_USER")
    password = os.getenv("DB_PASSWORD")
    name = os.getenv("DB_NAME")
    return f"postgresql+asyncpg://{user}:{password}@{host}:{port}/{name}"

def get_test_db_string() -> str:
    host = os.getenv("TEST_DB_HOST")
    port = os.getenv("TEST_DB_PORT")
    password = os.getenv("TEST_DB_PASSWORD")
    user = os.getenv("TEST_DB_USER")
    name = os.getenv("TEST_DB_NAME")
    return f"postgresql+asyncpg://{user}:{password}@{host}:{port}/{name}"