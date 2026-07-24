# PR 提交约定：aliyun/terraform-provider-alicloud

`aliyun/terraform-provider-alicloud` 仓有若干**独有的** PR 门禁与合并习惯，跟 GitHub 通用约定不同。SKILL.md Step 11-12 引用本文。

## 1. PR title / body 硬约定

### 1.1 title（新建资源）

```
New Resource: alicloud_<product>_<resource>; New Data Source: alicloud_<product>_<resources>
```

- **裸格式**——没有 `resource/…:` 前缀、没有 conventional-commit prefix、没有 Aone/工单/客户信息（GitHub public repo，SKILL.md Step 11.1 sanitize 规则）
- 单独资源改动用 `resource/alicloud_<x>: <what>`；单独 datasource 用 `datasource/alicloud_<xs>: <what>`
- 多个组件时用分号连接：`resource/…: …; datasource/…: …`

### 1.2 body

```markdown
## Summary
- <一行技术能力概述，别提客户/内部动机>
- <可选：第二条能力>

## Test plan
- [x] `TestAccAliCloud<X>_basic` PASS in <region> (<秒>s): create → update <field> → import → destroy.
- [x] `TestAccAliCloud<X>DataSource` PASS (<秒>s): ids exist/fake filter.
```

**Test plan** 里必须给出**已 PASS** 的用例名 + 时长——这是 reviewer 直接看到的验证证据。

## 2. 单 commit 门禁（`Pull Request Max Commits`）

CI job 严格限制**每个 PR 只允许 1 个 commit**（脚本 `if [[ ${commitNum} -gt 1 ]]; then exit 1`）。

### 2.1 首次提交前 squash

若 worktree 上多次迭代产生了多 commit，push 前 soft-reset 到 upstream base + 一次性 commit：

```bash
git fetch upstream master
MB=$(git merge-base HEAD upstream/master)
git reset --soft "$MB"
git commit -m "New Resource: alicloud_<x>; New Data Source: alicloud_<xs>"
bootstrap/github-identity.sh push api-tool-agent/terraform-provider-alicloud HEAD:<branch>
```

### 2.2 迭代阶段用 `--amend` + `--force-with-lease`

评审评论 / CI 失败 / ACC 重跑修完的每一次改动都 amend 到唯一那个 commit，然后 force-push：

```bash
git add <files>
git commit --amend --no-edit
bootstrap/github-identity.sh push api-tool-agent/terraform-provider-alicloud +HEAD:<branch>
```

- `+HEAD:<branch>` 显式 non-ff push，`bootstrap/github-identity.sh push` 内部会用 `--force-with-lease`
- 本约定**覆盖** CLAUDE.md「Prefer to create a new commit rather than amending」的通则——此 repo 强制 squash，amend 是唯一路径
- 只 force-push **自有 fork 的 PR-head 分支**（`api-tool-agent/…`），**绝不** force-push `aliyun/…` 上游

## 3. Rebase 应对 master 前进

master 提交速度快（约每天数个 PR merge）。PR 处在 review 期间 master 若合了改到 `alicloud/provider.go` 或 `alicloud/service_alicloud_<product>_v2.go` 的 PR，会撞冲突（`mergeStateStatus: DIRTY`）。

### 3.1 定期 rebase

```bash
git fetch upstream master
git rebase upstream/master
# 冲突多在 provider.go 的 DataSourcesMap / ResourcesMap，见 §4
git push --force-with-lease
```

### 3.2 不要 merge master

用 `git merge upstream/master` 会创出 merge commit → 违反 §2 的单 commit 门禁 → 又要 squash → 走了弯路。直接 rebase 保持单 commit。

## 4. `provider.go` / `service_alicloud_<product>_v2.go` 冲突解法

两个 map 尾巴上大家都在追加新条目 → 高频冲突。**冲突永远是"保留双方"**：

```go
DataSourcesMap: map[string]*schema.Resource{
<<<<<<< HEAD
    "alicloud_apig_services":            dataSourceAliCloudApigServices(),
=======
    "alicloud_apig_ai_model_providers":  dataSourceAliCloudApigAiModelProviders(),
>>>>>>> feat/apig-ai-model-provider
    ...
}
```

改成：

```go
DataSourcesMap: map[string]*schema.Resource{
    "alicloud_apig_ai_model_providers":  dataSourceAliCloudApigAiModelProviders(),
    "alicloud_apig_services":            dataSourceAliCloudApigServices(),
    ...
}
```

`ResourcesMap` 同理。map 里顺序无所谓，Go 会重排 hash。

## 5. `TestingCoverageRate` CI

- 依据 diff 找到本 PR 修改的资源；对每个资源检查所有 Optional/Required 属性是否至少出现在**一个** Step 的 Config 里
- `ImportStateVerifyIgnore` 只免"modified"部分，**不免**"testing"部分——见 `acc-test-writing-patterns.md` §3
- 报告格式：`<Attr> missing test cases` → 报告文件名找到那条 Optional，加进任一 Step 的 Config

## 6. `Content` CI（doc 三方一致性）

只对本 PR 改动的 doc 运行，检查：

- 标题 / description / subcategory / example 段完整
- 属性列表跟 schema 类型对齐
- 示例代码用 `terraform` 语言标签

修法：`go run scripts/document/document_check.go <doc>` 本地重现，逐条订正。

## 7. `Consistency` CI（doc ↔ schema 校验）

对每个改动的 resource/datasource，交叉核对 doc 声明的属性 / attributes 段是否跟 schema 一致（含 nested）。

- 常挂在：datasource nested schema 加了新字段但 doc `Attributes Reference` 段没同步（Bug 3 修完必看 doc 是否也补齐）
- 修法：手工同步 doc 的 attributes 段，运行 `scripts/consistency/consistency_check.go`——但**注意** consistency 脚本会 import alicloud 包（会编译，违反 Step 8），只在远程 CI 跑，本地读源码 diff 就行

## 8. `runner / errcheck` + `tflint`

标准 Go 静态检查——本 SOP Step 8 的 `gofmt -l` + `go vet` 覆盖不到 errcheck 的深度。CI 报错时按行修：一般是 `_ =` 忽略某个 error，或补 `if err != nil { return }`。

## 9. Reviewer 常见 comment 类型 + 快速回复

| Comment 类型 | 处理 |
|---|---|
| "Please add example for X" | 在 doc Example 段加最小工作 HCL，push |
| "Field X should be ForceNew" | schema 加 `ForceNew: true`，doc 改 `(Required, ForceNew)` |
| "Missing datasource test" | 补 `data_source_alicloud_<x>s_test.go`，跑远程 ACC |
| "Update commit message to match convention" | 只 amend commit message，不动内容：`git commit --amend -m "…" && push --force-with-lease` |
| APPROVED（无改动要求） | 无需 push；`wrap.sh sync` 回填 Aone 后交 PrWatchScheduler 跨会话看合并 |

## 10. 提交前最后一次自查

push 前**永远**跑一次 sanitize + `git log -p origin/master..HEAD` 扫敏感信息，见 SKILL.md Step 11.1 引用的 `bootstrap/pre-push-sanitize.sh`。
