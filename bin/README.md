# bin/a1id v2 —— a1 多身份并发切换器

## 布局

每个身份独立 config dir（`identities/<label>/`，内含完整 `auth.yaml`）；`as`/`--` 只把 `A1_CONFIG_DIR` 指到目标身份 dir 再 exec a1——不改 live（`~/.config/a1/auth.yaml`）、不改 `.active`（仅 `a1id use` 更新，脚本链路不读），**天然并发安全**。

启动时自动做**幂等**收编:发现旧单文件布局 `identities/<label>.auth.yaml` 而目录 `identities/<label>/` 尚未建 → `mkdir + cp`,旧文件**保留**(便于回退)。首跑 live 收编(仅 jarvis 且新旧布局全空时)落到 `identities/jarvis/auth.yaml`。

## 八个身份

| label | 期望 BUC 账号 | 别名 | 角色 |
|-------|---------------|------|------|
| `jarvis`       | `WORKER_1782379562571` | —  | 默认;Aone 数字员工;编排层(主会话)用 |
| `terraform-pd` | `WORKER_1783582374386` | pd | 产品数字人:分诊/查证/回复,对应 agent `terraform-pd.md` |
| `terraform-rd` | `WORKER_1783582458263` | rd | 研发数字人:开发/PR/CR 评审,对应 agent `terraform-rd.md` |
| `terraform-qa` | `WORKER_1783582593461` | qa | 质量数字人:AccTest 验证/验收,对应 agent `terraform-qa.md` |
| `chenyi` | `chenhanzhang.chz` | — | 陈汉璋(工号 320687);Jarvis 不得擅用,仅当面授权时 |
| `guozai` | `guozai.gzl`       | — | 郭子龙(工号 484483);Jarvis 不得擅用,仅当面授权时 |
| `linjun` | `lichaolin.lcl`    | — | 李超林(工号 429768);Jarvis 不得擅用,仅当面授权时 |
| `shanye` | `shanye.xzq`       | — | 杉也/徐茈琦(工号 414322);Jarvis 不得擅用,仅当面授权时 |

## 命令面

| 命令 | 作用 |
|------|------|
| `a1id login <id>`               | 交互 BUC SSO 登录,落盘 `identities/<id>/auth.yaml`;whoami 与期望账号不匹配则清盘并 die |
| `a1id use <id>`                 | 拷贝 `identities/<id>/auth.yaml` 到 live(仅影响人工直接跑 a1 的 live 会话;脚本链路不需要) |
| `a1id status`                   | 显示默认身份 / live active / A1ID_ROOT / 八身份登录表 |
| `a1id who [id]`                 | `a1 auth whoami`;缺省=默认身份 dir,指定=该身份 dir |
| `a1id ready <id>`               | 脚本探测:已登录退 0,否则退 1 |
| `a1id as <id> -- <a1 args...>`  | 以指定身份跑一条(严格;未登录直接 die,不回退) |
| `a1id -- <a1 args...>`          | 以默认身份跑(受 `JARVIS_A1_IDENTITY` 影响) |

## 环境变量

| 变量 | 作用 |
|------|------|
| `A1ID_ROOT`              | a1id 根目录(默认 `~/.config/a1`);测试用它隔离 |
| `A1_BIN`                 | a1 二进制路径(默认 `a1`);测试打桩用 |
| `JARVIS_A1_IDENTITY`     | 默认身份覆盖;canon+valid 校验;未登录时非 strict + jarvis 已登录 → 警告并回退 jarvis;个人身份 → stderr 打纪律告警(不阻断) |
| `JARVIS_A1_STRICT`       | 严格模式(=1):默认身份未登录时直接 die,不回退 |

## 一次性设置

```bash
bin/a1id login jarvis         # 数字员工(编排层默认)
bin/a1id login terraform-pd   # 产品数字人
bin/a1id login terraform-rd   # 研发数字人
bin/a1id login terraform-qa   # 质量数字人
# 个人身份按需登录:
# bin/a1id login chenyi
# bin/a1id login guozai
# bin/a1id login linjun
# bin/a1id login shanye
```

登录时浏览器 BUC 会话必须是对应期望账号;否则 `a1id` 会清掉这次污染的 auth.yaml 并给出修复步骤
(浏览器登出 / 无痕重登 / 重跑)。live 全程未动。

## 数字人 ↔ Agent 对应

| Aone 身份 | Claude Agent(`.claude/agents/`) | 主职 |
|-----------|--------------------------------|------|
| `terraform-pd` | `terraform-pd.md` | 产品:需求分析/工单分诊/三层查证 |
| `terraform-rd` | `terraform-rd.md` | 研发:代码开发/PR 提交/CR 评审 |
| `terraform-qa` | `terraform-qa.md` | 质量:AccTest 验证/验收确认 |

Agent 内推荐用法:

```bash
# 单条:以角色身份跑一条 a1
bin/a1id as terraform-pd -- project workitem comment create <id> -m "..."

# 整链路由:让 wrap.sh / claim.sh 等自动走该身份
JARVIS_A1_IDENTITY=terraform-rd bash bootstrap/wrap.sh sync <id> "MR: ..."

# 开工探测:未登录自动回退 jarvis,返回结果标 identity_fallback: jarvis
bin/a1id ready terraform-qa || echo "未登录,已回退 jarvis(标 identity_fallback: jarvis)"
```

## 并发原理(重点)

v2 每个身份一个独立 `A1_CONFIG_DIR`,`a1id as pd -- ...` 与 `a1id as rd -- ...` 并发跑各自 exec 独立 a1 子进程,读写各自 dir 内的 `auth.yaml`,互不干扰。live `~/.config/a1/auth.yaml` 只由 `a1id login`/`a1id use` 影响(人工态),脚本链路(`bootstrap/wrap.sh`、`bootstrap/claim.sh`、`bootstrap/scan.sh` 等)全部通过 `bin/a1id --` 或 `bin/a1id as ...` 走,不动 live。
