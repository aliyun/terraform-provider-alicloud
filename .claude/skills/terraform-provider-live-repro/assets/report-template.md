# {{TITLE}}

## 1. Scene status

- Created at: `{{CREATE_WINDOW}}`
- Disposition: `{{PRESERVE_OR_DESTROY}}`
- State resource count: `{{STATE_COUNT}}`
- Working directory: `{{WORKDIR}}`
- Warning: `{{DO_NOT_APPLY_OR_DESTROY_NOTICE}}`

## 2. Executive conclusion

1. Cloud API behavior: {{API_CONCLUSION}}
2. Provider behavior: {{PROVIDER_CONCLUSION}}
3. Terraform consequence: {{PLAN_CONCLUSION}}
4. Safe next step: {{NEXT_STEP}}

## 3. Environment and deviations

| Field | Value |
|---|---|
| Terraform | `{{TERRAFORM_VERSION}}` |
| Provider | `{{PROVIDER_VERSION}}` |
| Account identity | `{{ACCOUNT_IDENTITY}}` |
| Region / zone | `{{REGION_ZONE}}` |
| Aone | `{{AONE_ID}}` |
| Input deviations | {{INPUT_DEVIATIONS}} |

## 4. Instance inventory

| Type | Terraform address | Instance ID / endpoint | Status |
|---|---|---|---|
| {{TYPE}} | `{{ADDRESS}}` | `{{INSTANCE_ID}}` | `{{STATUS}}` |

## 5. Key create API evidence

- Time: `{{TIMESTAMP}}`
- API: `{{CREATE_API}}`
- RequestId: `{{CREATE_REQUEST_ID}}`
- Created ID: `{{CREATED_ID}}`
- Allowlisted request fields:

```json
{{CREATE_REQUEST_FIELDS_JSON}}
```

## 6. Key read API evidence

| Time | API | RequestId | Target | Relevant returned/missing fields |
|---|---|---|---|---|
| {{TIMESTAMP}} | `{{READ_API}}` | `{{READ_REQUEST_ID}}` | `{{TARGET_ID}}` | {{READ_FIELDS}} |

## 7. Apply API timeline

{{APPLY_TIMELINE}}

## 8. Direct post-create verification

{{DIRECT_API_VERIFICATION}}

## 9. Refresh API timeline

{{REFRESH_TIMELINE}}

## 10. Terraform drift

- Detailed exit code: `{{PLAN_EXIT_CODE}}`
- Summary: `{{PLAN_SUMMARY}}`

```json
{{REPLACEMENT_PATHS_JSON}}
```

## 11. Provider mechanism

- Schema: `{{SCHEMA_LOCATION}}`
- Read mapping: `{{READ_LOCATION}}`

```go
{{RELEVANT_PROVIDER_CODE}}
```

## 12. Questions for the product team

1. {{PRODUCT_QUESTION_1}}
2. {{PRODUCT_QUESTION_2}}
3. {{PRODUCT_QUESTION_3}}

## 13. Read-only verification commands

```bash
{{READ_ONLY_COMMANDS}}
```
