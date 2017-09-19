"""
Break crawler-user-agents.json in ways that validate.py should detect
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
def user_agent_list():
    # Load clean version of crawler-user-agents.json
    with open('crawler-user-agents.json') as f:
        original_json = json.load(f)
    # Yield copy of it for test function to break and overwrite file with
    yield deepcopy(original_json)
    # Overwrite broken version with clean version after test completes
    update_json_file(original_json)


def assert_validate_failed():
    assert subprocess.call(['python', 'validate.py']) != 0


def assert_validate_passed():
    assert subprocess.call(['python', 'validate.py']) == 0


def test_simple_pass(user_agent_list):
    # Check that a single pattern with more than 10 instances will pass
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
                                      'foot',]}]
    update_json_file(user_agent_list)
    assert_validate_passed()


def test_schema_violation_dict(user_agent_list):
    user_agent_list = {'foo': user_agent_list}
    update_json_file(user_agent_list)
    assert_validate_failed()


def test_schema_violation_int(user_agent_list):
    user_agent_list[0]['pattern'] = 2
    update_json_file(user_agent_list)
    assert_validate_failed()


def test_simple_duplicate_detection(user_agent_list):
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
                                      'foot',]},
                        {'pattern': 'foo'}]
    update_json_file(user_agent_list)
    assert_validate_failed()


def test_subset_duplicate_detection(user_agent_list):
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
                                      'foot',]},
                        {'pattern': 'afoot'}]
    update_json_file(user_agent_list)
    assert_validate_failed()
