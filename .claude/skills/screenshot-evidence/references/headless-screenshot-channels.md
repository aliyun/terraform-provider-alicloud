# headless 截图通道与 Playwright MCP 缺失根因

> 本文档对应工单：headless 缺少 Playwright MCP 导致无法截图。说明 headless runtime 与
> 交互会话的 MCP/tool 注入差异、缺失根因，以及仓库内可控的降级截图通道。

## 一、现象

`screenshot-evidence` skill 的「前置条件」长期写为 **Playwright MCP 已连接
（`mcp__playwright__*` 工具可用）**。在交互会话里这条成立；但在 bridge 拉起的 headless
会话里，`mcp__playwright__*` 工具族**根本不存在**，于是需要浏览器截图证据的流程
（尤其 Terraform 三层查证：OpenAPI / CloudSpec/ACube / Provider）无法产出截图，
任务卡在截图环节或被迫静默跳过——后者更糟，等于漏做证据。

## 二、根因：headless 不注入 Playwright MCP

headless 会话由 `bridge/jarvis_execution_runtime.py::run_claude_buffered` 拉起，命令由
`jarvis_cmd()` 拼装（`bridge/jarvis_execution_runtime.py:393-412`）：

```
claude --settings <provider-settings.json> --permission-mode bypassPermissions \
       -p <text> --output-format json --session-id <id>
```

再经 `bootstrap/jarvis-interactive-worker.py exec-headless`（`os.execvpe` 透传，
`bridge/jarvis-interactive-worker.py:2831-2849`）执行。**MCP server 只可能来自 settings 文件**，
而 provider settings 文件（默认 `~/.claude/idea_settings.json`，或 `JARVIS_SETTINGS` /
`JARVIS_SETTINGS_TF` 指向的文件）只载**模型 provider env**（`ANTHROPIC_BASE_URL` /
`ANTHROPIC_MODEL` / `ANTHROPIC_AUTH_TOKEN`），**不声明任何 `mcpServers`**；仓内也没有
`.mcp.json`，`exec-headless` 只是 `os.execvpe` 透传、不注入任何 MCP。

交互会话里的 Playwright MCP 是**仓库主人本地交互环境**的产物（装在该机器的 Claude
desktop/CLI 配置里），**不属于**仓库可控的 headless 启动链。换台 worker、或换一个没装
Playwright MCP 的交互环境，headless 一律拿不到 `mcp__playwright__*`。这就是「缺失」的根因：
不是配置漏了一行，而是 headless 启动链结构与交互会话不同，从来就没有 MCP 注入这一层。

## 三、修复策略：仓库内可控、可探测、可降级

不依赖交互态 MCP，改为仓库内可控的浏览器通道，按优先级探测：

| 优先级 | 通道 | 实现 | 适用 |
|--------|------|------|------|
| 1 | `playwright_python` | `playwright` Python 绑定 + 本地 chromium | 全页 + 元素级截图，SPA 等渲染稳定 |
| 2 | `chrome_binary` | headless Chrome/Chromium + DevTools Protocol | 真全页 + 目标文本元素截图；macOS app / Linux PATH / `JARVIS_CHROME_BIN` |
| — | 无 | `missing_capability` 收口 | 任务以可诊断原因收口，不静默跳过 |

- 探测与捕获实现在 `bridge/jarvis_screenshot.py`（单测 `bridge/test_jarvis_screenshot.py`）。
- skill 侧封装为 `.claude/skills/screenshot-evidence/scripts/capture.sh`，从任意 cwd 可调用，
  内部 `python3 -m bridge.jarvis_screenshot probe|capture`，并把仓库根插到 `PYTHONPATH` 最前，
  防止被继承的 `PYTHONPATH`（如 harness 注入的主仓）shadow 掉本仓模块。
- 退出码：`0` 通道可用/捕获成功；`3` `missing_capability`（无可用通道）；`1` `capture_error`
  （通道在但页面捕获失败）。`3` 与 preflight 的 `missing_capability` 口径一致。

## 四、与 skill 的集成

- **Step 0 能力探测**：截图前先 `capture.sh probe`；命中通道记名进 Step 2，命中
  `missing_capability` 则把该层 manifest 写 `n-a` + 原因，可继续只读文字查证，但 finalizer
  必须 `blocked / missing_capability`，不得 done。`validate-manifest.py` 接受带原因的
  `n-a` 行只代表结构完整。**绝不静默跳过**。
- **Step 2 截图**：交互会话仍可走 `mcp__playwright__*` 元素截图；headless / 无 MCP 时走
  `capture.sh capture`；`--full-page` 走 CDP 真全页，`--text` 截目标字段元素。
- 不影响交互模式与本机已有 Browser/Chrome 能力：交互态 Playwright MCP 与本机 Chrome 各自独立，
  headless 只新增仓库内通道，不改交互启动链。

## 五、worker 配置（可选，截图为降级能力，非硬门）

截图是降级能力，preflight 对其 WARN 而非 FAIL（见 `bootstrap/verify.sh` 的 screenshot 检查）。
需要 headless 实拍截图的 worker，二选一装一个通道即可：

```bash
# 通道 1：Playwright Python（推荐，支持元素级 + SPA 等待）
python3 -m pip install playwright
python3 -m playwright install chromium

# 通道 2：headless Chrome/Chromium 二进制（CDP 真全页 + 目标文本元素截图）
#   Linux: yum install -y chromium  或 apt install chromium-browser
#   macOS: /Applications/Google Chrome.app（本机已有）
#   也可 export JARVIS_CHROME_BIN=/path/to/chrome 显式指定
```

两通道都没装的 worker，截图诉求以 `blocked / missing_capability` 收口并被人工/上游看见，
不会静默漏做或误标完成。

## 六、验收映射

| 工单验收项 | 落地 |
|------------|------|
| 不预装/不注入 Playwright MCP 的 headless 仍能打开测试页并产出有效截图 | `chrome_binary` 通过 CDP 实拍 `example.com` → 真全页有效 PNG（单测 + 本机自测） |
| screenshot-evidence 最小链路在 headless 完成截图 + manifest | `capture.sh capture` 产 PNG，manifest `n-a` 行带原因被 `validate-manifest.py` 接受 |
| 能力完全不可用→可诊断 `missing_capability` 收口，不静默跳过 | `probe` exit 3 + `missing_capability:` 原因，skill 写 `n-a` 不删层，finalizer blocked |
| 自动化回归用例覆盖「无 Playwright MCP」场景 | `bridge/test_jarvis_screenshot.py`：无通道→exit 3、优先级、捕获派发、有效 PNG 产出 |
| 不影响交互模式与本机 Browser/Chrome | 不改交互启动链；通道独立探测，交互 Playwright MCP 不受影响 |
