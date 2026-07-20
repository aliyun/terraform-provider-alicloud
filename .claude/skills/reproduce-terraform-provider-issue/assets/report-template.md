# {{TITLE}}

## 1. 现场状态

- 创建时间：`{{CREATE_WINDOW}}`
- 处置方式：`{{PRESERVE_OR_DESTROY}}`
- State 资源数量：`{{STATE_COUNT}}`
- 工作目录：`{{WORKDIR}}`
- 警告：`{{DO_NOT_APPLY_OR_DESTROY_NOTICE}}`

## 2. 核心结论

1. 云产品 API 行为：{{API_CONCLUSION}}
2. Provider 行为：{{PROVIDER_CONCLUSION}}
3. Terraform 后果：{{PLAN_CONCLUSION}}
4. 安全后续动作：{{NEXT_STEP}}

## 3. 环境与输入偏差

| 字段 | 值 |
|---|---|
| Terraform 版本 | `{{TERRAFORM_VERSION}}` |
| Provider 版本 | `{{PROVIDER_VERSION}}` |
| 账号身份 | `{{ACCOUNT_IDENTITY}}` |
| Region / Zone | `{{REGION_ZONE}}` |
| Aone | `{{AONE_ID}}` |
| 输入偏差 | {{INPUT_DEVIATIONS}} |

## 4. 实例清单

| 类型 | Terraform 地址 | 实例 ID / 访问地址 | 状态 |
|---|---|---|---|
| {{TYPE}} | `{{ADDRESS}}` | `{{INSTANCE_ID}}` | `{{STATUS}}` |

## 5. 关键创建 API 证据

- 时间：`{{TIMESTAMP}}`
- API：`{{CREATE_API}}`
- RequestId：`{{CREATE_REQUEST_ID}}`
- 创建资源 ID：`{{CREATED_ID}}`
- 白名单请求字段：

```json
{{CREATE_REQUEST_FIELDS_JSON}}
```

## 6. 关键读取 API 证据

| 时间 | API | RequestId | 目标资源 | 相关返回或缺失字段 |
|---|---|---|---|---|
| {{TIMESTAMP}} | `{{READ_API}}` | `{{READ_REQUEST_ID}}` | `{{TARGET_ID}}` | {{READ_FIELDS}} |

## 7. Apply API 完整时间线

{{APPLY_TIMELINE}}

## 8. 创建后的直接 API 验证

{{DIRECT_API_VERIFICATION}}

## 9. Refresh API 完整时间线

{{REFRESH_TIMELINE}}

## 10. Terraform 漂移结果

- 详细退出码：`{{PLAN_EXIT_CODE}}`
- 摘要：`{{PLAN_SUMMARY}}`

```json
{{REPLACEMENT_PATHS_JSON}}
```

## 11. Provider 机制分析

- Schema 位置：`{{SCHEMA_LOCATION}}`
- Read 映射位置：`{{READ_LOCATION}}`

```go
{{RELEVANT_PROVIDER_CODE}}
```

## 12. 需要产品团队确认的问题

1. {{PRODUCT_QUESTION_1}}
2. {{PRODUCT_QUESTION_2}}
3. {{PRODUCT_QUESTION_3}}

## 13. 只读验证命令

```bash
{{READ_ONLY_COMMANDS}}
```
