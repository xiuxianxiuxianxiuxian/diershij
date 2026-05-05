"""Multi-cultivate + breakthrough test"""
import asyncio, json, websockets, requests, time, sys

BASE = "http://localhost:8081"
WS = "ws://localhost:8081"

async def main():
    ts = int(time.time())
    username = f"bt_test_{ts}"

    # ── Register ──
    print(f"注册用户: {username}")
    resp = requests.post(f"{BASE}/auth/register", json={"username": username, "password": "123456"})
    data = resp.json()
    assert data.get("success"), f"注册失败: {data}"
    token = data["token"]
    entity_id = data["entity"]["id"]
    print(f"  EntityID: {entity_id[:16]}...")

    # ── Connect WebSocket ──
    uri = f"{WS}/ws?token={token}"
    async with websockets.connect(uri) as ws:
        first = json.loads(await asyncio.wait_for(ws.recv(), timeout=5))
        print(f"  WS连接成功")

        def send_op(action_type, params=None):
            return json.dumps({"type": "operation", "payload": {"action_type": action_type, "params": params or {}}})

        async def recv(label, timeout=10):
            resp = json.loads(await asyncio.wait_for(ws.recv(), timeout=timeout))
            while resp.get("type") == "state_sync":
                resp = json.loads(await asyncio.wait_for(ws.recv(), timeout=timeout))
            p = resp.get("payload", {})
            return p

        # ── Learn a method & set main ──
        print(f"\n1. 学习功法")
        await ws.send(send_op("list_methods"))
        r = await recv("list_methods")
        methods = r.get("effects", {}).get("methods", [])
        if methods:
            mid = methods[0]["id"]
            mname = methods[0]["name"]
            await ws.send(send_op("learn_method", {"method_id": mid}))
            r = await recv("learn_method")
            print(f"   学习: {mname} -> {r.get('message')}")

            await ws.send(send_op("set_main_method", {"method_id": mid}))
            r = await recv("set_main_method")
            print(f"   设置主修: {r.get('message')}")

        # ── Multiple cultivations ──
        print(f"\n2. 多次修炼 (积累修为到 100 进行突破)")
        total_progress = 0
        rounds = 0
        while total_progress < 100 and rounds < 200:
            await ws.send(send_op("cultivate"))
            r = await recv("cultivate")
            gain = r.get("effects", {}).get("cultivation_gain", 0)
            total_progress = r.get("effects", {}).get("progress", 0)
            rounds += 1
            if rounds <= 5 or rounds % 20 == 0 or total_progress >= 100:
                print(f"   第{rounds:2d}次修炼: +{gain:.4f} 修为, 当前进度={total_progress:.2f}")

        print(f"\n   {rounds}次修炼后, 修为进度={total_progress:.2f}")

        # ── Breakthrough ──
        print(f"\n3. 尝试突破!")
        await ws.send(send_op("breakthrough"))
        r = await recv("breakthrough")
        success = r.get("success", False)
        msg = r.get("message", "")
        effects = r.get("effects", {})

        print(f"   结果: {'成功!' if success else '失败'}, 消息: {msg}")
        if effects:
            for k, v in effects.items():
                print(f"   {k}: {v}")

        # ── If breakthrough succeeded, check new realm ──
        if success:
            print(f"\n4. 突破后状态:")
            resp = requests.post(f"{BASE}/auth/login", json={"username": username, "password": "123456"})
            data = resp.json()
            entity = data.get("entity", {})
            new_realm = entity.get("realm", "?")
            new_attrs = entity.get("attributes", {})
            print(f"   新境界: {new_realm}")
            print(f"   新MaxQi: {new_attrs.get('max_qi', '?')}")
            print(f"   新MaxSpiritualPower: {new_attrs.get('max_spiritual_power', '?')}")
            print(f"   新寿命: {new_attrs.get('max_lifespan', '?')}")

        print(f"\n测试完成")

if __name__ == "__main__":
    asyncio.run(main())
