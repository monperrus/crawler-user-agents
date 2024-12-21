"""
Simple tests for python harness

Usage:
$ pytest test_harness.py

"""

from crawleruseragents import is_crawler, matching_crawlers



def test_nomatch():
    assert is_crawler("!!!!!!!!!!!!") is False


def test_case_sensitive():
    assert is_crawler("test googlebot/2.0 test") is False


def test_case_insensitive():
    assert is_crawler("test googlebot/2.0 test", case_sensitive=False) is True


def test_matching_crawlers_match_case_sensitive():
    result = matching_crawlers("test Googlebot/2.0 test")
    assert isinstance(result, list)
    assert len(result) > 0
    assert all(isinstance(val, int) for val in result)


def test_matching_crawlers_match_case_insensitive():
    result = matching_crawlers("test googlebot/2.0 test", False)
    assert isinstance(result, list)
    assert len(result) > 0
    assert all(isinstance(val, int) for val in result)

def test_matching_crawlers_match_lower_case_agent():
    result = matching_crawlers("test googlebot/2.0 test")
    assert isinstance(result, list)
    assert len(result) == 0


def test_matching_crawlers_nomatch():
    result = matching_crawlers("!!!!!!!!!!!!")
    assert isinstance(result, list)
    assert len(result) == 0


def test_matching_crawlers_case():
    result = matching_crawlers("test googlebot/2.0 test")
    assert isinstance(result, list)
    assert len(result) == 0
