# 五阶段流程

## 0. 自举 (Bootstrap)
- 安装依赖 (install)
- 验证环境 (verify)

## 1. 扫描 (Scan)
- 执行 `scan.sh`
- 并行扫描 5 个池
- 带优先级 (priority) 和标签 (tag)

## 2. 去重 + 分诊门 (Deduplication + Triage Gate)
- `plan.sh` 去重
- Claude 按优先级分诊
- **[人工点]** Supervised 授权

## 3. 逐条 Triage (Per-Item Triage)
- Claim 工作项
- `aone-triage` 查证
- Autonomy 判定
- 执行: `run_done` 或 `escalate`
- Release 版本

## 4. 硬门 (Hard Gate)
- **[人工点]** `release_prod` 永停，需人工审批
- 正式发布

## 5. 收工 (Cleanup)
- `sweep` 僵尸任务
- Escalation 处理
- Runs 审计

---

## 人工审核点
1. **阶段 2**: 分诊授权 (Supervised authorization before distribution)
2. **阶段 4**: 正式发布 (Manual approval for production release via release_prod)
