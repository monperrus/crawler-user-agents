"""
Validate JSON to ensure that patterns all work
"""
from __future__ import print_function

import json
import re
from collections import Counter

from jsonschema import validate


JSON_SCHEMA = {
    "type": "array",
    "items": {
        "type": "object",
        "properties": {
            "pattern": {"type": "string"},
            "instances": {"type": "array"},
            "url": {"type": "string"},
            "addition_date": {"type": "string"}
        },
        "required": ["pattern"]
    }
}


def main():
    with open('crawler-user-agents.json') as f:
        json_data = json.load(f)

    # check format using JSON Schema
    validate(json_data, JSON_SCHEMA)

    # check for simple duplicates
    pattern_counts = Counter(entry['pattern'] for entry in json_data)
    for pattern, count in pattern_counts.most_common():
        if count > 1:
            raise ValueError('Pattern {!r} appears {} times'.format(pattern,
                                                                    count))

    # check for duplicates with different capitalization
    patterns = sorted(entry['pattern'].lower() for entry in json_data)
    last = ""
    for pattern in patterns:
        if pattern == last:
            raise ValueError('Pattern {!r} is duplicated with different capitalization'.format(pattern))
        last = pattern

    # check that we match the given instances
    num_instances = 0
    for entry in json_data:
        pattern = entry['pattern']
        instances = entry.get('instances')
        if instances:
            for instance in instances:
                num_instances += 1
                if not re.search(pattern, instance):
                    raise ValueError('Pattern {!r} misses instance {!r}'
                                     .format(pattern, instance))
                # TODO: Check for re2 matching here

    # Make sure we have at least 10 instances in file
    if num_instances < 10:
        raise ValueError('Only had {} instances in JSON'.format(num_instances))

    # Check for patterns that match other patterns
    for entry1 in json_data:
        for entry2 in json_data:
            if entry1 != entry2 and re.search(entry1['pattern'],
                                              entry2['pattern']):
                raise ValueError('Pattern {!r} is a subset of {!r}'
                                 .format(entry2['pattern'], entry1['pattern']))

    print('Validation passed')


if __name__ == '__main__':
    main()
