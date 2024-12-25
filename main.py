import os

app_path = "src.saucer_api.api:app"

os.system(f'uvicorn {app_path} --reload --host=0.0.0.0 --port=8000')
