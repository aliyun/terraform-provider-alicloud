# cap-claim-ledger-owner-scoping — 认领台账加实例归属 + Stop 闸门收窄到"自己"

- **缺口类型**:流程缺陷(整机共享状态无归属 → 跨会话串味)
- **触发场景**:处理工单 83899246 收尾时,本会话的 Stop 钩子(`wrap-check.sh`)报出 **83900538 未收尾** —— 而 83900538 是**另一个 jarvis 流程/会话**认领并在跑的 in-flight 单(ACube AccTest 超时修复),它用 `wrap.sh sync` 同步过进展(sync 不写 run_done),本会话被迫替它背锅。
- **置信度**:high(根因已在源码层定位,复现机制清晰)

---

## 一、根因

`.my-day/claims-*.json` 认领台账是**整机一份、扁平、无归属**的列表,条目只有 `{"id","done"}`;`bootstrap/wrap-check.sh`(Stop 钩子)扫**整机所有** `claims-*.json`(跨日期)+ `touched-*.json`,把每一条 `done:false 且无 run_done` 的认领一律 `exit 2` 阻断,**不区分是哪个实例/会话开的**。

claim 留痕里的 `U-4L600YWP-2032` 是 `scutil --get ComputerName`(**整机一个**),更区分不了会话。

`coord.sh` 的实例 id 是 `hostname-$$`(带 PID,**每进程唯一**),写在 `.my-day/instances/`,`coord.sh dead <id>` 可判存活(hb 文件龄 > TTL 或 PID 不在 → 视为死)—— 但**认领台账没有把这个归属记下来**,Stop 闸门也没用它。

### 两种串味

| 场景 | 后果 |
|------|------|
| **真并发**:A 正做单 X(已认领、未 run_done),B 干完单 Y 触发 Stop | B 被 A 的在途单卡住,B 无过错 |
| **顺序遗留**(本次):早先会话 sync 过但没 run_done 就退了 | 下一个会话 Stop 替它背锅,甚至(如本次)被迫给别人的单补 run_done,污染其审计线 |

外加隐患:多个进程并发 read-modify-write 同一份 `claims-*.json`,`_ledger_upsert` 虽单次写原子(mktemp+mv),但跨进程仍可能丢更新(lost update)。

### 为什么现在没天天爆

现架构默认**一机一个活跃 triage 实例**:`coord.sh` 是**顺序接管**模型(`list-orphans` 只挑 owner 进程已死的 task 让新实例 `adopt`),dispatch(headless)实例"只心跳不 adopt"、由 bridge 串行喂单。假设一旦被打破(本机 `.my-day/instances/` 有 ~10 个 `jarvis.local-<PID>`),缺陷即咬人。

---

## 二、建议补丁(分期,保守、不回归)

### Phase 1(本单实现):归属 + Stop 收窄
1. **`claim.sh _ledger_upsert` 记 owner**:条目加 `owner`,值取 `${COORD_ID:-}`(coord.sh register 返回、loop 导出的实例 id;交互/无 coord 会话则为空)。
   - INSERT 时写 owner;UPDATE(release/finish)时 **coalesce 保留原 owner**(`(.owner // "")==""` 才回填),避免异实例收尾时覆盖原认领者。
2. **`wrap-check.sh` 按 owner 收窄**:self = `${COORD_ID:-}`。对每条 open(`done:false`)/ touched 的 id 取其台账 owner:
   - `owner` **非空且 ≠ self** → **不 block**,仅打 WARN(可见性),交给 `reconcile` 兜底;
   - `owner == self` **或 owner 为空(遗留/交互)** → 仍要求 run_done,缺则 block(**不回归**:自己开的单、以及无归属的历史条目照旧问责)。
   - touched 无 owner 字段 → 查台账 owner 套同规则;查不到 → 当遗留 → block。

净效果:一个会话的 Stop **只硬卡自己(本实例)开的单和无归属历史条目,绝不再卡别的具名实例的单**。修掉带 COORD_ID 的并发/顺序串味主场景。

### Phase 2(后续单):并发写锁 + 交互会话稳定归属
- `_ledger_upsert` 加 mkdir-lock(macOS 无 flock)消除跨进程丢更新。
- 交互(无 coord register)会话缺稳定 owner 的问题:考虑 claim.sh 无 `COORD_ID` 时懒注册一个 coord 实例或从稳定来源(tty/父 shell)派生 owner,让 Phase 1 的收窄对交互会话也生效。
- `reconcile stale/orphan` 增补:显式认领"异实例死后遗留的 open claim"归口,与 Stop 的 WARN 联动。

---

## 三、验收

- `test/claim_test.sh`(owner 字段写入/coalesce)+ `test/wrap_check_test.sh`(self / 异活实例 / 异死实例 / 空 owner 四种归属的 block-vs-skip 矩阵)全绿。
- 无回归:自己开的单、无归属历史条目仍被 Stop 拦截。

## 四、状态

- 入队:cap 缺口 + Phase 1 补丁草案 → MR 待人工评审。
- 红线:**绝不自动合 master**(self-improve.md)。
