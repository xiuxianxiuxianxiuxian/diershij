import asyncio, json, websockets, requests

async def test():
    resp = requests.post("http://localhost:8081/auth/register", json={"username": f"tester_{int(asyncio.get_event_loop().time())}", "password": "123456"})
    if not resp.json().get("success"):
        resp = requests.post("http://localhost:8081/auth/login", json={"username": "tester2", "password": "123456"})
    data = resp.json()
    token = data["token"]
    entity_id = data["entity"]["id"]
    print(f"[OK] Login: entity={entity_id[:12]}")
    if not token:
        return

    uri = f"ws://localhost:8081/ws?token={token}"
    async with websockets.connect(uri) as ws:
        # consume state_sync
        first = json.loads(await asyncio.wait_for(ws.recv(), timeout=5))
        print(f"[OK] WS connected, state_sync received")

        def send_op(action_type, params={}):
            return json.dumps({"type": "operation", "payload": {"action_type": action_type, "params": params}})

        async def recv_result(label):
            resp = json.loads(await asyncio.wait_for(ws.recv(), timeout=10))
            # If we get state_sync, it's a push - consume next message
            if resp.get("type") == "state_sync":
                resp = json.loads(await asyncio.wait_for(ws.recv(), timeout=10))
            p = resp.get("payload", {})
            success = p.get("success")
            msg = p.get("message", "")
            status = "OK" if success == True else ("FAIL" if success == False else "ERROR")
            print(f"  [{status}] {label}: {msg[:60]}")
            return resp

        # Test operations with valid regions
        ops = [
            ("cultivate", {}),
            ("move", {"region_id": "qingyun_town"}),
            ("move", {"region_id": "spirit_mist_mountain"}),
            ("meditate", {}),
            ("sleep", {}),
            ("explore", {}),
            ("gather", {"resource_type": "herb", "quantity": 1}),
        ]

        for action_type, params in ops:
            await ws.send(send_op(action_type, params))
            await recv_result(action_type)

        # Test breakthrough
        await ws.send(send_op("breakthrough", {}))
        await recv_result("breakthrough")

        # Test chat
        await ws.send(json.dumps({"type": "chat", "payload": {"content": "大家好!", "channel": "world"}}))
        resp = json.loads(await asyncio.wait_for(ws.recv(), timeout=5))
        ct = resp.get("type")
        print(f"  [{'OK' if ct=='chat' else 'FAIL'}] chat: type={ct}")

        print("\nAll tests done.")

asyncio.run(test())
