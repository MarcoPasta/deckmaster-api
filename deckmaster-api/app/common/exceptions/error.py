from dataclasses import dataclass


@dataclass(frozen=True, slots=True)
class DeckmasterError:
    code: str
    message: str
    params: dict[str, str] | None = None