import requests
import json

# 요청할 URL
url = 'http://127.0.0.1/postInstruction'

# 전송할 데이터
data = {
    "agentUUID": "12342",
    "procedureID": "P_Collection_Kimsuky_001",
    "messageUUID": "f5556669-ffbe-4d24-b833-fc9888fdeaef"
}

# 헤더 설정
headers = {
    'Content-Type': 'application/json'
}

# POST 요청 전송
response = requests.post(url, headers=headers, data=json.dumps(data))

# 응답 출력
if response.status_code == 200:
    print("Request succeeded:", response.json())
else:
    print(f"Request failed with status code {response.status_code}: {response.text}")
