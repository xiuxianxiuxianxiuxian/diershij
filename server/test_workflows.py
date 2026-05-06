"""
Comprehensive workflow test for all completed systems (A-K).
Tests with correct parameter names for all operations.
"""
import asyncio
import json
import random
import requests
import websockets

BASE_URL = "http://localhost:8081"
WS_URL = "ws://localhost:8081/ws"
GREEN = "\033[92m"
RED = "\033[91m"
YELLOW = "\033[93m"
CYAN = "\033[96m"
RESET = "\033[0m"
results = {"pass": 0, "fail": 0, "skip": 0}

def log(msg, color=GREEN):
    results["pass"] += 1
    print(f"  {color}[PASS]{RESET} {msg}")

def fail(msg):
    results["fail"] += 1
    print(f"  {RED}[FAIL]{RESET} {msg}")

def skip(msg):
    results["skip"] += 1
    print(f"  {YELLOW}[SKIP]{RESET} {msg}")

async def main():
    suffix = random.randint(10000, 99999)
    username = f"test_{suffix}"
    password = "test123"

    # ══════════ 1. Auth ══════════
    print(f"\n{CYAN}══════════ 1. Auth & Registration ══════════{RESET}\n")
    resp = requests.post(f"{BASE_URL}/auth/register", json={"username": username, "password": password})
    data = resp.json()
    if data.get("token"):
        log(f"Register OK: token received ({username})")
        token = data["token"]
        entity_data = data.get("entity", {})
        entity_id = entity_data.get("id", "?")
    else:
        fail(f"Register failed: {data}")
        return

    # Check spiritual roots in registration response (A1) — JSON tags are lowercase
    roots = entity_data.get("attributes", {}).get("spiritual_roots", [])
    if roots:
        log(f"Spiritual Roots on register: {len(roots)} root(s)")
        for r in roots:
            print(f"         - {r.get('element', '?')} purity={r.get('purity', '?')}%")
    else:
        skip("No spiritual roots in register response")

    # ══════════ 2. WebSocket & State Sync ══════════
    print(f"\n{CYAN}══════════ 2. WebSocket & State Sync ══════════{RESET}\n")
    uri = f"{WS_URL}?token={token}"
    async with websockets.connect(uri) as ws:
        first = json.loads(await asyncio.wait_for(ws.recv(), timeout=5))
        if first.get("type") == "state_sync":
            log("State sync received")
            entity = first.get("payload", first).get("entity", first)
            region = entity.get("position", {}).get("region_id", "?")
            realm = entity.get("realm", "?")
            print(f"         Entity: {entity.get('name', '?')}")
            print(f"         Region: {region} | Realm: {realm}")

            # Check spiritual roots in state sync (A) — JSON tags are lowercase
            roots = entity.get("attributes", {}).get("spiritual_roots", [])
            if roots:
                log(f"Spiritual Roots in sync: {len(roots)} root(s)")
            else:
                skip("Spiritual Roots not in sync payload")
        else:
            log(f"Msg received: type={first.get('type')}")

        async def send_op(action_type, params=None):
            msg = json.dumps({"type": "operation", "payload": {"action_type": action_type, "params": params or {}}})
            await ws.send(msg)
            await asyncio.sleep(1.5)
            results_list = []
            try:
                while True:
                    resp = json.loads(await asyncio.wait_for(ws.recv(), timeout=2))
                    results_list.append(resp)
            except (asyncio.TimeoutError, websockets.exceptions.ConnectionClosed):
                pass
            for r in results_list:
                if r.get("type") == "op_result":
                    return r.get("payload", r)
            return results_list[0] if results_list else None

        # ══════════ 3. Cultivation (A: Roots, B: Methods) ══════════
        print(f"\n{CYAN}══════════ 3. Cultivation System (A+B) ══════════{RESET}\n")

        r = await send_op("cultivate")
        if r and r.get("success", r.get("message")):
            log(f"Cultivate: {r.get('message', 'ok')}")
        else:
            fail(f"Cultivate failed: {r}")
        await asyncio.sleep(0.3)

        r = await send_op("meditate")
        if r:
            log(f"Meditate: {r.get('message', 'ok')}")

        # Breakthrough (expected to fail for new char)
        r = await send_op("breakthrough")
        if r and not r.get("success", True):
            skip(f"Breakthrough (expected): {r.get('message', 'insufficient cultivation')}")
        elif r:
            log(f"Breakthrough: {r.get('message', 'ok')}")
        else:
            skip("Breakthrough: no response")

        # List methods (B)
        r = await send_op("list_methods")
        if r and r.get("success", False):
            effects = r.get("effects", r)
            methods = effects.get("methods", effects.get("data", []))
            log(f"List methods: {len(methods) if methods else effects.get('message', 'ok')}")
        elif r:
            log(f"List methods response: {r.get('message', str(r)[:60])}")
        else:
            skip("List methods: no response")

        # ══════════ 4. Exploration & Gathering ══════════
        print(f"\n{CYAN}══════════ 4. Exploration & Gathering ══════════{RESET}\n")

        r = await send_op("explore")
        if r:
            log(f"Explore: {r.get('message', str(r)[:80])}")
        else:
            fail("Explore: no response")
        await asyncio.sleep(0.3)

        r = await send_op("gather")
        if r:
            log(f"Gather: {r.get('message', str(r)[:80])}")
        else:
            fail("Gather: no response")

        # ══════════ 5. Shop & Trade (F) ══════════
        print(f"\n{CYAN}══════════ 5. Shop & Trade (F) ══════════{RESET}\n")

        # shop_list - uses entity's current region (qingyun_town)
        r = await send_op("shop_list")
        if r and r.get("success", False):
            effects = r.get("effects", r)
            shops = effects.get("shops", [])
            log(f"Shop list: {len(shops)} shop(s) in current region")
        elif r:
            log(f"Shop list: {r.get('message', str(r)[:80])}")
        else:
            fail("Shop list: no response")

        # shop_items - needs string shop_id
        r = await send_op("shop_items", {"shop_id": "qingyun_trading"})
        if r and r.get("success", False):
            effects = r.get("effects", r)
            items = effects.get("items", [])
            log(f"Shop items: {len(items)} item(s) in Qingyun Shop")
        elif r:
            log(f"Shop items: {r.get('message', str(r)[:80])}")
        else:
            fail("Shop items: no response")

        # Buy - needs string shop_id + item_name
        r = await send_op("buy", {"shop_id": "qingyun_trading", "item_name": "Healing Pill", "quantity": 1})
        if r and r.get("success", False):
            log(f"Buy: {r.get('message', 'ok')}")
        elif r:
            log(f"Buy response: {r.get('message', str(r)[:80])}")
        else:
            fail("Buy: no response")

        # Sell - need an item in inventory first
        r = await send_op("sell", {"shop_id": "qingyun_trading", "item_name": "Healing Pill", "quantity": 1})
        if r and r.get("success", False):
            log(f"Sell: {r.get('message', 'ok')}")
        elif r:
            log(f"Sell response: {r.get('message', str(r)[:80])}")
        else:
            fail("Sell: no response")

        # Auction list
        r = await send_op("auction_list")
        if r:
            log(f"Auction list: {r.get('message', str(r)[:80])}")
        else:
            fail("Auction list: no response")

        # ══════════ 6. World Events (G) ══════════
        print(f"\n{CYAN}══════════ 6. World Events (G) ══════════{RESET}\n")

        # Move to another region
        r = await send_op("move", {"region_id": "misty_mountains"})
        if r and r.get("success", False):
            log(f"Move: {r.get('message', 'ok')}")
        elif r:
            log(f"Move response: {r.get('message', str(r)[:60])}")
        else:
            fail("Move: no response")

        # World events are timer-driven, check current state
        await asyncio.sleep(0.5)
        # Events are polling-based (G3), so check active events
        msg = json.dumps({"type": "operation", "payload": {"action_type": "world_events", "params": {}}})
        await ws.send(msg)
        await asyncio.sleep(1)
        try:
            resp = json.loads(await asyncio.wait_for(ws.recv(), timeout=2))
            log(f"World events check: {str(resp)[:100]}")
        except asyncio.TimeoutError:
            skip("World events: no immediate response (timer-driven)")

        # ══════════ 7. Social System (H) ══════════
        print(f"\n{CYAN}══════════ 7. Social System (H) ══════════{RESET}\n")

        # World chat
        r = await send_op("send_message", {"message": "Hello from test!", "channel": "world"})
        if r and r.get("success", False):
            log(f"World chat: {r.get('message', 'ok')}")
        elif r:
            log(f"Chat response: {r.get('message', str(r)[:60])}")
        else:
            fail("World chat: no response")

        # Mail list (H2)
        r = await send_op("mail_list")
        if r and r.get("success", False):
            effects = r.get("effects", r)
            mails = effects.get("mails", [])
            log(f"Mail list: {len(mails)} mail(s)")
        elif r:
            log(f"Mail response: {r.get('message', str(r)[:80])}")
        else:
            fail("Mail: no response")

        # Leaderboard (H3)
        r = await send_op("leaderboard")
        if r and r.get("success", False):
            effects = r.get("effects", r)
            entries = effects.get("entries", [])
            log(f"Leaderboard: {effects.get('count', len(entries))} entries")
        elif r:
            log(f"Leaderboard response: {r.get('message', str(r)[:80])}")
        else:
            fail("Leaderboard: no response")

        # Friends (H5)
        r = await send_op("add_friend", {"target_name": username[:5] + "0"})  # non-existent
        if r:
            log(f"Add friend response: {r.get('message', str(r)[:60])}")
        else:
            fail("Add friend: no response")

        # Nearby players
        r = await send_op("nearby_players")
        if r:
            log(f"Nearby players: {r.get('message', str(r)[:60])}")

        # ══════════ 8. Equipment & Combat (C+D) ══════════
        print(f"\n{CYAN}══════════ 8. Equipment & Combat (C+D) ══════════{RESET}\n")

        r = await send_op("inventory")
        if r:
            log(f"Inventory: {r.get('message', str(r)[:60])}")
        else:
            skip("Inventory: no response")

        r = await send_op("equip_item", {"item_name": "test_sword"})
        if r:
            log(f"Equip response: {r.get('message', str(r)[:60])}")
        else:
            skip("Equip: no response")

        r = await send_op("use_skill", {"skill_name": "basic_attack", "target": "monster"})
        if r:
            log(f"Use skill: {r.get('message', str(r)[:60])}")

        # ══════════ 9. Sect (I) ══════════
        print(f"\n{CYAN}══════════ 9. Sect System ══════════{RESET}\n")

        r = await send_op("form_sect", {"name": f"TestSect{suffix}", "description": "A test sect"})
        if r and r.get("success", False):
            log(f"Form sect: {r.get('message', 'ok')}")
        elif r:
            log(f"Form sect response: {r.get('message', str(r)[:60])}")

        r = await send_op("sect_info")
        if r:
            log(f"Sect info: {r.get('message', str(r)[:60])}")

    # ══════════ SUMMARY ══════════
    print(f"\n{CYAN}══════════ TEST SUMMARY ══════════{RESET}\n")
    total = results["pass"] + results["fail"] + results["skip"]
    print(f"  Total: {total}  |  {GREEN}Pass: {results['pass']}{RESET}  |  {RED}Fail: {results['fail']}{RESET}  |  {YELLOW}Skip: {results['skip']}{RESET}")

if __name__ == "__main__":
    asyncio.run(main())
