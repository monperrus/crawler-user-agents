import json
from  jsonschema import validate

data = json.load(open('crawler-user-agents.json'))

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
