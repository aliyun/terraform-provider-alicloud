# jarvis

可被 Claude 完全接替日常工作的自包含起点母版。

## 启动

```bash
bridge/run.sh start          # 数字员工常驻：钉钉收消息 → 跑 jarvis → 流式回卡
claude                       # 当面用：进目录起会话，CLAUDE.md 接管 scan→plan→triage
```

一次性 setup：`bash bootstrap/install.sh`（装依赖）→ `cp bootstrap/.env.example bootstrap/.env` 填密钥 → `bin/a1id login jarvis`（只登 jarvis；钉钉桥另填 `bridge/jarvis.env`）。

## triage loop

人确认凭证后，Claude 自动执行以下流程：`bootstrap/scan.sh` 拉取当前 Aone 入箱清单 → `bootstrap/plan.sh` 生成执行计划（含动作、置信度、auto/stop、不可逆点）→ supervised 模式下等你逐条授权 → 按 `loops/aone-triage.md` 处理授权项，全链跑到预发/CR。**正式发布永停**，须人工介入。
