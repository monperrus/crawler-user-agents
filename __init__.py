import re
import json
from pathlib import Path


def load_json():
    cwd = Path(__file__).parent
    user_agents_file_path = cwd / "crawler-user-agents.json"
    with user_agents_file_path.open() as patterns_file:
        return json.load(patterns_file)


CRAWLER_USER_AGENTS_DATA = load_json()


def is_crawler(user_agent: str) -> bool:
    for crawler_user_agent in CRAWLER_USER_AGENTS_DATA:
        if re.search(crawler_user_agent["pattern"], user_agent, re.IGNORECASE):
            return True
    return False


def is_crawler2(s):
    regexp = re.compile("|".join([i["pattern"] for i in CRAWLER_USER_AGENTS_DATA]))
    return regexp.search(s) is not None
