import json
import re2 as re
from jsonschema import validate

re.set_fallback_notification(re.FALLBACK_EXCEPTION)

data = json.load(open('crawler-user-agents.json'))

# check for format using JSON Schema
schema = {
    "type": "array",
    "items": {
        "type": "object",
        "properties" :  {
            "pattern" : {"type" : "string"},
        },
        "required" : [ "pattern" ]
    }
}
validate(data, schema)

data = [i['pattern'].lower() for i in data]

for i in data:
    for j in data:
        # check re2 library (https://github.com/google/re2) compatiblity
        try:
            match = re.match(i, j, re.IGNORECASE)
        except re.RegexError:
            raise Exception('regex "{}" is not compatible with re2 library'.format(i))

        # check for duplicates
        if i != j and match:
            raise Exception('duplicate found "{}" and "{}"'.format(i, j))

