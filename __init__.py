import re
import json
from functools import cached_property
from pathlib import Path


class CrawlerPatterns:
    def __init__(self):
        pass

    @cached_property
    def case_insensitive(self):
        return re.compile(
            "|".join(i["pattern"] for i in CRAWLER_USER_AGENTS_DATA),
            re.IGNORECASE
        )

    @cached_property
    def case_sensitive(self):
        return re.compile("|".join(i["pattern"] for i in CRAWLER_USER_AGENTS_DATA))


def load_json():
    cwd = Path(__file__).parent
    user_agents_file_path = cwd / "crawler-user-agents.json"
    with user_agents_file_path.open() as patterns_file:
        return json.load(patterns_file)


CRAWLER_USER_AGENTS_DATA = load_json()
CRAWLER_PATTERNS = CrawlerPatterns()


def is_crawler(user_agent: str, case_sensitive: bool = True) -> bool:
    """Return True if the given User-Agent matches a known crawler."""
    if case_sensitive:
        return bool(re.search(CRAWLER_PATTERNS.case_sensitive, user_agent))
    return bool(re.search(CRAWLER_PATTERNS.case_insensitive, user_agent))


def matching_crawlers(user_agent: str, case_sensitive: bool = True) -> list[int]:
    """
    Return a list of the indices in CRAWLER_USER_AGENTS_DATA of any crawlers
    matching the given User-Agent.
    """
    result = []
    if is_crawler(user_agent, case_sensitive):
        for num, crawler_user_agent in enumerate(CRAWLER_USER_AGENTS_DATA):
            if re.search(crawler_user_agent["pattern"], user_agent, 0 if case_sensitive else re.IGNORECASE):
                result.append(num)
    return result
