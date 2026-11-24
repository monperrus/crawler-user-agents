"""
Simple tests for python harness

Usage:
$ pytest test_harness.py

"""

from crawleruseragents import is_crawler, matching_crawlers



def test_nomatch():
    assert is_crawler("!!!!!!!!!!!!") is False


def test_android_instagram_build_fingerprint_is_not_crawler():
    # AP3A.240617.008 is an Android 15 platform build id on ordinary Instagram traffic.
    user_agent = "Mozilla/5.0 (Linux; Android 15; CPH2557 Build/AP3A.240617.008; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/142.0.7444.142 Mobile Safari/537.36 Instagram 406.0.0.58.159 Android (35/15; 480dpi; 1080x2400; OPPO; CPH2557; OP573DL1; mt6833; en_MY; 822918295; IABMV/1) NV/1"
    assert is_crawler(user_agent) is False


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
