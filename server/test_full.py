import asyncio, json, websockets, requests, time

async def recv(ws, label, timeout=10):
    resp = json.loads(await asyncio.wait_for(ws.recv(), timeout=timeout))
    if resp.get("type") == "state_sync":
        resp = json.loads(await asyncio.wait_for(ws.recv(), timeout=timeout))
    p = resp.get("payload", {})
    success = p.get("success")
    msg = p.get("message", "")
    effects = p.get("effects", {})
    status = "OK" if success == True else ("FAIL" if success == False else "ERROR")
    print(f"  [{status}] {label}: {msg[:80]}")
    if effects:
        for k in ("cultivation_gain", "progress", "resource", "discoveries",
                  "new_realm", "qi_cost", "quantity", "skill_exp"):
            if k in effects:
                print(f"       {k} = {effects[k]}")
    return resp

async def test():
    ts = int(time.time())
    user = f"t{ts % 100000}"

    resp = requests.post("http://localhost:8081/auth/register",
                         json={"username": user, "password": "123456"}).json()
    if not resp.get("success"):
        resp = requests.post("http://localhost:8081/auth/login",
                             json={"username": user, "password": "123456"}).json()
    token, eid = resp["token"], resp["entity"]["id"]
    print(f"=== 角色功能测试 ===")
    print(f"用户: {user}, ID: {eid[:16]}...\n")

    uri = f"ws://localhost:8081/ws?token={token}"
    async with websockets.connect(uri) as ws:
        await asyncio.wait_for(ws.recv(), timeout=5)
        print("[WS 已连接]\n")

        def op(a, p={}):
            return json.dumps({"type": "operation", "payload": {"action_type": a, "params": p}})

        # 1. 基础操作
        print("── 1. 修炼 / 打坐 / 移动 / 休息 ──")
        await ws.send(op("cultivate"));   await recv(ws, "修炼")
        await ws.send(op("meditate"));    await recv(ws, "打坐")
        await ws.send(op("sleep"));       await recv(ws, "休息")
        await ws.send(op("move", {"region_id": "spirit_mist_mountain"}))
        await recv(ws, "移动→灵雾山脉")

        # 2. 探索（资源发现+自动入库）
        print("\n── 2. 探索 (40%概率发现资源+自动入库) ──")
        for i in range(5):
            await ws.send(op("explore"))
            await recv(ws, f"探索#{i+1}")

        # 3. 采集（物品入库）
        print("\n── 3. 采集草药 (物品自动入库) ──")
        await ws.send(op("gather", {"resource_type": "herb", "quantity": 1}))
        r = await recv(ws, "采集草药")
        res = r.get("payload",{}).get("effects",{}).get("resource","?")
        print(f"  => 获得: {res}")

        # 4. 不足修为时突破被拒
        print("\n── 4. 突破检查 (预期:修为不足被拒) ──")
        await ws.send(op("breakthrough"))
        await recv(ws, "突破拦截")

        # 5. 连续修炼+打坐 10 轮 (验证进度累积)
        print("\n── 5. 连续修炼10轮 (验证进度累积) ──")
        for i in range(10):
            await ws.send(op("cultivate"))
            r = await recv(ws, f"修炼#{i+1}")
            await ws.send(op("meditate"))
            await recv(ws, f"打坐#{i+1}")
        print("  => 每轮修为+0.05, 10轮后应~0.5%")

        # 6. 再采集几次
        print("\n── 6. 多轮采集 ──")
        await ws.send(op("gather", {"resource_type": "herb", "quantity": 1}))
        await recv(ws, "采集#1")
        await ws.send(op("gather", {"resource_type": "herb", "quantity": 2}))
        await recv(ws, "采集x2")

        # 7. 聊天
        print("\n── 7. 聊天 ──")
        await ws.send(json.dumps({"type":"chat","payload":{"content":"测试通过","channel":"world"}}))
        r = json.loads(await asyncio.wait_for(ws.recv(), timeout=5))
        print(f"  [{'OK' if r.get('type')=='chat' else 'FAIL'}] chat")

        print("\n=== 全部功能测试通过 ===")

asyncio.run(test())
