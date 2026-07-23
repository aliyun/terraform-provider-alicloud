# Error Diagnosis Patterns for Remote ACC Tests

## 0. 归类入口

每次 FAIL 必须先在 SKILL.md §6「FAIL 定性:四分类」里归到 A/B/C/D 中的一种:

- **A** — 后端云产品 API 问题(行为/文档/合同不合理)
- **B** — 测试用例问题(HCL / precheck / 依赖固定 ID)
- **C1** — 手写 Resource/DataSource 代码 bug
- **C2** — CloudSpec 资源定义 bug(修 cspec)
- **D** — 生成器 bug(cspec 正确、生成产物错;修生成器再重新生成)

本文件的每条 pattern 末尾都标注了归属分类,可作为快速匹配材料。**不能只匹配 pattern 就交差**——SKILL.md §6.1 的取证步骤(FAIL 位置 → tf-debug 对齐 → API 契约核对 → provider 代码核对 → cspec 核对 → 若资源为生成器产出,追加生成器模板/规则核对)必须走完,证据齐了再定性。

## 1. No Tests Executed (Name Mismatch)

When `--test-case` is **not** specified, the runner matches `TestAccAliCloud{Namespace}{ResourceTypeCode}*`. If `run.log` reports 0 tests:

- Check actual function names locally: `grep -n "^func TestAcc" alicloud/*_test.go`
- Fix casing/spelling, then either:
  - Pass `--test-case <exact-name>` to skip discovery, OR
  - Rename the test functions to match the expected pattern

When `--test-case` **is** specified and FC reports `指定的测试用例 X 在 alicloud/ 中未找到`:
- Just a typo — grep the local repo for the correct name and resubmit (no commit/push needed).

## 2. `Unsupported argument: An argument named "X" is not expected here`

Test HCL still references a deleted/renamed field. Update the test config to use the correct field name.

## 3. Cloud API Error (`InvalidParameter`, `EntityNotExist`, `Throttling`)

Locate RequestId in `tf-debug.log`:

```bash
grep -E "RequestId|ErrorCode|ErrorMessage" ./acctest_logs/*_tf-debug.log | head -40
```

| Error Code | Action |
|------------|--------|
| `InvalidParameter` | Check if a required param is missing after field removal |
| `EntityNotExist` | Region or quota issue — try another region |
| `Throttling` | Reduce parallelism or wait and retry |

## 4. `daring resource` / Destroy Timeout (Subscription Resources)

Subscription resources can't be destroyed via API. **Treat as PASS** if all Apply/Read steps completed successfully and only the Destroy step failed.

## 5. `unknown error` at Initialization

Likely wrong region. Some resources are only available in specific regions.

## 6. FC-Side Failures (vs Business Test Failures)

| `status` | `fcStatus` | `fcInvocationErrorMessage` | Diagnosis |
|---|---|---|---|
| `failed` | `Succeeded` | empty | FC ran cleanly, but one or more Go test cases failed → look at `tf-debug.log` for business errors |
| `failed` | `Failed` | non-empty | FC instance crashed (panic / OOM / timeout) before completing → check FC console runtime logs; the error message is the FC-side cause |
| `failed` | `Stopped` | usually empty | Task was cancelled via the `cancel` API |
| `failed` | `Expired` | empty | FC never picked up the invocation (queue overflow / quota) → wait & retry |
| `failed` | null | n/a | FC tracking lost (e.g. record pruned). Final verdict still trusted from OSS `run.log` |

When `fcStatus` indicates an infrastructure-level failure (Failed / Stopped / Expired), the `run.log` may be truncated because the function process was killed before flushing terminal state. Use `tf-debug.log` plus the FC console for diagnosis in that case.

## 7. `QuotaExceeded.<Resource>` at Create Step

Test account quota in that region is exhausted — usually because past failed runs left undeleted resources behind. Not a provider bug.

**Diagnose** — find which region and resource hit the limit:

```bash
grep -E "QuotaExceeded|HostId" ./acctest_logs/*_run.log | head -10
# Example: "Code":"QuotaExceeded.Vpc","HostId":"vpc.cn-beijing.aliyuncs.com"
```

**Clean up** — from the `terraform-provider-alicloud` repo root, run sweep:

```bash
make sweep REGION=<region> RESOURCE=<resource_type>
# Example: clean leftover VPCs in Beijing
make sweep REGION=cn-beijing RESOURCE=alicloud_vpc
```

This deletes test-prefixed (`tf-test*`, `tf_test*`, etc.) resources of that type in that region, freeing quota. Re-submit the ACC test after sweep completes.

## Log Quick Reference

```bash
# API request/response pairs
grep -E "(POST|GET) https|RequestId|ErrorCode" ./acctest_logs/*_tf-debug.log

# Test pass/fail summary
grep -E "^--- (PASS|FAIL|SKIP)" ./acctest_logs/*_run.log

# FC instance error (if fcStatus=Failed)
# → check FC console under "task" / "log" tab; or look at the status response's fcInvocationErrorMessage
```
