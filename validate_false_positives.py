"""
Validate that our crawler patterns do not produce too many false positives
on real production traffic datasets.

For each dataset, compute the weighted false positive rate (fraction of
non-bot requests that our patterns incorrectly flag as bots). Fail if any
dataset exceeds the threshold.
"""
from __future__ import annotations

import json
import sys
import urllib.request

from crawleruseragents import is_crawler

# Maximum allowed false positive rate per dataset (as a percentage)
FALSE_POSITIVE_THRESHOLD_PCT = 1.0

DATASETS = [
    {
        "name": "top_user_agents (mwnciau/user_agent_dumps)",
        "url": (
            "https://raw.githubusercontent.com/mwnciau/user_agent_dumps"
            "/refs/heads/main/top_user_agents/user_agents_parsed.json"
        ),
        # each entry: {"userAgent": "...", "count": N, "isBot": bool, ...}
        "user_agents_key": "userAgents",
        "ua_field": "userAgent",
        "count_field": "count",
        "is_bot_field": "isBot",
    },
]


def check_dataset(dataset: dict) -> tuple[float, list[str]]:
    """Return (false_positive_rate_pct, list_of_false_positive_ua_strings)."""
    with urllib.request.urlopen(dataset["url"]) as resp:
        data = json.load(resp)

    entries = data[dataset["user_agents_key"]]
    ua_field = dataset["ua_field"]
    count_field = dataset["count_field"]
    is_bot_field = dataset["is_bot_field"]

    total_requests = 0
    fp_requests = 0
    fp_uas: list[str] = []

    for entry in entries:
        if entry.get(is_bot_field):
            continue
        ua = entry[ua_field]
        count = entry.get(count_field, 1)
        total_requests += count
        if is_crawler(ua):
            fp_requests += count
            fp_uas.append(ua)

    rate = (fp_requests / total_requests * 100) if total_requests else 0.0
    return rate, fp_uas


def main() -> None:
    failed = False
    for dataset in DATASETS:
        print(f"Checking: {dataset['name']}")
        rate, fp_uas = check_dataset(dataset)
        print(f"  False positive rate: {rate:.4f}%  (threshold: {FALSE_POSITIVE_THRESHOLD_PCT}%)")
        if fp_uas:
            print("  False positives:")
            for ua in fp_uas:
                print(f"    {ua}")
        if rate > FALSE_POSITIVE_THRESHOLD_PCT:
            print(f"  FAIL: rate {rate:.4f}% exceeds threshold {FALSE_POSITIVE_THRESHOLD_PCT}%")
            failed = True
        else:
            print("  OK")

    if failed:
        sys.exit(1)
    print("False positive validation passed")


if __name__ == "__main__":
    main()
