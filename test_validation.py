"""
Break crawler-user-agents.json in ways that validate.py should detect

Usage:
$ pytest test_validation.py

"""
from __future__ import print_function

import json
import subprocess
from copy import deepcopy

import pytest


def update_json_file(data):
    with open('crawler-user-agents.json', 'w') as f:
        json.dump(data, f, indent=2)

@pytest.fixture
def restore_original_json():
    # Load original version of crawler-user-agents.json
    with open('crawler-user-agents.json') as f:
        original_json = json.load(f)

    # By using a yield statement instead of return, all the code after
    # the yield statement serves as the teardown code:
    yield None

    # tear down code: restore original version of crawler-user-agents.json
    update_json_file(original_json)


def assert_validate_failed():
    assert subprocess.call(['python', 'validate.py']) != 0


def assert_validate_passed():
    assert subprocess.call(['python', 'validate.py']) == 0


def test_simple_pass(restore_original_json):
    # the json must be an array of objects containing "pattern"
    # there must be more than 10 instances to pass
    user_agent_list = [{'pattern': 'foo',
                        'instances': ['foo',
                                      'afoo',
                                      'foob',
                                      'cfood',
                                      '/foo/',
                                      '\\foo\\',
                                      ':foo:',
                                      'foo.',
                                      '!foo',
                                      '/foo',
                                      'foo\\',
                                      'FoofooFoo',
                                      'foot']}]
    update_json_file(user_agent_list)
    assert_validate_passed()

def test_simplest_pass(restore_original_json):
    # the simplest crawler file passes
    user_agent_list = [{'pattern': 'foo', 'instances': []}]
    update_json_file(user_agent_list)
    assert_validate_passed()

def test_schema_violation_dict1(restore_original_json):
    # contract: the json must be an array
    user_agent_list = {'foo':None}
    update_json_file(user_agent_list)
    assert_validate_failed()

def test_schema_violation_dict2(restore_original_json):
    # contract: the json must be an array of objects containing "pattern"
    user_agent_list = [{'foo':None}]
    update_json_file(user_agent_list)
    assert_validate_failed()

def test_schema_violation_dict3(restore_original_json):
    # contract: the json must be an array of objects containing "pattern" and valid properties
    user_agent_list = [{'pattern':'foo', 'foo':3}]
    update_json_file(user_agent_list)
    assert_validate_failed()

def test_simple_duplicate_detection(restore_original_json):
    # contract: if we have the same pattern twice, it fails
    user_agent_list = [{'pattern': 'foo',
                        'instances': ['foo',
                                      'afoo',
                                      'foob',
                                      'cfood',
                                      '/foo/',
                                      '\\foo\\',
                                      ':foO:',
                                      'foo.',
                                      '!foo',
                                      '/foo',
                                      'foo\\',
                                      'FoofooFoo',
                                      'foot']},
                        {'pattern': 'foo'}]
    update_json_file(user_agent_list)
    assert_validate_failed()


def test_simple_duplicate_detection2(restore_original_json):
    # contract: if we have the same pattern twice, it fails (even w/o instances)
    user_agent_list = [{'pattern': 'foo',
                        'instances': ['foo',
                                      'afoo',
                                      'foob',
                                      'cfood',
                                      '/foo/',
                                      '\\foo\\',
                                      ':foo:',
                                      'foo.',
                                      '!foo',
                                      '/foo',
                                      'foo\\',
                                      'FoofooFoo',
                                      'foot']},
                       {'pattern': 'bar'},
                       {'pattern': 'bar'}]
    update_json_file(user_agent_list)
    assert_validate_failed()


def test_subset_duplicate_detection(restore_original_json):
    # contract: if a pattern matches another pattern, it fails
    user_agent_list = [{'pattern': 'foo',
                        'instances': ['foo',
                                      'afoo',
                                      'foob',
                                      'cfood',
                                      '/foo/',
                                      '\\foo\\',
                                      ':foo:',
                                      'foo.',
                                      '!foo',
                                      '/foo',
                                      'foo\\',
                                      'FoofooFoo',
                                      'foot']},
                       {'pattern': 'afoot'}]
    update_json_file(user_agent_list)
    assert_validate_failed()


def test_case_sensitivity(restore_original_json):
    # contract: the patterns are case sensitive
    user_agent_list = [{'pattern': 'foo',
                        'instances': ['FOO',
                                      'aFoo',
                                      'foob',
                                      'cfood',
                                      '/foo/',
                                      '\\foo\\',
                                      ':foo:',
                                      'fOo.',
                                      '!FOO',
                                      '/foo',
                                      'foo\\',
                                      'FoofooFoo',
                                      'foot']}]
    update_json_file(user_agent_list)
    assert_validate_failed()

def test_duplicate_case_insensitive_detection(restore_original_json):
    # contract: fail if we have patterns that differ only in capitailization
    user_agent_list = [{'pattern': 'foo',
                        'instances': ['foo',
                                      'afoo',
                                      'foob',
                                      'cfood',
                                      '/foo/',
                                      '\\foo\\',
                                      ':foo:',
                                      'foo.',
                                      '!foo',
                                      '/foo',
                                      'foo\\',
                                      'FoofooFoo',
                                      'foot']},
                       {'pattern': 'fOo'}]
    update_json_file(user_agent_list)
    assert_validate_failed()
