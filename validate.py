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
            "pattern": {"type": "string"}, # required
            "instances": {"type": "array"}, # required
            "url": {"type": "string"}, # optional
            "description": {"type": "string"}, # optional
            "addition_date": {"type": "string"}, # optional
            "depends_on": {"type": "array"} # allows an instance to match twice
        },
        "required": ["pattern", "instances"]
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
    pattern_counts = Counter(entry['pattern'].lower() for entry in json_data)
    for pattern, count in pattern_counts.most_common():
        if count > 1:
            raise ValueError('Pattern {!r} is duplicated {} times with different capitalization'
                             .format(pattern, count))
    
    # checks that no pattern contains unescaped slash / 
    for entry in json_data:
        pattern = entry['pattern']
        if re.search('[^\\\\]/', pattern):
            raise ValueError('Pattern {!r} has an unescaped slash character'.format(pattern))

    # check that we match the given instances
    num_instances = 0
    for entry in json_data:
        pattern = entry['pattern']
        
        # canonicalize entry
        if 'depends_on' not in entry: entry['depends_on'] = []
            
        # check that we have only the rights properties (not handled by default in module jsonschema)
        assert set([str(x) for x in entry.keys()]).issubset(set(JSON_SCHEMA['items']['properties'].keys())), "the entry contains unknown properties"  
        instances = entry.get('instances')
        if instances:
            # check that there is no duplicate
            if not len(instances) == len(set(instances)):
                raise Exception("duplicate instances in "+pattern)
            for instance in instances:
                num_instances += 1
                if not re.search(pattern, instance):
                    raise ValueError('Pattern {!r} misses instance {!r}'
                                     .format(pattern, instance))
                

    # Make sure we have at least one pattern
    if len(json_data) < 1:
        raise Exception("no pattern")

    # Check for patterns that match other patterns
    for entry1 in json_data:
        for entry2 in json_data:
            if entry1 != entry2 and re.search(entry1['pattern'],
                                              entry2['pattern'],re.IGNORECASE):
                raise ValueError('Pattern {!r} is a subset of {!r}'
                                 .format(entry2['pattern'], entry1['pattern']))

    print('Validation passed')


if __name__ == '__main__':
    main()
