# cap-bootstrap-test-drift

## 缺口类型

bootstrap 测试与当前脚本行为漂移。

## 阻塞任务

本次 Terraform 资源开发 skill 改造的新增测试和 `bootstrap/verify.sh` 均通过,但执行全量 `test/*.sh` 时,既有 `plan_test.sh` 与 `scan_test.sh` 失败。失败集中在测试期望与当前 plan/scan 输出结构不一致,不属于本次 provider-resource-dev skill 改动。

## 当前发现

- `test/plan_test.sh` 失败:计划文件中未按旧期望包含 WI-100/WI-400/WI-500,且 action/confidence 字段断言不再匹配当前输出。
- `test/scan_test.sh` 失败:当前 scan 会按 pool/category 展开多条记录,测试仍按单条记录和旧字段数量断言;缓存行为也导致 empty inbox 用例读到前一次 scan 结果。
- `test/verify_test.sh` 中的 `FAIL gh-cred` 是测试用 stub 刻意制造的预期失败;真实 `bootstrap/verify.sh` 已通过。

## 建议补丁

1. 单独校准 `plan_test.sh` 与 `scan_test.sh` 的断言,明确当前行为是“多池/多 category 展开”还是应在脚本层去重。
2. 为 scan 测试隔离 cache 目录,避免一个用例污染下一个用例。
3. 给全量测试入口增加 fail-fast 或汇总退出码,避免子测试失败但外层命令退出 0。

## 置信度

high:基于本 worktree 内执行 `for f in test/*.sh; do bash "$f"; done` 的输出。

## 关联

Aone: https://project.aone.alibaba-inc.com/v2/project/2100304/req/83695664
