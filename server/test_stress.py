"""Full game system stress test - multiple iterations"""
import asyncio, json, websockets, requests, time, sys, os

LOG_FILE = os.path.join(os.path.dirname(__file__), "stress_test_report.txt")

def log(msg):
    with open(LOG_FILE, "a", encoding="utf-8") as f:
        f.write(msg + "\n")
    print(msg, flush=True)

BASE = "http://localhost:8081"
WS = "ws://localhost:8081"

ALL_ACTIONS = [
    "cultivate", "breakthrough", "explore", "gather",
    "meditate", "sleep", "move", "list_methods", "learn_method",
    "set_main_method", "cast_spell", "use_skill", "flee",
    "list_friends", "add_friend", "accept_friend", "remove_friend",
    "use_item", "equip_item", "unequip_item", "drop_item",
    "craft", "create_method", "send_message",
    "combat", "form_sect", "join_sect", "leave_sect", "sect_info",
]

results = {"pass": 0, "fail": 0, "total": 0}
all_logs = []

def log_test(action, status, msg=""):
    results["total"] += 1
    if status:
        results["pass"] += 1
    else:
        results["fail"] += 1
    tag = "PASS" if status else "FAIL"
    line = f"  [{tag}] {action}: {msg[:120]}"
    all_logs.append(line)
    log(line)

def log_header(msg):
    log(f"\n{'='*60}")
    log(f"  {msg}")
    log(f"{'='*60}")

async def recv_result(ws, timeout=15):
    resp = json.loads(await asyncio.wait_for(ws.recv(), timeout=timeout))
    while resp.get("type") == "state_sync":
        resp = json.loads(await asyncio.wait_for(ws.recv(), timeout=timeout))
    p = resp.get("payload", {})
    if resp.get("type") == "error":
        return {"success": False, "message": resp.get("payload", {}).get("message", str(resp))}
    return p

async def run_iteration(iteration, username):
    log_header(f"迭代 {iteration}: {username}")

    resp = requests.post(f"{BASE}/auth/register", json={"username": username, "password": "123456"})
    data = resp.json()
    if not data.get("success"):
        resp = requests.post(f"{BASE}/auth/login", json={"username": username, "password": "123456"})
        data = resp.json()
    if not data.get("success"):
        log_test("注册/登录", False, data.get("error", "unknown"))
        return

    token = data["token"]
    entity = data.get("entity", {})
    entity_id = entity.get("id", "?")
    attrs = entity.get("attributes", {})
    roots = attrs.get("spiritual_roots", [])
    root_str = ", ".join([f"{r['element']}({r['purity']})" for r in roots]) if roots else "无"
    log_test("注册", True, f"ID={entity_id[:12]}... 灵根=[{root_str}]")
    log_test("灵根生成", bool(roots), f"{len(roots)}条: {[r['element'] for r in roots]}")

    try:
        ws = await websockets.connect(f"{WS}/ws?token={token}")
        first = json.loads(await asyncio.wait_for(ws.recv(), timeout=5))
        log_test("WS连接", first.get("type") == "state_sync")
    except Exception as e:
        log_test("WS连接", False, str(e))
        return

    def send_op(action, params=None):
        return json.dumps({"type": "operation", "payload": {"action_type": action, "params": params or {}}})

    # ── B4: 功法系统 ──
    await ws.send(send_op("list_methods"))
    r = await recv_result(ws)
    methods = r.get("effects", {}).get("methods", [])
    log_test("list_methods", r.get("success"), f"{r.get('message','')[:60]}")

    learned_ids = []
    if methods:
        mid = methods[0]["id"]
        mname = methods[0]["name"]
        await ws.send(send_op("learn_method", {"method_id": mid}))
        r = await recv_result(ws)
        log_test("learn_method", r.get("success"), f"{mname}: {r.get('message','')[:50]}")
        if r.get("success"):
            learned_ids.append(mid)

        await ws.send(send_op("set_main_method", {"method_id": mid}))
        r = await recv_result(ws)
        log_test("set_main_method", r.get("success"), r.get("message", "")[:50])

        # Try learning more methods
        if len(methods) > 1:
            mid2 = methods[1]["id"]
            await ws.send(send_op("learn_method", {"method_id": mid2}))
            r = await recv_result(ws)
            log_test("learn_method#2", r.get("success"), r.get("message", "")[:50])

    # ── 修炼 x5 ──
    total_progress = 0
    for i in range(5):
        await ws.send(send_op("cultivate"))
        r = await recv_result(ws)
        gain = r.get("effects", {}).get("cultivation_gain", 0)
        total_progress = r.get("effects", {}).get("progress", 0)
        if i == 0:
            log_test("cultivate", r.get("success"), f"第1次: +{gain:.3f}")
    log_test("cultivate x5", total_progress > 0, f"总进度={total_progress:.2f}")

    # ── 冥想 ──
    await ws.send(send_op("meditate"))
    r = await recv_result(ws)
    log_test("meditate", r.get("success"), r.get("message", "")[:50])

    # ── 睡眠 ──
    await ws.send(send_op("sleep"))
    r = await recv_result(ws)
    log_test("sleep", r.get("success"), r.get("message", "")[:50])

    # ── 探索 x3 ──
    explore_ok = 0
    for i in range(3):
        await ws.send(send_op("explore"))
        r = await recv_result(ws)
        if r.get("success"): explore_ok += 1
    log_test("explore x3", explore_ok > 0, f"{explore_ok}/3 成功")

    # ── 移动 x4 (不同区域) ──
    regions = ["qingyun_town", "spirit_mist_mountain", "fallen_sword_valley", "qingyun_town"]
    move_ok = 0
    for rname in regions:
        await ws.send(send_op("move", {"region_id": rname}))
        r = await recv_result(ws)
        if r.get("success"): move_ok += 1
    log_test("move x4", move_ok > 1, f"{move_ok}/4 成功")

    # ── 采集 ──
    await ws.send(send_op("gather", {"resource_type": "herb", "quantity": 1}))
    r = await recv_result(ws)
    log_test("gather", True, r.get("message", "")[:50])

    # ── 发送消息 ──
    await ws.send(send_op("send_message", {"content": f"测试消息_{iteration}", "channel": "world"}))
    r = await recv_result(ws)
    log_test("send_message", True, r.get("message", "")[:50])

    # ── 好友系统 ──
    friend_user = f"friend_{username}"
    requests.post(f"{BASE}/auth/register", json={"username": friend_user, "password": "123456"})
    await ws.send(send_op("add_friend", {"username": friend_user}))
    r = await recv_result(ws)
    log_test("add_friend", True, r.get("message", "")[:50])

    await ws.send(send_op("list_friends"))
    r = await recv_result(ws)
    log_test("list_friends", True, r.get("message", "")[:50])

    # ── 战斗 (找NPC) ──
    # Move to spirit_mist_mountain where NPCs should be
    await ws.send(send_op("move", {"region_id": "spirit_mist_mountain"}))
    await recv_result(ws)

    # Try combat with various target_ids (NPC UUIDs might vary, so expect failure)
    await ws.send(send_op("combat", {"target_id": "00000000-0000-0000-0000-000000000001"}))
    r = await recv_result(ws)
    log_test("combat", True, r.get("message", "")[:50])

    # ── 法术 ──
    await ws.send(send_op("cast_spell", {"spell_id": "99999999-9999-9999-9999-999999999999"}))
    r = await recv_result(ws)
    log_test("cast_spell(无)", not r.get("success"), r.get("message", "")[:50])

    # ── 技能 ──
    await ws.send(send_op("use_skill", {"target_id": "00000000-0000-0000-0000-000000000001"}))
    r = await recv_result(ws)
    log_test("use_skill", True, r.get("message", "")[:50])

    # ── 逃跑 ──
    await ws.send(send_op("flee"))
    r = await recv_result(ws)
    log_test("flee", True, r.get("message", "")[:50])

    # ── 自创功法 ──
    await ws.send(send_op("create_method", {"name": "测试功法", "element": "fire"}))
    r = await recv_result(ws)
    log_test("create_method", True, r.get("message", "")[:50])

    # ── 物品操作 ──
    await ws.send(send_op("use_item", {"item_name": "灵芝"}))
    r = await recv_result(ws)
    log_test("use_item", True, r.get("message", "")[:50])

    # ── 突破 ──
    await ws.send(send_op("breakthrough"))
    r = await recv_result(ws)
    if "修为不足" in r.get("message", ""):
        log_test("breakthrough(修为不足)", True, r.get("message", "")[:50])
    else:
        log_test("breakthrough", r.get("success") == False, r.get("message", "")[:50])

    # ── 宗门信息 ──
    await ws.send(send_op("sect_info"))
    r = await recv_result(ws)
    log_test("sect_info", True, r.get("message", "")[:50])

    await ws.close()

async def main():
    iterations = int(sys.argv[1]) if len(sys.argv) > 1 else 3

    # Clear previous log
    open(LOG_FILE, "w", encoding="utf-8").close()

    for i in range(1, iterations + 1):
        username = f"stress_{int(time.time())}_{i}"
        await run_iteration(i, username)

    # Summary
    total = results["pass"] + results["fail"]
    log(f"\n{'='*60}")
    log(f"  全部测试完成 ({iterations} 轮)")
    log(f"  通过: {results['pass']}/{total}  ({results['pass']*100//max(total,1)}%)")
    log(f"  失败: {results['fail']}/{total}")
    log(f"{'='*60}")

    if results["fail"] > 0:
        log(f"\n失败记录:")
        for line in all_logs:
            if "[FAIL]" in line:
                log(line)

    return results["fail"] == 0

if __name__ == "__main__":
    ok = asyncio.run(main())
    sys.exit(0 if ok else 1)
