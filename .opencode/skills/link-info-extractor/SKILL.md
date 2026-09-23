---
name: link-info-extractor
description: Extract requirement details from a URL (Aone workitem, GitLab Code Review, etc.) using MCP tools. Use when the user provides a link to a requirement or review system.
metadata:
  version: "2.0.0"
  domain: requirements
  triggers: aone link, workitem, code review, codereview, requirement url, extract from link
---

# Link Info Extractor

When a user provides a URL pointing to a requirement or review system, use this skill to fetch and extract structured information from that link.

## Supported Sources

| Source | URL Pattern | Parameters to Extract |
|--------|------------|----------------------|
| Aone Workitem | `https://project.aone.alibaba-inc.com/v2/project/<projectId>/req/<workitemId>` | `projectId`, `workitemId` |
| GitLab Code Review | `https://code.alibaba-inc.com/<group>/<project>/codereview/<reviewId>` | `group`, `project`, `reviewId` |

## Step 1: Identify the Source

Match the URL against known patterns above to determine the source type.

## Step 2: Fetch Information via MCP Tools

### Source A: Aone Workitem

1. Call `coop_query_workitem_detail(workitemId)` to get the workitem description, title, and status
2. Call `coop_get_workitem_comments(workitemId)` to get all comments

> **Note:** The MCP tool name prefix (e.g., `coop_`) may vary across different app registrations, but the tool name suffix (`query_workitem_detail`, `get_workitem_comments`) stays the same. Always match by suffix.

### Source B: GitLab Code Review

Use the `alibaba_code` MCP to fetch Code Review information:

1. Get Code Review details (description, status, changed files)
2. Get comments (may contain technical specs or requirement details)

> **Note:** The MCP name prefix may vary across different app registrations, but the underlying tool names stay the same. Always match by suffix.

## Step 3: Extract Structured Requirements

From the fetched information, extract a structured technical spec. Focus on:

1. **Attribute definitions**
   - Attribute name (Terraform-side `snake_case` naming)
   - Type (Boolean / String / Integer / List / Map, etc.)
   - Required / Optional
   - ForceNew (resource must be recreated on change)
   - Computed (value returned by server)

2. **API mapping**
   - `requestPath` — parameter name for Create/Update API calls (typically PascalCase)
   - `responsePath` — field name in Read API response
   - Which API operations are involved: Create / Read / Update

3. **Constraints**
   - Default values
   - Value range restrictions
   - Mutual exclusion or dependencies with other attributes
   - Minimum supported provider version

## Output Format

Return structured requirement info:

```
## Requirement Summary

**Source**: [<source description>](<url>)
**Title**: <title>

### New Attributes

| Attribute | Type | Required | ForceNew | Computed | Description |
|-----------|------|----------|----------|----------|-------------|
| <name> | <type> | <Y/N> | <Y/N> | <Y/N> | <description> |

### API Mapping

- **Create**: `<API name>` — request path: `<requestPath>`
- **Read**: `<API name>` — response path: `<responsePath>`
- **Update**: `<API name>` — request path: `<requestPath>`

### Constraints

- <constraint 1>
- <constraint 2>
- ...
```

## Important Notes

1. Technical specs may appear in the description OR in comments — always check both
2. If the requirement is unclear, flag uncertain parts and ask the user to confirm
3. Extracted API parameter names must match Alibaba Cloud OpenAPI documentation
4. Distinguish which attributes support Update vs. Create-only (ForceNew)
5. This skill is NOT needed when the user describes requirements directly without a link
