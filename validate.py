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
    x.instances = i['instances'] if "instances" in i else []
    data2.append(x)

nbinstances = 0
for i in data2:
    # contract #1 do we match the given instances?
    for instance in i.instances:
        nbinstances += 1
        assert re.search(i.pattern, instance, re.IGNORECASE), i.pattern+" does not match "+instance
            
        # ongoing work by @vetty
        # assert re2.match(i.pattern, instance, re.IGNORECASE)
        
    for j in data2:
        if i!=j and re.search(i.pattern, j.pattern, re.IGNORECASE):
            raise Exception('duplicate '+i.pattern+' '+j.pattern)

assert nbinstances>10

