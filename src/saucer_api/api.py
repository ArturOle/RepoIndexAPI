import json
import os

from fastapi import FastAPI, HTTPException
from fastapi.responses import JSONResponse, RedirectResponse
from fastapi.middleware.cors import CORSMiddleware

from slowapi import Limiter, _rate_limit_exceeded_handler
from slowapi.util import get_remote_address
from slowapi.middleware import SlowAPIMiddleware
from slowapi.errors import RateLimitExceeded

import httpx

app = FastAPI()

limiter = Limiter(key_func=get_remote_address)
app.state.limiter = limiter
app.add_exception_handler(RateLimitExceeded, _rate_limit_exceeded_handler)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
app.add_middleware(SlowAPIMiddleware)


CLIENT_ID = os.getenv("GITHUB_CLIENT_ID")
CLIENT_SECRET = os.getenv("GITHUB_CLIENT_SECRET")
REDIRECT_URI = "http://localhost:8000/auth/github/callback"


@app.get("/auth/github/login")
def login_with_github():
    url = f"https://github.com/login/oauth/authorize?client_id={CLIENT_ID}&redirect_uri={REDIRECT_URI}&scope=read:user"
    return RedirectResponse(url)


@app.get("/auth/github/callback")
def github_callback(code: str):
    token_url = "https://github.com/login/oauth/access_token"
    headers = {"Accept": "application/json"}
    data = {
        "client_id": CLIENT_ID,
        "client_secret": CLIENT_SECRET,
        "code": code,
    }

    response = httpx.post(token_url, headers=headers, data=data)
    response_data = response.json()

    if "access_token" not in response_data:
        raise HTTPException(
            status_code=400,
            detail="Failed to retrieve access token"
        )

    access_token = response_data["access_token"]

    user_info = httpx.get(
        "https://api.github.com/user",
        headers={"Authorization": f"Bearer {access_token}"}
    ).json()

    return {"user_info": user_info, "access_token": access_token}


@app.get("/repos/random")
def get_random_repositories(quantity: int = 15):
    response = [
        {
            'title': f'Random repository {i}',
            'describtion': f'This is a random repository {i} doing random things',
        } for i in range(quantity)
    ]

    return JSONResponse(content=response)
