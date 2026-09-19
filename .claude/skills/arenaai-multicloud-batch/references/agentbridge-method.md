# AgentBridge 通路（chrome.debugger 经桌面扩展桥接）

> 当 `chrome:control-chrome` 因 Chrome 153 封锁 `--remote-debugging-port`、或其 relay 拥塞整批崩溃时，改用本通路驱动用户真实 Chrome。仍是"原生 fill → press Enter → 核验清空"序列，只是换了传输层。

## 何时启用

满足任一即可改用本通路（`chrome:control-chrome` 为默认，本通路为备选）：

- Chrome 153+ 在 Default profile 上拒绝 `--remote-debugging-port`（managed/headless Chrome 拿不到 CDP 端口）。
- 同时对战 ≥15 页时 `chrome:control-chrome` 的 relay 变慢，`get_text` 30s 超时并整批崩溃。
- 用户希望复用已登录 BUC 的真实 Chrome，而非另开 managed profile 重新登录。

不满足上述条件时仍走 `chrome:control-chrome`（见 [页面操作与恢复](arenaai-operations.md)）。

## 关键差异（为何 evaluate 在此通路能用）

ArenaAI 页面 CSP 不含 `'unsafe-eval'`。**页面上下文**的 `evaluate`（`chrome:control-chrome` 的 page eval、playwright 的 page.evaluate）受 CSP 拦截，任意 JS 无法运行——这是 [页面操作与恢复](arenaai-operations.md) "CSP 约束" 一节的事实基础，对那条通路仍然成立。

但本通路用的是扩展的 `chrome.debugger.attach({tabId},"1.3")`，其 `Runtime.evaluate` 是 **CDP 级**（devtools 协议层），**不受页面 CSP 限制**——即使页面 CSP 禁 `unsafe-eval`，CDP evaluate 仍可执行任意表达式。这是两条通路的核心区别：同样写"原生 setter + dispatchEvent"片段，page-context eval 被 CSP 拦，CDP eval 能跑。

**但**：仍应使用原生 fill（native setter + input/change events）而非任意 JS 点击——原因是 React 受控组件的兼容性，不是 CSP。`press Enter` 经 `chrome.debugger` 的 Input.dispatchKeyEvent 也能触发 React onSubmit。不要因为 evaluate 能跑了就改用 JS 模拟点击跳过原生命令序列。

## 一次性准备

1. **桌面扩展**：用户本地 `~/Desktop/browser-control-extension/`（manifest_v3，permissions 含 `debugger`、`tabs`、`scripting`、`<all_urls>`）。在用户真实 Chrome 的 `chrome://extensions` 加载为"已解压的扩展程序"，保持启用。扩展用 `chrome.debugger.attach` 附加各 tab——这是扩展权限级，**不依赖** `--remote-debugging-port`，故绕过 Chrome 153 的封锁。
2. **daemon 脚本**：把下文 [daemon 脚本](#daemon-脚本) 写到 `/tmp/agentbridge_daemon.py`，`pip install websockets` 后 `nohup python3 /tmp/agentbridge_daemon.py > /tmp/agentbridge_daemon.stdout 2>&1 &` 启动。WS 监听 `127.0.0.1:10087`（扩展连过来），HTTP 监听 `127.0.0.1:10088`（curl 调过来）。日志 `/tmp/agentbridge_daemon.log`。
3. **连接验证**：`curl -s http://127.0.0.1:10088/status` 应返回 `{"ok":true,"connected":true}`。`connected:false` 表示扩展未连（检查扩展是否启用、端口是否被占）。
4. **BUC 登录**：真实 Chrome 的 Default profile 已持久登录 ArenaAI（cookie 跨重启存活，title 为 "Arena" 而非登录页）。本通路直接复用，**无需重登**。若 cookie 确实过期，请用户在真实 Chrome 里手动完成 BUC SSO，不要在 daemon/扩展侧读 cookie。

## 协议

- 扩展 → daemon：`{type:"hello",data:{...}}`（连接元信息）；`{type:"response",requestId,data}` 或 `{...,error}`（命令回复）。
- daemon → 扩展：`{type:"command",requestId,session,action,args}`。
- daemon 用 `requestId` 把命令-回复配对到 `asyncio.Future`，`send_cmd` await 该 future，超时由 `timeout` 控制。

## HTTP API（curl 侧）

- `POST /cmd`，body `{"action","args","session","timeout"}`，args 里通常带 `tabId`。返回 `{"ok":true,"data":{...}}` 或 `{"ok":false,"error":...}`。
- `GET /status` → `{"ok":true,"connected:<bool>}`。
- `GET /shutdown` → daemon 退出。

shell 包装示例：

```bash
ab() {  # ab <action> '<args json>' [timeout]
  curl -s -m 45 -X POST http://127.0.0.1:10088/cmd \
    -d "{\"action\":\"$1\",\"args\":${2:-{\}},\"timeout\":${3:-30}}"
}
ab_status() { curl -s -m 5 http://127.0.0.1:10088/status; }
```

## 命令面（扩展实现，本通路用到的子集）

| action | 关键 args | 说明 |
|---|---|---|
| `evaluate` | `code`、`tabId` | CDP `Runtime.evaluate`；**注意表达式/函数体模式坑**（见下） |
| `press` | `key`、`selector`、`tabId` | 经 debugger 的键盘事件；`Enter` 触发 React onSubmit |
| `navigate` | `url`、`newTab`、`timeoutMs` | 同窗口新 tab 或现有 tab 跳转 |
| `list_tabs` | — | 列当前窗口所有 tab（拿 tabId） |
| `fill` | `selector`、`value`、`tabId` | chrome.scripting 注入 fill；**但推荐用 evaluate 跑原生 setter**（见下） |

## 填入（evaluate 投递原生 setter）

不用扩展的 `fill` action（其 commit 策略与 React 受控组件偶有脱节），改用 `evaluate` 执行原生 setter片段——这是 React 兼容写法，CSP 下也能跑：

```python
import json
SEL = 'textarea[placeholder="输入问题，即刻开始"]'
def fill_native(tid, q):
    code = ("var ta=document.querySelector(" + json.dumps(SEL) + ");"
            "if(!ta) return 'NO_TA';"
            "var s=Object.getOwnPropertyDescriptor(window.HTMLTextAreaElement.prototype,'value').set;"
            "s.call(ta," + json.dumps(q) + ");"
            "ta.dispatchEvent(new Event('input',{bubbles:true}));"
            "ta.dispatchEvent(new Event('change',{bubbles:true}));"
            "return ta.value;")
    return ev(code, tid)   # ev = evaluate 封装
```

`HTMLTextAreaElement.prototype.value` 的原生 setter set + `input`/`change` 事件，是 React 受控组件把外部设值收进内部 state 的标准姿势。**不要**用 `ta.value = q`（React 会覆盖）或 IIFE 包裹（见下）。

## 提交与核验

```python
# 提交前读 uc 基线
uc0 = ev("return document.querySelectorAll('.user-message').length;", tid)
# focus 后 press Enter
ev("var ta=document.querySelector("+json.dumps(SEL)+");if(ta)ta.focus();return ta?'ok':'NO_TA';", tid)
ab('press', json.dumps({"key":"Enter","selector":SEL,"tabId":tid}))
sleep(6)
uc_a, ta_val = read_uc_ta(tid)   # 读 uc 与 textarea 值
```

**提交主信号**：textarea 清空（`ta_val==''`）或 `uc_a > uc0`。**详见 [页面操作与恢复](arenaai-operations.md) 的"幂等发送"与"状态定义"——本通路沿用同一判据与状态机，不重复定义。**

## 坑：performCdpEvaluate 表达式 vs 函数体模式

`chrome.debugger` 的 `Runtime.evaluate`（扩展里 `performCdpEvaluate`）按代码内容分两种模式：

- **表达式模式**：代码**不含** `return` 关键字 → 整段当单个表达式，其值即返回值。例 `document.querySelector('textarea').value`。
- **函数体模式**：代码**含** `return` 关键字 → 整段当函数体，须显式 `return` 才有返回值。例 `var ta=...; return ta.value;`。

**IIFE 陷阱**：`(function(){return ta.value;})()` **含** `return` → 触发函数体模式 → 整段被当函数体，而作为最后一条表达式语句的 IIFE **本身没有 return** → 返回 `undefined`。R4 首次 fill 全员 NOTSYNC 就是因为包了 IIFE。

**解法**：写成一串语句 + 尾部 `return`，**不要包 IIFE**：

```javascript
// 正确（函数体模式，有 return）
var ta=document.querySelector(sel);
if(!ta) return 'NO_TA';
var s=Object.getOwnPropertyDescriptor(window.HTMLTextAreaElement.prototype,'value').set;
s.call(ta,q);
ta.dispatchEvent(new Event('input',{bubbles:true}));
return ta.value;
```

返回值若是对象，evaluate 的 `data.result` 可能是 `{value:...}` 结构，调用侧要再取一层 `.value`。

## 坑：提交后 React 重渲染间隙（假阴性）

`press Enter` 后约 6s，React 把新 user message 挂进 DOM 列表，**重渲染瞬间 `.user-message` 列表会被清空再重挂**——此时读 `uc` 可能得到 **0**（比基线还小）。R5 tab 14 就出现 `uc 1->0`，被误判 `submitted=False`，但 t+75s 验收时 `uc=2 / wrap=4 / maxLen=18047` 证明实际已提交并获答。

**解法**：

1. 读 `uc_a` 后若 `uc_a < uc0`（异常回退）或 `ta_val` 仍非空，**sleep 2s 重读一次**再判 submitted。
2. **t+75s 的验收是 ground truth**：`completed = uc>=基线+1 and !hasStop and taEmpty and maxLen>50`。Stage 2 的 `submitted` 测量瑕疵不推翻 Stage 4 的 completed 判定。
3. 验收脚本里 `uc` 与 `wrap` 与 `maxMsgLen` 须一致递增（每轮 +1 user / +2 wrap），三者不一致即怀疑测量落在重渲染窗口，补一次重读。

## 坑：问题必须是字符串，不是 dict

问题集 JSON 每条是 `{"i":..,"tab":..,"theme":..,"question":"..."}`。fill 前必须取 `d['question']`（字符串）传入，**不能**直接把整个 dict 当问题文本——否则 `json.dumps(dict)` 进 setter，textarea 被填成字面量 `[object Object]`（恰好 15 字符），synced 全 False。R4 首轮踩过。

## 坑：macOS 无 `timeout` 命令

daemon 超时由 `asyncio.wait_for` 在脚本内控；shell 侧若要限时跑命令，用 `perl -e 'alarm shift @ARGV; exec @ARGV' <N> <cmd>`，不要用 `timeout`（macOS 无）。

## 验收计数（与 [页面操作与恢复](arenaai-operations.md) "最终验收"一致）

每轮跑完后汇报：`fill synced N/N | submit N/N | completed N/N`，并落盘 `/tmp/arena_ledger_rN.json`（含 `tabId`、`q`、`hash`、`state`、`submit`、`verify`、`retry`）。`submit` 字段的 `submitted:false` 若伴随 `verify.completed:true`，按 completed 上报并在汇报里注明该 tab 为 Stage 2 瞬时假阴性（不重发、不重试，验收已证完成）。

## daemon 脚本

自包含，写到 `/tmp/agentbridge_daemon.py` 即可：

```python
#!/usr/bin/env python3
"""AgentBridge daemon: WS server (10087) for the Browser Control extension,
HTTP API (10088) for local curl-driven commands. Bridges the extension's
chrome.debugger / chrome.scripting capabilities to simple curl calls so we
can drive the user's real Chrome (Default profile, already logged into BUC)
without playwright-cli's congesting relay and without --remote-debugging-port
(which Chrome 153 blocks on the default profile)."""
import asyncio
import json
import sys
import os

try:
    import websockets
except ImportError:
    print("FATAL: websockets not installed", file=sys.stderr)
    sys.exit(1)

SESSION = "arena"
ext_ws = {"ws": None}
pending = {}
req_counter = [0]
LOG = open("/tmp/agentbridge_daemon.log", "a", buffering=1)


def log(msg):
    line = f"[{os.getpid()}] {msg}"
    print(line, flush=True)
    LOG.write(line + "\n")


async def ws_handler(websocket, *args):
    ext_ws["ws"] = websocket
    log(f"extension CONNECTED path={getattr(websocket, 'path', '?')}")
    try:
        async for msg in websocket:
            try:
                data = json.loads(msg)
            except Exception:
                log(f"non-json msg: {str(msg)[:120]}")
                continue
            t = data.get("type")
            if t == "hello":
                meta = data.get("data") or {}
                log(f"hello from extension: {json.dumps(meta)[:200]}")
                continue
            if t == "response":
                rid = data.get("requestId")
                fut = pending.get(rid)
                if fut and not fut.done():
                    if data.get("error"):
                        fut.set_exception(Exception(str(data["error"])))
                    else:
                        fut.set_result(data.get("data"))
                continue
            log(f"ext msg type={t}: {json.dumps(data)[:160]}")
    except websockets.ConnectionClosed:
        pass
    except Exception as e:
        log(f"ws handler error: {e}")
    finally:
        ext_ws["ws"] = None
        log("extension DISCONNECTED")


async def send_cmd(action, args=None, session=SESSION, timeout=30):
    ws = ext_ws["ws"]
    if ws is None:
        raise RuntimeError("extension not connected (load the extension in Chrome)")
    req_counter[0] += 1
    rid = req_counter[0]
    loop = asyncio.get_running_loop()
    fut = loop.create_future()
    pending[rid] = fut
    payload = {"type": "command", "requestId": rid, "session": session,
               "action": action, "args": args or {}}
    await ws.send(json.dumps(payload))
    try:
        return await asyncio.wait_for(fut, timeout)
    finally:
        pending.pop(rid, None)


async def http_handler(reader, writer):
    try:
        line = await asyncio.wait_for(reader.readline(), 5)
        if not line:
            writer.close()
            return
        headers = {}
        while True:
            h = await asyncio.wait_for(reader.readline(), 5)
            if h in (b"\r\n", b"\n", b""):
                break
            k, _, v = h.decode(errors="ignore").partition(":")
            headers[k.strip().lower()] = v.strip()
        cl = int(headers.get("content-length", "0") or "0")
        body = b""
        if cl > 0:
            body = await asyncio.wait_for(reader.readexactly(cl), 5)
        parts = line.decode(errors="ignore").split()
        path = parts[1] if len(parts) > 1 else "/"

        if path == "/cmd":
            try:
                req = json.loads(body.decode() or "{}")
                action = req["action"]
                args = req.get("args", {})
                session = req.get("session", SESSION)
                timeout = float(req.get("timeout", 30))
                data = await send_cmd(action, args, session, timeout)
                resp = json.dumps({"ok": True, "data": data})
            except Exception as e:
                resp = json.dumps({"ok": False, "error": str(e)})
        elif path == "/status":
            resp = json.dumps({"ok": True, "connected": ext_ws["ws"] is not None})
        elif path == "/shutdown":
            resp = json.dumps({"ok": True, "bye": True})
            payload = resp.encode()
            hdr = (f"HTTP/1.1 200 OK\r\nContent-Type: application/json\r\n"
                   f"Content-Length: {len(payload)}\r\n\r\n").encode()
            writer.write(hdr + payload)
            await writer.drain()
            writer.close()
            os._exit(0)
        else:
            resp = json.dumps({"ok": False, "error": "use POST /cmd or GET /status"})
        payload = resp.encode()
        hdr = (f"HTTP/1.1 200 OK\r\nContent-Type: application/json\r\n"
               f"Content-Length: {len(payload)}\r\n"
               f"Access-Control-Allow-Origin: *\r\n\r\n").encode()
        writer.write(hdr + payload)
        await writer.drain()
    except Exception as e:
        try:
            err = json.dumps({"ok": False, "error": "http: " + str(e)[:200]}).encode()
            hdr = (f"HTTP/1.1 200 OK\r\nContent-Type: application/json\r\n"
                   f"Content-Length: {len(err)}\r\n\r\n").encode()
            writer.write(hdr + err)
            await writer.drain()
        except Exception:
            pass
    finally:
        try:
            writer.close()
        except Exception:
            pass


async def main():
    ws_server = await websockets.serve(ws_handler, "127.0.0.1", 10087,
                                       ping_interval=20, ping_timeout=20,
                                       max_size=16 * 1024 * 1024)
    http_server = await asyncio.start_server(http_handler, "127.0.0.1", 10088)
    log("daemon UP: WS=10087 (extension), HTTP=10088 (curl)")
    log("waiting for Browser Control extension to connect to ws://127.0.0.1:10087/ws ...")
    async with ws_server:
        await asyncio.gather(ws_server.wait_closed(), http_server.serve_forever())


if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        log("interrupted, exit")
```
