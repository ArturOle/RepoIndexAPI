from pydantic import BaseModel


class Repo(BaseModel):
    name: str
    owner: str
    url: str
