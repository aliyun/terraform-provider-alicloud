# plugin-dev 插件开发 loop

> 在 alibabacloud-agent-toolkit 里新增/改一个 Claude Code 插件（MCP+skills+hooks）的标准流程：调研→照模子起骨架→写 skill+参考→本地门禁→真凭证实测→推分支等评审。
> 与 `loops/adhoc-intake.md` 衔接：意图=开发且落点是 agent-toolkit 时走本 loop。脏活交子代理，主 Agent 只编排（CLAUDE.md 工作纪律）。

---

## 一、触发

| 方式 | 说明 |
| --- | --- |
| 手动 | 「在 toolkit 里做个管 X 的插件」「按方案A出个 Plugin 版本」 |
| 链接 | 钉钉设计文档 / Aone 需求（建插件类） |
| 缺口 | 现有插件不覆盖某云产品场景 |

落点工作区 `agent_toolkit`（config/workspaces.json）；改文件先开 worktree。

---

## 二、调研（先调研，禁臆造）

1. 看模子：读 plugins/alibabacloud-core 的 manifest/.mcp.json/skills/hooks，定为照搬模板。
2. 看接口：目标云产品的 OpenAPI 真表 — `api.aliyun.com/meta/v1/products/<P>/versions/<ver>/api-docs.json` 拉全量算子，别信记忆。
3. 多代理并行 fan-out（插件骨架/skills/hooks/接口/方案映射），合成 design brief。
4. 凡接口/动词，**以 meta + 实测为准**，文档里所有名都待验证。

---

## 三、起骨架（照 core 模子）

`plugins/<name>/`：`.claude-plugin/plugin.json` + `.codex-plugin/plugin.json`（name/version/description/author/interface 全）+ `.mcp.json`（uvx alibabacloud.mcp-proxy，`--safety-policy <product>:*=allow,*=deny`）+ `skills/<name>/` + `hooks/`**逐字节 cp core**（CI verify-hooks 校验）+ README。两处 marketplace 追加：`.claude-plugin/marketplace.json`、`.agents/plugins/marketplace.json`。

约束（validate.py）：codex 清单必含 5 键；skill name=目录、kebab、desc≥20；mcp stdio 必有 command。

---

## 四、写 skill + 参考

skill 名=目录名；frontmatter desc 含中英 triggers + allowed-tools 全限定 `mcp__plugin_<name>_<server>__...`。重逻辑落 `references/*.md`，镜像 spec-ops iac-service-api.md 风格。命名先调研同场景 aws/azure/gcp 约定再定。

---

## 五、门禁（全绿才推）

```bash
python3 tools/validate.py --plugin <name> && python3 tools/validate.py
npx markdownlint-cli2 'plugins/<name>/**/*.md'
bash tools/dev-hooks/verify-hooks.sh
```

---

## 六、实测（真凭证，跑完整旅程）

profile=TerraformUT 走端到端，记真实体验+卡点，回写文档与 skill：
- 动词只认 kebab；写动词带 `--client-token` UUID；plan/apply 闸 + destroy 清理，账号不留残留。
- 测出的纠偏直接固化进 reference/skill（标 verified）。

---

## 七、Done

| 结果 | 说明 |
| --- | --- |
| 完成 | 门禁绿+实测通过 → push 分支，`run_done` 入 runs/，关 Aone |
| escalation | 缺接口/低置信 → escalation/ + self-improve；**正式合 marketplace=红线永停** |

---

## 八、工具链速查

| 工具 | 作用 |
| --- | --- |
| plugins/alibabacloud-core | 照搬模子 |
| api.aliyun.com/meta api-docs.json | 算子真表 |
| tools/validate.py / verify-hooks.sh | 门禁 |
| TerraformUT profile | 实测凭证 |
| config/workspaces.json `agent_toolkit` | 工作区登记 |
