# Evidence and sanitization contract

## Evidence layers

Keep these layers separate so a product team can identify ownership without reverse-engineering Terraform output:

| Layer | Required evidence | What it proves |
|---|---|---|
| Terraform input | Provider version and trigger fields | What the reporter asked the Provider to send |
| Create API | API name, allowlisted request fields, RequestId, created ID | The service accepted the request and produced an object |
| Read API | API name, RequestId, returned/missing fields | What the service can round-trip after creation |
| Related API | RequestId and relationship fields | Whether data exists on a child/attachment instead of the primary object |
| Provider state | State fields immediately after Read/refresh | How the Provider mapped the API response |
| Plan | detailed exit code, actions, replace paths | The Terraform lifecycle consequence |

## Timeline row contract

Every row should contain:

```text
timestamp | API | RequestId | target IDs | status | allowlisted request | observation
```

Preserve repeated polling calls. They prove whether a missing field was transient or stable and give the product team multiple RequestIds.

## Allowed request fields

The bundled parser allows only diagnostic fields such as:

- region/zone;
- resource IDs and names;
- resource type/protocol/storage/encryption selections;
- VPC/VSwitch and related topology IDs;
- access-group/rule semantics;
- CIDR blocks used by the reproduction.

Do not expand the allowlist to include authentication, authorization, signature, session, cookie, token, password, user-data, private-key, or arbitrary payload fields.

## Sensitive artifacts

Assume the following contain secrets even when console output looks redacted:

- saved Terraform plan binaries;
- raw `TF_LOG` / Provider debug traces;
- shell environment dumps;
- provider configuration values;
- signed URLs and authorization headers.

Use raw artifacts only locally while extracting allowlisted evidence. Delete them after producing the sanitized timeline/report. Never upload raw traces to AutomationAgent or attach them to Aone.

## Report requirements

The report must include:

1. scene status and preservation/cleanup warning;
2. environment and any deviation from reporter input;
3. complete instance inventory;
4. create request/response evidence;
5. first read response showing the disputed fields;
6. complete apply API timeline, including polls;
7. direct post-create API verification;
8. complete refresh API timeline;
9. state diff and replacement paths;
10. Provider code/schema locations when known;
11. questions for the owning product team;
12. read-only commands for continued investigation.

## Conclusion discipline

- Say “API returned empty/missing” only when a response and RequestId prove it.
- Say “Provider cleared state” only after inspecting state after Read/refresh.
- Say “Terraform replaces” only after recording the plan action and replace paths.
- Do not flatten “API contract ambiguity” and “Provider mapping defect” into one claim.
- If an isolated network replaced inaccessible reporter IDs, say so near the top of the report.
