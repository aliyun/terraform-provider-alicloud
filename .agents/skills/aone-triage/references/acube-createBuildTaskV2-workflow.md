# 临钧路由(生成器产出) · acube V2 完整流程

> 从 `tf-customer-request-routing.md` 分支 D-临钧段抽出的完整实现细节。
> 主 skill 只写"生成器产出走 acube V2 接口,jarvis 不手动建单",本文件放
> 详细 bash 步骤 + 关键纪律,单点维护避免主 skill 膨胀。

## 背景

生成器产出资源交给 acube 的 `TerraformVendorBuildTaskOpenapiController#createBuildTaskV2` 接口——接口内部**自动**在 terraform-alicloud (528766) 建关联单、指派临钧(429768)、触发生成/PR 工作流,jarvis 只负责查回 aoneId 并做源单关联+指派。**严禁**同时走 `a1 workitem create` 手动建单流程,否则双单污染临钧队列。

## 服务端接口(邻仓 `a-cube-aliyun-com`)

- `POST /api/v1/terraform_vendor_build/createBuildTaskV2` — body `TerraformVendorBuildTaskDTO`,返回 `ResultDTO<Long>` (taskId,同步返回)
- `GET  /api/v1/terraform_vendor_build/queryAoneByTaskId?taskId={taskId}` — 返回 `{taskId, aoneId, aoneUrl}`,aoneId 异步产生(acube 内部建单完成后回写),需轮询

## 完整 bash 流程

```bash
# 0. 拿 jarvis 工号(acube 侧 workId/workName 用当前 a1 身份,便于事后追溯)
#    jarvis 默认身份的 Emp ID 是 WORKER_ 前缀长 id(非 5-11 位数字),同样合法——
#    acube workId 接受任意字符串
jarvis_empid=$(bin/a1id -- auth whoami 2>/dev/null | python3 -c '
import sys,re
raw=sys.stdin.read()
# 兼容 whoami 输出的多种格式:WB 外包工号 / WORKER_ agent 身份 / 5-11 位数字工号
m=re.search(r"\b(WB\d+|WORKER_\d+|\d{5,11})\b", raw)
print(m.group(1) if m else "")')
[ -z "$jarvis_empid" ] && echo "jarvis 未登录 a1(bin/a1id login jarvis),阻断" && exit 1

# 1. 触发 build 任务(acube 自动建单+指派临钧+跑生成器)
#    必填字段: namespace / resourceTypeCode / resourceTypeVersion / osType / flowType / workId / workName
#    resourceTypeVersion 走"生成器产出待跑"场景填 0.0.0(acube 会跑首版生成)
task_id=$(curl -s -X POST "https://acube.aliyun-inc.com/api/v1/terraform_vendor_build/createBuildTaskV2" \
  -H "Content-Type: application/json" -H "accept: */*" \
  -d "{
    \"namespace\":\"<product>\",
    \"resourceTypeCode\":\"<PascalCase Resource>\",
    \"resourceTypeVersion\":\"0.0.0\",
    \"osType\":\"Linux\",
    \"flowType\":\"ACubeRelease\",
    \"workId\":\"$jarvis_empid\",
    \"workName\":\"jarvis\"
  }" | python3 -c '
import json,sys
d=json.load(sys.stdin)
if d.get("code")!="SUCCESS":
  sys.stderr.write(f"createBuildTaskV2 failed: {d.get(\"code\")} {d.get(\"message\")}\n"); sys.exit(1)
print(d.get("data"))')
[ -z "$task_id" ] && echo "acube createBuildTaskV2 未返回 taskId,阻断" && exit 1
echo "taskId=$task_id"

# 2. 轮询查 aoneId(taskId 立返,aoneId 需等 acube 异步建单完成;60s 内应有值)
NEW_ID=""
for i in 1 2 3 4 5 6; do
  NEW_ID=$(curl -s "https://acube.aliyun-inc.com/api/v1/terraform_vendor_build/queryAoneByTaskId?taskId=${task_id}" \
    -H "accept: */*" | python3 -c '
import json,sys
d=json.load(sys.stdin); data=(d.get("data") or {})
print(data.get("aoneId") or "")')
  [ -n "$NEW_ID" ] && break
  sleep 10
done
if [ -z "$NEW_ID" ]; then
  echo "acube 60s 内未返回 aoneId,升级 escalation/(不要回退到手动 a1 workitem create,可能双建)"
  bootstrap/log.sh escalate <源工单ID> "acube build task $task_id 60s 内未返回 aoneId,人工排查"
  exit 1
fi
echo "临钧关联单 aoneId=$NEW_ID"

# 3. 关联到源客户单(aone 自动双向,单次 relation add 即建 A↔B)
bin/a1id -- project workitem relation add <源工单ID> relate:$NEW_ID

# 4. 源工单同步指派临钧 + 状态(评论 @临钧 走 Step 4 模板 B)
bin/a1id -- project workitem update <源工单ID> --assignee 429768
bin/a1id -- project workitem update <源工单ID> --status 问题解决中
```

## 环境

- 正式走 `acube.aliyun-inc.com`,预发把域名换成 `pre-acube.aliyun-inc.com`(路径/参数/返回结构一致)
- **pre-acube 不是无副作用沙箱**:`createBuildTaskV2` 在预发同样写真 Aone(528766 建单并指派临钧);连通性测试只打只读 `queryAoneByTaskId`
- `/api/v1/**` 免鉴权,内网 DNS(需办公网/VPN)

## 关键纪律

- acube 自动建单+指派+触发工作流是**原子动作**,jarvis 只做"查 aoneId + 关联源单"善后
- 60s 内没查到 aoneId → 直接升级 escalation,**禁**回退手动 `a1 workitem create`,双建会污染临钧研发队列
- workId/workName 填当前 jarvis 身份工号,acube 侧任务日志能追到调用方
