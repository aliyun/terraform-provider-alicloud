# cap-a1id-identity-guard

> **[归档]** 本 cap 针对的 v1 `bin/a1id`(复制覆盖 live 的单文件布局)已被 v2 整体重写(每身份独立 `A1_CONFIG_DIR`,commit dd3d2f5):P0 账号映射 + login whoami 校验已内置;P1 verify 诉求由 `a1id ready/status` + `verify.sh` 日检等价覆盖;P2/P3 随架构变更失效(无 store→live 拷贝链路)。归档保留正文作历史缺口记录。

## 缺口类型

`bin/a1id` 缺身份一致性护栏,SSO 时静默把任意 BUC 会话存进任意身份槽位,导致 jarvis 公用身份被个人身份凭据污染。

## 阻塞任务

`/aone-triage` 开局身份验证阶段发现:

- 本机 `~/.config/a1/identities/` 只存了 `jarvis.auth.yaml`,但内容是 `guozai.gzl`(过载本人)的凭据(sha256 与浏览器 live BUC 会话完全一致,`diff` 为空)。
- 意味着 jarvis 公用身份(BUC Account = `WORKER_1782379562571`,Aone 数字员工)在本机**根本没登进过**——所谓"jarvis 已登录"是假象。
- 直接跑 `/aone-triage` 会以 guozai 名义扫入箱 + claim + 回复,踩 CLAUDE.md 工作纪律 #6 红线。

已做的即时修复:备份 `~/.config/a1/identities/` 到 `.my-day/a1id-store-backup-<ts>/`,删除污染的 `jarvis.auth.yaml`,清空 `.active`。下次任何 a1id 命令(`use jarvis` / `--`)都会 `die` 提示"身份 'jarvis' 未登录",强制走 `bin/a1id login jarvis`。

## 当前发现

### Root cause

`bin/a1id` 55–59 行 `login` 分支:

```bash
login)
  id="${1:-}"; valid "$id" || die "用法: a1id login <jarvis|chenyi|guozai|linjun>"
  "$A1" auth login --buc
  cp "$LIVE" "$(store_for "$id")"; echo "$id" > "$ACTIVE"   # 登完落盘到该身份 store
  echo "a1id: '$id' 登录已保存"; "$A1" auth whoami ;;
```

- `a1 auth login --buc` 只走浏览器 BUC SSO,继承当前 BUC 会话——**没有** `--account` 之类的目标账号锁定参数。
- 命令行标签(`login jarvis`)只影响本地文件命名,与实际登入的账号解耦。
- `cp` 前无校验,回执 `whoami` 印在事后当装饰,不阻断落盘。
- 结果:浏览器 BUC 登的是谁,`<label>.auth.yaml` 就存谁的凭据。

### 复现路径

1. 浏览器 BUC 会话为个人账号(guozai)。
2. `bin/a1id login jarvis`。
3. 弹网页 SSO,自动通过(已登)。
4. `jarvis.auth.yaml` 被写入 guozai 凭据,`a1id: 'jarvis' 登录已保存` 提示成功。
5. 之后 `bin/a1id -- <args>` 全部以 guozai 名义跑,`whoami` 无从察觉。

本次踩坑时间线:`~/.config/a1/identities/jarvis.auth.yaml` mtime `2026-07-02 07:20`,即本机上一次 `a1id login jarvis` 的时点。

### 二级影响

- `use` / `as` / `--` 分支通过 `activate()`(46–51 行)只做 store → live 拷贝,不校验 live 内容与 label 是否一致——**污染一旦落盘,后续无法自愈**。
- `who` 分支(62 行)透传 `a1 auth whoami`,能看出真实账号,但**没人会主动跑**,污染继续。
- `chenyi` / `guozai` / `linjun` 三个个人身份同样脆弱,只是 jarvis 是默认身份+使用最频繁,污染代价最大。

## 建议补丁

### P0 — 已知账号映射 + login 后校验(必做)

在 `bin/a1id` 里维护 `id → 期望 BUC 账号名` 映射:

```bash
expected_account_for(){
  case "$1" in
    jarvis) echo "WORKER_1782379562571" ;;   # BUC Account:字段,Aone 数字员工 ID
    chenyi) echo "chenhanzhang.chz" ;;
    guozai) echo "guozai.gzl" ;;
    linjun) echo "lichaolin.lcl" ;;
    *) return 1 ;;
  esac
}
```

改造 `login` 分支:先把 LIVE 备份,登完抽 `a1 auth whoami` 的 `Account:` 与 expected 对比,不匹配立即还原 LIVE 并 die,给清晰恢复指令。

```bash
login)
  id="${1:-}"; valid "$id" || die "用法: a1id login <jarvis|chenyi|guozai|linjun>"
  expected="$(expected_account_for "$id")"
  # 备份 LIVE 以便回滚
  live_bak=""; [ -f "$LIVE" ] && live_bak="$(mktemp)" && cp "$LIVE" "$live_bak"
  "$A1" auth login --buc
  got=$("$A1" auth whoami 2>/dev/null | awk '/Account:/{print $2}')
  if [ "$got" != "$expected" ]; then
    [ -n "$live_bak" ] && cp "$live_bak" "$LIVE" && rm -f "$live_bak"
    die "身份不匹配:请求登入 '$id'(期望账号 $expected),实际浏览器 BUC 会话是 '$got'。
     修:1) 浏览器登出 BUC (https://buc.alibaba-inc.com/);
        2) 或用无痕/其它浏览器登 '$expected';
        3) 重跑 'bin/a1id login $id'。
     现已还原 live auth.yaml,原身份未受影响。"
  fi
  [ -n "$live_bak" ] && rm -f "$live_bak"
  cp "$LIVE" "$(store_for "$id")"; echo "$id" > "$ACTIVE"
  echo "a1id: '$id'($got) 登录已保存"; "$A1" auth whoami ;;
```

### P1 — `verify` 子命令(供 preflight 集成)

新增 `bin/a1id verify` 遍历 `$STORE/*.auth.yaml`,对每个身份短暂 activate + `a1 auth whoami`,校验实际账号 = expected;不一致列表输出退非零。可挂到 `bootstrap/preflight.sh`,一年跑一次也能兜住存量污染。

### P2 — `activate` 加一次性校验(可选)

`activate()` 里在 `cp store → live` 之后 `a1 auth whoami` 快速 check,不匹配 warn(不 die,因为切换本身不写盘,只是提醒)。低价值高噪声,存疑,先不做。

### P3 — 文档修 stale 引用

`bin/a1id` 头部注释 10 / 14 / 17 行"纪律见 CLAUDE.md 第8条",当前是**第 6 条**(numbering 已收敛)。顺手改。

## 置信度

high:root cause 通过读 `bin/a1id:55-59` + 本机实际复现(`diff ~/.config/a1/auth.yaml ~/.config/a1/identities/jarvis.auth.yaml` 为空)双层确认;补丁方案可在本仓 mock `a1` 二进制跑测试用例。

## 关联

- 待建 Aone 需求单(api_toolkit 池 2100304 或 tf_provider 池 528766):记录问题+护栏计划。
- 后续 MR:实施 P0/P1 补丁,worktree 分支名 `worktree-a1id-identity-guard-impl`。
- 相关工作纪律:CLAUDE.md #6 身份纪律。
