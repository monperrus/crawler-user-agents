import crawleruseragents
import re
import json
from importlib import resources 
    
def load_json():
  return json.loads(resources.read_text(crawleruseragents,"crawler-user-agents.json"))

DATA = load_json()

def is_crawler(s):
  # print(s)
  for i in DATA:
    test=re.search(i["pattern"],s,re.IGNORECASE)
    if test:
      return True
  return False

def is_crawler2(s):
  regexp = re.compile("|".join([i["pattern"] for i in DATA]))
  return regexp.search(s) != None

