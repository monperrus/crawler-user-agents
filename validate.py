import json
import re
from  jsonschema import validate

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

# check for duplicates
class Entry: # we need this to be able to use i!=j with exact same content
    pass
data2 = []
for i in data:
    x=Entry()
    x.pattern = i['pattern'].lower()
    data2.append(x)
for i in data2:
    for j in data2:
        if i!=j and re.match(i.pattern, j.pattern, re.IGNORECASE):
            raise Exception('duplicate '+i.pattern+' '+j.pattern)

