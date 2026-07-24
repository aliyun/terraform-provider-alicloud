# Generator 产物后处理清单

Terraform generator-v4 有若干已知产出缺陷，跑完 Step 6 后 SOP 必须逐项过一遍。这些不是每个资源都命中，但每次都得手工 grep 一遍——**漏掉一项等到 ACC/CI 才发现，代价是重跑一轮**。

**判定原则**：这些缺陷都在**生成产物**（`alicloud/resource_alicloud_<x>.go` / `alicloud/data_source_alicloud_<x>s.go` / 两个 doc），**不是 cspec 定义问题**——不用回 CloudSpec 侧。修法：手工订正生成产物，同时把 fix pattern 汇报给 generator 维护方（`provider-resource-dev`），下次生成前修好模板。

## 检查清单

### Bug 1：Create-only 属性缺 `ForceNew: true`

**症状**：cspec 里 `@rac({operatePrivateType: ["create"]})` 或 `@readonly` + `@rac({create})` 的字段，在生成的 `*.go` schema 里只标 `Required: true` 没标 `ForceNew: true`。用户改这些字段 → 走 Update → 服务端悄悄保留旧值 → state drift。

**检查**（在 provider worktree）：

```bash
# cspec 里 create-only 属性列表（假设 CSPEC_DIR 已定位）
grep -B3 "operatePrivateType.*\"create\"" "$CSPEC_DIR/resources/<Resource>.cspec" \
  | grep -oE "^\s+[A-Z][a-zA-Z0-9]+:" | tr -d ' :'
# 换成 snake_case 对着 alicloud/resource_alicloud_<x>.go 检查是否有 ForceNew: true
```

**修法**：schema block 里加 `ForceNew: true`。同时在 doc 上把 `(Required)` 改成 `(Required, ForceNew)`（配合 Bug 4）。

### Bug 2：Read 不反向映射反向关系字段

**症状**：Create 用 `service_ids` 发 body，但 Get API 返回的是 `boundServices[*].serviceId`。生成的 Read 只 `d.Set` 了直接同名字段，`service_ids` 从来不回填 → apply 后 state 里 `service_ids=[]`，下次 plan 永远想 Update。

**检查**：

```bash
# grep 生成的 Read 函数，看是否有 boundServices→service_ids 类翻译
grep -A3 "^func.*Read" alicloud/resource_alicloud_<x>.go | grep -E "d\.Set"
# 对比 Create 里发出去但 Read 没读回的字段
```

**修法**：在 Read 末尾追加反向映射，例：

```go
serviceIds := make([]string, 0)
if boundServicesRaw, ok := objectRaw["boundServices"].([]interface{}); ok {
    for _, item := range boundServicesRaw {
        if svc, ok := item.(map[string]interface{}); ok {
            if id, ok := svc["serviceId"].(string); ok && id != "" {
                serviceIds = append(serviceIds, id)
            }
        }
    }
}
d.Set("service_ids", serviceIds)
```

### Bug 3：Datasource nested schema 缺字段

**症状**：`data_source_alicloud_<x>s.go` 的 `providers` / `items` / `<plural>` 嵌套 Elem schema 只声明部分字段，但 Read 函数里 `mapping["display_name"]` / `mapping["bound_services"]` 都在赋值——运行时 `d.Set(...)` 报 "invalid or unknown key"，datasource 直接不可用。

**检查**：

```bash
# 对比 schema 声明的 key 和 Read 里 mapping[]= 赋值的 key
sed -n '/Elem: &schema.Resource{/,/^\t*}$/p' alicloud/data_source_alicloud_<x>s.go \
  | grep -oE '"[a-z_]+":' | sort -u > /tmp/schema.keys
grep -oE 'mapping\["[a-z_]+"\]' alicloud/data_source_alicloud_<x>s.go \
  | grep -oE '"[a-z_]+"' | sort -u > /tmp/read.keys
diff /tmp/schema.keys /tmp/read.keys
```

Read 独有的 key → 补到 schema 里。**特别关注嵌套 struct/array**：`bound_services` / `model_cards` 这类，Read 里遍历子字段赋值，schema 侧必须声明 `Type: schema.TypeList, Elem: &schema.Resource{Schema: {...}}` 及所有子字段。

### Bug 4：Docs 缺 `(Required, ForceNew)` + 多余 NOTE

**症状 A**：`(Required)` 但没 `ForceNew`（跟 Bug 1 联动）；`(Optional)` 但缺 `Computed`。

**症状 B**：generator 会给每个 Optional/Required 加一句：

```markdown
-> **NOTE:** This parameter is only evaluated during resource creation and update.
   Modifying it in isolation will not trigger any action.
```

对普通字段没意义（改了必然触发 Update），是模板残留。

**症状 C**：immutable 字段的 NOTE：

```markdown
-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.
```

只要 schema 已经 `ForceNew: true`，这条 NOTE 就是废话——`(Required, ForceNew)` 已经说明了。

**修法**：修完 Bug 1 后把这些 NOTE 全删干净。

### Bug 5：Docs Example 是占位文字

**症状**：`## Example Usage` 段下面：

```
没有资源测试用例，请先通过资源测试用例后再生成示例代码。
```

或 datasource：

```
No resource test case, please pass the resource test case first before generating example code.
```

这是 generator 抓不到测试用例时的占位。**必须手写一份 example HCL** 才能过 doc lint 和用户使用。

**修法**：变量化 example，把复杂依赖（AI Gateway 等）拆成 `variable` 由用户注入：

```terraform
variable "name"       { default = "terraform-example" }
variable "gateway_id" {
  description = "The ID of an existing AI-type APIG gateway"
  type        = string
}

resource "alicloud_apig_ai_model_provider" "default" {
  gateway_id     = var.gateway_id
  model_provider = "openai"
  display_name   = var.name
}
```

### Bug 6：`d.SetId` 之后 Create 里额外调用 Read 之前的 `addDebug` 顺序

（暂无稳定 repro，观察项。）

### Bug 7：Update 的 `update = true` flag 只在部分字段变化时触发

**症状**：生成的 Update 函数里，`request["displayName"] = d.Get("display_name")` 总是无条件赋值，但只有 `d.HasChange("service_ids")` 触发 `update = true` → 只改 `display_name` 时 `update` 保持 `false` → PUT 永远不发 → 服务端保留旧值 → 测试的第 2 步（改 display_name）assertion 失败。

**检查**：

```bash
# 看 Update 函数里所有 request[...]=d.Get(...) 的字段，是否都对应有 d.HasChange 触发 update=true
grep -B1 -A15 "^func.*Update" alicloud/resource_alicloud_<x>.go | grep -E "HasChange|update = true|request\["
```

**修法**：为每个可 update 字段加 `d.HasChange` 分支：

```go
request["displayName"] = d.Get("display_name")
if d.HasChange("display_name") {
    update = true
}
if d.HasChange("service_ids") {
    update = true
    ...
}
```

## 完成信号

修完清单后运行本地静态检查（**不**跑 `go build/test`——见 SKILL.md Step 8）：

```bash
gofmt -l alicloud/                                    # 应输出空
go run scripts/document/document_check.go \
   website/docs/r/<name>.html.markdown                # 期望 "Finished!"
go run scripts/document/document_check.go \
   website/docs/d/<name>s.html.markdown               # 期望 "Finished!"
```

然后进 Step 9（写测试）。发现的每一类 bug 都在 Aone 留一条评论：`generator_bug_id`、`file:line`、修法 diff——积累样本给 `provider-resource-dev` 团队修模板。
