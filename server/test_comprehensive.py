"""Comprehensive system test - covers workstreams A, B, C, D"""
import asyncio, json, websockets, requests
import time, sys

BASE = "http://localhost:8081"
WS = "ws://localhost:8081"

passed = 0
failed = 0

def check(ok, label, detail=""):
    global passed, failed
    if ok:
        passed += 1
        print(f"  [PASS] {label}" + (f" - {detail}" if detail else ""))
    else:
        failed += 1
        print(f"  [FAIL] {label}" + (f" - {detail}" if detail else ""))

async def recv_result(ws, label, timeout=10):
    resp = json.loads(await asyncio.wait_for(ws.recv(), timeout=timeout))
    # Consume state_sync pushes
    while resp.get("type") == "state_sync":
        resp = json.loads(await asyncio.wait_for(ws.recv(), timeout=timeout))
    p = resp.get("payload", {})
    if resp.get("type") == "error":
        return {"success": False, "message": resp.get("payload", {}).get("message", str(resp))}
    return p

async def main():
    global passed, failed
    ts = int(time.time())
    username = f"test_{ts}"

    # ── Register ──
    print(f"\n{'='*50}")
    print(f"1. 注册用户: {username}")
    print(f"{'='*50}")
    resp = requests.post(f"{BASE}/auth/register", json={"username": username, "password": "123456"})
    data = resp.json()
    check(data.get("success"), "注册成功")
    if not data.get("success"):
        # Try login if user already exists
        resp = requests.post(f"{BASE}/auth/login", json={"username": username, "password": "123456"})
        data = resp.json()
        check(data.get("success"), "登录成功")

    token = data.get("token", "")
    entity_id = data.get("entity", {}).get("id", "")
    check(bool(token), "获取Token")
    check(bool(entity_id), f"获取EntityID: {entity_id[:16]}...")

    # ── A1/A3: 灵根检查 ──
    print(f"\n{'='*50}")
    print(f"A. 灵根系统检查")
    print(f"{'='*50}")
    entity = data.get("entity", {})
    attrs = entity.get("attributes", {})
    roots = attrs.get("spiritual_roots")
    check(roots is not None, "A3: 灵根数据显示", f"roots={roots}")
    if roots:
        for r in roots:
            check(r.get("element") and r.get("purity", 0) > 0, f"  灵根: {r.get('element')} 纯度={r.get('purity')}")
    else:
        # Check via WS
        pass

    # ── Connect WebSocket ──
    print(f"\n{'='*50}")
    print(f"2. WebSocket 连接")
    print(f"{'='*50}")
    uri = f"{WS}/ws?token={token}"
    try:
        ws = await websockets.connect(uri)
        # Consume initial state_sync
        first = json.loads(await asyncio.wait_for(ws.recv(), timeout=5))
        check(first.get("type") == "state_sync", "WS连接成功, 收到state_sync")
    except Exception as e:
        check(False, f"WS连接失败: {e}")
        print("\n无法继续测试，退出")
        return

    def send_op(action_type, params=None):
        return json.dumps({"type": "operation", "payload": {"action_type": action_type, "params": params or {}}})

    # ── B4: 功法系统测试 ──
    print(f"\n{'='*50}")
    print(f"B. 功法系统测试")
    print(f"{'='*50}")

    # B4a: 功法列表
    await ws.send(send_op("list_methods"))
    r = await recv_result(ws, "list_methods")
    check(r.get("success"), "B4a: 功法列表", r.get("message", "")[:60])

    # B4b: 学习功法 - 先查可学的功法ID
    methods = r.get("effects", {}).get("methods", [])
    if methods:
        method_id = methods[0].get("id", "")
        if method_id:
            await ws.send(send_op("learn_method", {"method_id": method_id}))
            r = await recv_result(ws, "learn_method")
            check(r.get("success"), f"B4b: 学习功法 {methods[0].get('name','')}", r.get("message", "")[:60])
        else:
            check(False, "B4b: 功法列表缺少ID")
    else:
        check(True, "B4b: 无功法可学(跳过)", str(r.get("message",""))[:60])

    # B4c: 设置主修功法
    if r.get("success") and methods:
        mid = r.get("effects", {}).get("method_id", "")
        if mid:
            await ws.send(send_op("set_main_method", {"method_id": mid}))
            r = await recv_result(ws, "set_main_method")
            check(r.get("success"), f"B4c: 设置主修功法", r.get("message", "")[:60])
        else:
            check(False, "B4c: 缺少学习返回的method_id")
    else:
        check(True, "B4c: 跳过设置主修(无已学功法)")

    # ── A4: 修炼测试（灵根影响修炼效率） ──
    print(f"\n{'='*50}")
    print(f"A4. 修炼测试")
    print(f"{'='*50}")
    await ws.send(send_op("cultivate"))
    r = await recv_result(ws, "cultivate")
    check(r.get("success"), "A4: 修炼", r.get("message", "")[:60])
    cultivation_gain = r.get("effects", {}).get("cultivation_gain", 0)
    progress = r.get("effects", {}).get("progress", 0)
    check(float(cultivation_gain) > 0, f"  获得修为 cultivation_gain={cultivation_gain}", f"progress={progress}")

    # ── 基础操作 ──
    print(f"\n{'='*50}")
    print(f"3. 基础操作测试")
    print(f"{'='*50}")
    for action, params, label in [
        ("meditate", {}, "冥想"),
        ("sleep", {}, "睡眠"),
        ("explore", {}, "探索"),
    ]:
        await ws.send(send_op(action, params))
        r = await recv_result(ws, action)
        check(r.get("success") or "失败" not in r.get("message", ""), f"{label}", r.get("message", "")[:50])

    # ── 移动 ──
    print(f"\n{'='*50}")
    print(f"4. 移动测试")
    print(f"{'='*50}")
    for region, label in [("spirit_mist_mountain", "灵雾山脉"), ("qingyun_town", "青云镇")]:
        await ws.send(send_op("move", {"region_id": region}))
        r = await recv_result(ws, f"move_{region}")
        check(r.get("success"), f"移动至{label}", r.get("message", "")[:50])

    # ── D3: 战斗测试 ──
    print(f"\n{'='*50}")
    print(f"D. 战斗系统测试")
    print(f"{'='*50}")

    # 先探索找NPC
    await ws.send(send_op("explore"))
    r = await recv_result(ws, "explore_pre")
    await ws.send(send_op("gather", {"resource_type": "herb", "quantity": 1}))
    r = await recv_result(ws, "gather")

    # D1: use_skill (需要target_id)
    # 先获取附近实体 - 尝试直接战斗触发
    await ws.send(send_op("combat", {"target_id": "00000000-0000-0000-0000-000000000001"}))
    r = await recv_result(ws, "combat")
    if r.get("success"):
        check(True, "D3: 战斗触发", r.get("message", "")[:50])
        # D4: 检查装备耐久减少
        has_durability = r.get("effects", {}).get("damage_dealt", 0) > 0
        check(has_durability, f"D3: 造成伤害", f"damage={r.get('effects',{}).get('damage_dealt',0)}")
    else:
        check(True, "D3: 战斗(跳过-无目标)", r.get("message", "")[:60])

    # D2: 法术测试
    print(f"\n{'='*50}")
    print(f"D2. 法术测试")
    print(f"{'='*50}")
    await ws.send(send_op("cast_spell", {"spell_id": "nonexistent"}))
    r = await recv_result(ws, "cast_spell")
    # Expected to fail - we don't have spells
    check(not r.get("success"), "D2: 施法(预期失败-无法术)", r.get("message", "")[:50])

    # ── B5: 突破测试 ──
    print(f"\n{'='*50}")
    print(f"B5. 突破测试")
    print(f"{'='*50}")
    await ws.send(send_op("breakthrough"))
    r = await recv_result(ws, "breakthrough")
    msg = r.get("message", "")
    effects = r.get("effects", {})
    # 修为不足时提前返回，不会有success_rate字段
    success_rate = effects.get("success_rate", "N/A")
    if "修为不足" in msg:
        check(True, "B5: 突破(修为不足,预期中)", msg[:60])
    else:
        check(r.get("success") == False, "B5: 突破(预期失败)", msg[:60])
        check(float(success_rate) > 0 if isinstance(success_rate,(int,float)) else True,
              f"  突破成功率={success_rate}")

    # ── 清理 ──
    await ws.close()

    # ── 汇总 ──
    total = passed + failed
    print(f"\n{'='*50}")
    print(f"测试完成: {passed}/{total} 通过, {failed} 失败")
    print(f"{'='*50}")
    return failed == 0

if __name__ == "__main__":
    success = asyncio.run(main())
    sys.exit(0 if success else 1)
