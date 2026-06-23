# 并发认领:池级 scan + tag 认领 + 僵尸清扫

> 设计 v0 · 2026-06-22 · 待评审

## 目标
多机各自 clone 的 Claude 并行处理 5 池,不抢不重,挂掉的认领能被回收。

## 已验事实
- tag 过滤须 `--project`(池级资源)→ scan 必须逐池
- 无 untag 命令:领=tag jarvis-claimed,完=retag jarvis-done(两标签已建)
- comment create/list 可用 → 认领带时间戳

## 机制
- **认领** `claim <id> <proj>`:`update --tag jarvis-claimed` + 评论 `jarvis-claim <host> <UTC>`;回读确认是己
- **释放** `release <id> <proj>`:`update --tag jarvis-done`
- **池级 scan**:循环 pools.json 项目,`list --project P --filter "NOT tag=jarvis-claimed"` 合并
- **僵尸清扫** `sweep`:逐池 `list --tag jarvis-claimed`,读 claim 评论时间,超 TTL(默认 45min)→ 入 escalation/ 待人回收(不自动改挂,避免误抢活着的)

## 改动
1. config 加 `claim.done_tag=jarvis-done`、`claim.ttl_min=45`
2. bootstrap/claim.sh:claim/release,回读校验
3. scan.sh 改池级合并
4. bootstrap/sweep.sh:僵尸→escalation/
5. loops/autonomy 串入:领→done,sweep 定期跑

## 验收
两 clone 不抢同单;领→done 闭环;杀一个处理中→sweep 45min 判僵尸入队;scan 跳已领。

## YAGNI
不做:自动回收僵尸、心跳续约、跨池并行调用。
