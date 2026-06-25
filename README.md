# jarvis

可被 Claude 完全接替日常工作的自包含起点母版。

## 配置（一次性）

```bash
bash bootstrap/install.sh                       # 装依赖（a1/aliyun/cloudspec/gh）
cp bootstrap/.env.example bootstrap/.env        # 填 GH_TOKEN + 阿里云 AK
bin/a1id login jarvis && bin/a1id login chenyi  # 两个身份各登一次
```

钉钉数字员工还需 `cp bridge/jarvis.env.example bridge/jarvis.env`，填钉钉 appKey/appSecret/卡片模板 id。

## 启动

**当面用** —— 进目录起会话，CLAUDE.md 接管 scan→plan→triage：

```bash
claude
```

**数字员工** —— 常驻钉钉，收消息 → 跑 jarvis → 流式回卡：

```bash
bridge/run.sh start
```

## triage loop

人确认凭证后，Claude 自动执行以下流程：`bootstrap/scan.sh` 拉取当前 Aone 入箱清单 → `bootstrap/plan.sh` 生成执行计划（含动作、置信度、auto/stop、不可逆点）→ supervised 模式下等你逐条授权 → 按 `loops/aone-triage.md` 处理授权项，全链跑到预发/CR。**正式发布永停**，须人工介入。
