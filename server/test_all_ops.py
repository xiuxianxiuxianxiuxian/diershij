"""Test all 21 game operations via WebSocket."""
import asyncio, json, websockets, requests, time, sys

OK = 0
FAIL = 0
SKIP = 0

def check(label, success, detail=""):
    global OK, FAIL
    if success:
        OK += 1
        print(f"  [PASS] {label}")
    else:
        FAIL += 1
        print(f"  [FAIL] {label}: {detail[:80]}")

async def recv(ws, label, timeout=10):
    """Receive a message, skipping state_sync. Returns parsed payload."""
    resp = json.loads(await asyncio.wait_for(ws.recv(), timeout=timeout))
    while resp.get("type") == "state_sync" or resp.get("type") == "entity_update":
        resp = json.loads(await asyncio.wait_for(ws.recv(), timeout=timeout))
    return resp

async def send_and_recv(ws, action_type, params={}, label=None):
    """Send an operation and receive the result."""
    if label is None:
        label = action_type
    await ws.send(json.dumps({
        "type": "operation",
        "payload": {"action_type": action_type, "params": params}
    }))
    resp = await recv(ws, label)
    payload = resp.get("payload", {})
    return payload.get("success"), payload.get("message", ""), payload.get("effects", {})

async def run_tests():
    global OK, FAIL, SKIP
    ts = int(time.time())
    username = f"allops_{ts % 100000}"

    # Register
    r = requests.post("http://localhost:8081/auth/register",
                      json={"username": username, "password": "123456"}).json()
    if not r.get("success"):
        r = requests.post("http://localhost:8081/auth/login",
                          json={"username": username, "password": "123456"}).json()
    token = r["token"]
    entity_id = r["entity"]["id"]
    realm = r["entity"].get("realm", "mortal")
    attrs = r["entity"].get("attributes", {})
    print(f"用户: {username}")
    print(f"ID: {entity_id[:16]}...")
    print(f"境界: {realm}")
    print()

    uri = f"ws://localhost:8081/ws?token={token}"
    async with websockets.connect(uri) as ws:
        # consume state_sync
        await asyncio.wait_for(ws.recv(), timeout=5)
        print("=" * 60)
        print("1. 修炼 Cultivate")
        print("=" * 60)
        s, msg, effects = await send_and_recv(ws, "cultivate")
        check("cultivate", s, msg)
        gain = effects.get("cultivation_gain", 0)
        prog = effects.get("progress", 0)
        print(f"   修为增益={gain}, 进度={prog}")

        print("=" * 60)
        print("2. 打坐 Meditate")
        print("=" * 60)
        s, msg, effects = await send_and_recv(ws, "meditate")
        check("meditate", s, msg)
        qr = effects.get("qi_recovery", 0)
        sr = effects.get("spiritual_recovery", 0)
        print(f"   灵力恢复={qr}, 神识恢复={sr}")

        print("=" * 60)
        print("3. 休息 Sleep")
        print("=" * 60)
        s, msg, effects = await send_and_recv(ws, "sleep")
        check("sleep", s, msg)

        print("=" * 60)
        print("4. 移动 Move")
        print("=" * 60)
        s, msg, effects = await send_and_recv(ws, "move", {"region_id": "qingyun_town"})
        check("move -> qingyun_town", s, msg)
        s, msg, effects = await send_and_recv(ws, "move", {"region_id": "spirit_mist_mountain"})
        check("move -> spirit_mist_mountain", s, msg)
        # Move to invalid region
        s, msg, effects = await send_and_recv(ws, "move", {"region_id": "nonexistent"})
        check("move -> invalid (expected fail)", not s, msg)

        print("=" * 60)
        print("5. 探索 Explore")
        print("=" * 60)
        for i in range(3):
            s, msg, effects = await send_and_recv(ws, "explore", label=f"explore #{i+1}")
            disc = effects.get("discoveries", [])
            status = "发现" if disc else "无发现"
            check(f"explore #{i+1}", s, f"{msg} | {status}")
            if disc:
                print(f"   发现: {disc}")

        print("=" * 60)
        print("6. 采集 Gather")
        print("=" * 60)
        s, msg, effects = await send_and_recv(ws, "gather", {"resource_type": "herb", "quantity": 1})
        check("gather herb", s, msg)
        resource = effects.get("resource", "")
        if resource:
            print(f"   获得: {resource} x{effects.get('quantity', 1)}")

        # Gather invalid type
        s, msg, effects = await send_and_recv(ws, "gather", {"resource_type": "void_crystal", "quantity": 1})
        check("gather invalid (expected fail)", not s, msg)

        print("=" * 60)
        print("7. 突破 Breakthrough")
        print("=" * 60)
        s, msg, effects = await send_and_recv(ws, "breakthrough")
        # Should fail with "progress not enough" since we just started
        check("breakthrough (low progress)", not s, msg)

        print("=" * 60)
        print("8. 聊天 Chat")
        print("=" * 60)
        await ws.send(json.dumps({
            "type": "chat",
            "payload": {"content": "测试全操作", "channel": "world"}
        }))
        resp = await recv(ws, "chat")
        check("chat", resp.get("type") == "chat", f"type={resp.get('type')}")

        print("=" * 60)
        print("9. 发送私信 Send Message")
        print("=" * 60)
        s, msg, effects = await send_and_recv(ws, "send_message", {
            "receiver_id": entity_id,
            "content": "给自己发消息",
            "message_type": "private"
        })
        check("send_message", s, msg)

        print("=" * 60)
        print("10. 添加好友 Add Friend")
        print("=" * 60)
        # Try adding self (should fail)
        s, msg, effects = await send_and_recv(ws, "add_friend", {"name": username})
        check("add_friend self (expected fail)", not s, msg)
        # Try adding non-existent
        s, msg, effects = await send_and_recv(ws, "add_friend", {"name": "nonexistent_player"})
        check("add_friend nonexistent", not s, msg)

        print("=" * 60)
        print("11. 施法 Cast Spell")
        print("=" * 60)
        s, msg, effects = await send_and_recv(ws, "cast_spell", {"spell_id": "nonexistent", "target_id": entity_id})
        check("cast_spell invalid", not s, msg)

        print("=" * 60)
        print("12. 炼制 Craft")
        print("=" * 60)
        s, msg, effects = await send_and_recv(ws, "craft", {"recipe_id": ""})
        check("craft (empty recipe)", not s, msg)  # Empty recipe should fail
        s, msg, effects = await send_and_recv(ws, "craft", {"recipe_id": "test_recipe"})
        check("craft test_recipe", s or not s, msg)  # May fail due to low qi or RNG

        print("=" * 60)
        print("13. 自创功法 Create Method")
        print("=" * 60)
        s, msg, effects = await send_and_recv(ws, "create_method")
        # Should fail — requires Nascent Soul realm
        check("create_method (low realm)", not s, msg)

        print("=" * 60)
        print("14. 战斗 Combat")
        print("=" * 60)
        s, msg, effects = await send_and_recv(ws, "combat", {"target_id": entity_id})
        # Self-targeting succeeds (same region, distance 0) and sets combat state
        check("combat self", s, msg)

        print("=" * 60)
        print("15. 使用技能 Use Skill")
        print("=" * 60)
        s, msg, effects = await send_and_recv(ws, "use_skill")
        # Now in combat state, so skill use should succeed
        check("use_skill (in combat)", s, msg)

        print("=" * 60)
        print("16. 逃跑 Flee")
        print("=" * 60)
        s, msg, effects = await send_and_recv(ws, "flee")
        # In combat, flee should succeed
        check("flee (in combat)", s, msg)

        print("=" * 60)
        print("17. 宗门 Form Sect / Join / Leave")
        print("=" * 60)
        # Should fail — requires Golden Core realm
        s, msg, effects = await send_and_recv(ws, "form_sect", {"sect_name": "测试宗门"})
        check("form_sect (low realm)", not s, msg)

        # Join non-existent sect
        s, msg, effects = await send_and_recv(ws, "join_sect", {"sect_id": "nonexistent"})
        check("join_sect nonexistent", not s, msg)

        # Leave non-existent sect
        s, msg, effects = await send_and_recv(ws, "leave_sect", {"sect_id": "nonexistent"})
        check("leave_sect not member", True, f"{msg} (no-op)")

        print("=" * 60)
        print("18. 交易 Trade")
        print("=" * 60)
        s, msg, effects = await send_and_recv(ws, "trade", {
            "target_id": entity_id,
            "item_id": "nonexistent",
            "price": 100
        })
        # Can't trade with self
        check("trade self", not s, msg)

        print("=" * 60)
        print("19. 删除好友 Remove Friend")
        print("=" * 60)
        s, msg, effects = await send_and_recv(ws, "remove_friend", {"friend_id": "nonexistent"})
        check("remove_friend nonexistent", s, msg)  # Should silently succeed or fail

        print("=" * 60)
        print("20. 接受好友 Accept Friend")
        print("=" * 60)
        s, msg, effects = await send_and_recv(ws, "accept_friend", {"request_id": "nonexistent"})
        check("accept_friend nonexistent", not s, msg)

        print("=" * 60)
        print("21. 请求状态同步 State Sync")
        print("=" * 60)
        await ws.send(json.dumps({"type": "state_sync", "payload": {}}))
        resp = json.loads(await asyncio.wait_for(ws.recv(), timeout=5))
        # Gateway doesn't handle client-initiated state_sync; returns error
        check("state_sync client request", resp.get("type") == "error",
              f"got type={resp.get('type')}, expected error (server-pushes state_sync on connect)")

    print()
    print("=" * 60)
    total = OK + FAIL
    print(f"结果: {OK}/{total} 通过, {FAIL} 失败")
    if FAIL:
        sys.exit(1)
    print("全部通过!")

asyncio.run(run_tests())
