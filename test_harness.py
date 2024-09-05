"""
Simple tests for python harness

Usage:
$ pytest test_harness.py

"""
from crawleruseragents import is_crawler, matching_crawlers


def test_match():
    assert is_crawler("test Googlebot/2.0 test") is True


def test_nomatch():
    assert is_crawler("!!!!!!!!!!!!") is False


def test_case():
    assert is_crawler("test googlebot/2.0 test") is False


def test_matching_crawlers_match():
    result = matching_crawlers("test Googlebot/2.0 test")
    assert isinstance(result, list)
    assert len(result) > 0
    assert all(isinstance(val, int) for val in result)


def test_matching_crawlers_nomatch():
    result = matching_crawlers("!!!!!!!!!!!!")
    assert isinstance(result, list)
    assert len(result) == 0


def test_matching_crawlers_case():
    result = matching_crawlers("test googlebot/2.0 test")
    assert isinstance(result, list)
    assert len(result) == 0
