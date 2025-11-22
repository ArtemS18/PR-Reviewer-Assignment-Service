import os
from pathlib import Path
from dotenv import load_dotenv

BASE_DIR = Path(__file__).resolve().parent.parent 
ENV_PATH = BASE_DIR.parent / ".env"

load_dotenv(ENV_PATH, encoding= "utf-8") 


TEAMS_COUNT = 20
USERS_PER_TEAM = 10  # 20 * 10 = 200
SEED_DATA_FILE = "seed_data.json"

TEST_TIME = "1m"

HTTP_PORT = os.getenv("HTTP_PORT", "8080").strip()
BASE_URL =f"http://localhost:{HTTP_PORT}"

DB_HOST = os.getenv("DB_HOST", "localhost").strip()
DB_PORT = os.getenv("DB_OPEN_PORT", "5432").strip()
DB_NAME = os.getenv("DB_NAME", "postgres").strip()
DB_USER = os.getenv("DB_USER", "postgres").strip()
DB_PASS = os.getenv("DB_PASS", "postgres").strip()