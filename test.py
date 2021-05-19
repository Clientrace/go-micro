
import requests


resp = requests.post(
    'https://dtfixrp294.execute-api.ap-southeast-1.amazonaws.com/dev/user',
    headers={
        'Content-Type': 'application/json'
    },
    json={
        'testField': 'testVal'
    }
)

print(resp.json())




