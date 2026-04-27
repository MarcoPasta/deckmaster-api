from dataclasses import asdict

from fastapi import HTTPException

from common.exceptions.error import DeckmasterError


class DeckmasterException(HTTPException):
    status_code: int = 500
    def __init__(self, error: DeckmasterError):
        detail = [{"error": asdict(error)}]
        super().__init__(status_code=self.status_code, detail=detail)

class ConflictException(DeckmasterException):
    status_code = 409