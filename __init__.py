import re
import json
from pathlib import Path


def load_json():
    cwd = Path(__file__).parent
    user_agents_file_path = cwd / "crawler-user-agents.json"
    with user_agents_file_path.open() as patterns_file:
        return json.load(patterns_file)


CRAWLER_USER_AGENTS_DATA = load_json()
CRAWLER_USER_AGENTS_REGEXP = re.compile(
    "|".join(i["pattern"] for i in CRAWLER_USER_AGENTS_DATA)
)


def is_crawler(user_agent: str) -> bool:
    """Return True if the given User-Agent matches a known crawler."""
    return bool(CRAWLER_USER_AGENTS_REGEXP.search(user_agent))


def matching_crawlers(user_agent: str) -> list[int]:
    """
    Return a list of the indices in CRAWLER_USER_AGENTS_DATA of any crawlers
    matching the given User-Agent.
    """
    result = []
    if is_crawler(user_agent):
        for num, crawler_user_agent in enumerate(CRAWLER_USER_AGENTS_DATA):
            if re.search(crawler_user_agent["pattern"], user_agent):
                result.append(num)
    return result
