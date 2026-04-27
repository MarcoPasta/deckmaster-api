from common.exceptions.error import DeckmasterError


def integrity_error() -> DeckmasterError:
    return DeckmasterError(
        code="INTEGRITY_ERROR",
        message="A database integrity error occurred. Resource already exists or violates constraints.",
    )
