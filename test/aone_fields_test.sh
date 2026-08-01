#!/usr/bin/env bash
# Hermetic tests for bootstrap/aone-fields.sh. All workitem/field data is synthetic.

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root="$(cd "$test_dir/.." && pwd)"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
mkdir -p "$tmp/bin"
log="$tmp/a1.log"

# Keep the field-metadata cache out of the developer's real .my-day/cache, and
# off by default here: each case below drives the a1 stub with different env
# (A1_OPTIONS_FAIL_FIELD, A1_GENERIC_FIXTURES, ...) under the same project/type,
# so a shared cache would serve case N's payload to case N+1 and silently void
# the fail-closed assertions. The cache path itself is covered by its own case
# at the end of this file.
export JARVIS_CACHE_DIR="$tmp/cache"
export JARVIS_FIELD_META_TTL=0

cat > "$tmp/workitem.json" <<'JSON'
{"fields":[
  {"identifier":"workitemType","value":"36","displayValue":"功能缺陷"},
  {"identifier":"space","value":"528766","displayValue":"synthetic-project"},
  {"identifier":"140282","value":"","displayValue":""},
  {"identifier":"140283","value":"already-set","displayValue":"已设置"},
  {"identifier":"140284","value":"","displayValue":""}
]}
JSON

cat > "$tmp/fields.json" <<'JSON'
[
  {"identifier":"status","displayName":"状态","isRequired":true,"sourceType":"system","format":"list","options":[]},
  {"identifier":"140282","displayName":"需求分类","description":"synthetic","isRequired":true,"sourceType":"team","format":"list","options":[]},
  {"identifier":"140283","displayName":"已有字段","isRequired":true,"sourceType":"team","format":"list","options":[]},
  {"identifier":"140284","displayName":"来源","isRequired":true,"sourceType":"team","format":"list","options":[{"value":"manual","displayValue":"人工"}]},
  {"identifier":"140285","displayName":"选填字段","isRequired":false,"sourceType":"team","format":"input","options":[]}
]
JSON

cat > "$tmp/bin/a1" <<'STUB'
#!/usr/bin/env bash
echo "$*" >> "$A1_TEST_LOG"
if [ "$1 $2 $3" = "project workitem get" ]; then
  if [ -f "$A1_TEST_STATE.updated" ] && [ "${A1_READBACK_FAIL:-}" = "1" ]; then
    exit 8
  fi
  get_count="$(cat "$A1_TEST_STATE.get_count" 2>/dev/null || printf '0')"
  get_count=$((get_count + 1))
  printf '%s\n' "$get_count" > "$A1_TEST_STATE.get_count"
  # Read 1 is preflight's own context read; read 2 is the explicit pre-mutation
  # drift check (noop readback, or the read immediately before update). Drift is
  # injected from read 2 so it lands on that check. This used to be read 3 only
  # because preflight reached the field list through a nested `missing` that
  # re-fetched the same workitem.
  if [ "${A1_DRIFT_BEFORE_UPDATE:-}" = "project" ] \
      && [ ! -f "$A1_TEST_STATE.updated" ] && [ "$get_count" -ge 2 ]; then
    jq '(.fields[] | select(.identifier=="space") | .value) = "528766"' "$A1_TEST_WI"
    exit 0
  fi
  if [ "${A1_DRIFT_AFTER_UPDATE:-}" = "type" ] \
      && [ -f "$A1_TEST_STATE.updated" ]; then
    jq '(.fields[] | select(.identifier=="workitemType") | .value) = "99"' "$A1_TEST_WI"
    exit 0
  fi
  if [ "${A1_CONFLICT_AFTER_UPDATE:-}" = "140282" ] \
      && [ -f "$A1_TEST_STATE.updated" ]; then
    jq '(.fields[] | select(.identifier=="140282") | .value) = "manual"' "$A1_TEST_WI"
    exit 0
  fi
  if [ "${A1_CONFLICT_AFTER_UPDATE:-}" = "140097" ] \
      && [ -f "$A1_TEST_STATE.updated" ]; then
    jq '(.fields[] | select(.identifier=="140097") | .value) = "redis"' "$A1_TEST_WI"
    exit 0
  fi
  cat "$A1_TEST_WI"; exit 0
fi
if [ "$1 $2 $3 $4" = "project workitem field list" ]; then
  cat "$A1_TEST_FIELDS"; exit 0
fi
if [ "$1 $2 $3 $4" = "project workitem field options" ]; then
  if [ "${A1_OPTIONS_FAIL_FIELD:-}" = "$5" ]; then
    exit 7
  fi
  if [ "${A1_GENERIC_FIXTURES:-}" = "1" ]; then
    case "$5" in
      140097) printf '[{"value":"","identifier":"terraform","displayValue":"Terraform"},{"value":"mongodb","identifier":"ignored-second-key","displayValue":"MongoDB"}]\n' ;;
      140282) printf '[{"value":"defined","displayValue":"有OpenAPI，资源未定义，开放平台维护"},{"value":"manual","displayValue":"手工资源"}]\n' ;;
      107239) printf '[{"Identifier":"906688","Name":"API 工具","Path":"产品/API 工具"},{"Identifier":"891438","Name":"云数据库 MongoDB 版","Path":"产品/MongoDB"}]\n' ;;
    esac
    exit 0
  fi
  if [ "${A1_AUTO_FIXTURES:-}" != "1" ]; then
    printf '[{"identifier":"bug","value":"bug","displayValue":"缺陷"}]\n'; exit 0
  fi
  case "$5" in
    140097) printf '[{"value":"mongodb","displayValue":"MongoDB/云数据库 MongoDB 版"},{"value":"redis","displayValue":"Redis/云数据库 Redis 版"}]\n' ;;
    140282) printf '[{"value":"defined","displayValue":"有OpenAPI，资源未定义，开放平台维护"},{"value":"manual","displayValue":"手工资源"}]\n' ;;
    107239) printf '[{"value":"891438","displayValue":"云数据库 MongoDB 版"},{"value":"906688","displayValue":"API 工具"}]\n' ;;
    *) printf '[{"identifier":"bug","value":"bug","displayValue":"缺陷"}]\n' ;;
  esac
  exit 0
fi
if [ "$1 $2 $3" = "project workitem update" ]; then
  [ "${A1_UPDATE_FAIL:-}" != "1" ] || exit 6
  touch "$A1_TEST_STATE.updated"
  if [ "${A1_MUTATE_ON_UPDATE:-}" = "1" ]; then
    shift 3
    while [ "$#" -gt 0 ]; do
      if [ "$1" = "--cfs" ] && [ "$#" -ge 2 ]; then
        field_id="${2%%=*}"
        field_value="${2#*=}"
        next="$A1_TEST_STATE.next"
        jq --arg id "$field_id" --arg value "$field_value" '
          (.fields // []) |= map(
            if ((.fieldIdentifier // .identifier // "") | tostring) == $id
            then .value = $value | .displayValue = $value
            else .
            end)
        ' "$A1_TEST_WI" > "$next" && mv "$next" "$A1_TEST_WI"
        shift 2
      else
        shift
      fi
    done
  fi
  exit 0
fi
exit 9
STUB
chmod +x "$tmp/bin/a1"

export A1_TEST_LOG="$log" A1_TEST_WI="$tmp/workitem.json" A1_TEST_FIELDS="$tmp/fields.json"
export A1_TEST_STATE="$tmp/a1-state"
export JARVIS_A1=a1 JARVIS_ROOT="$root"
PASS=0; FAIL=0
ok() { echo "PASS: $1"; PASS=$((PASS + 1)); }
bad() { echo "FAIL: $1"; FAIL=$((FAIL + 1)); }

out="$(PATH="$tmp/bin:$PATH" bash "$root/bootstrap/aone-fields.sh" missing 9001 2>"$tmp/missing.err")"; rc=$?
[ "$rc" -eq 0 ] && ok "missing exits 0" || bad "missing rc=$rc"
ids="$(printf '%s' "$out" | jq -r 'map(.id) | sort | join(",")')"
[ "$ids" = "140282,140284" ] && ok "only empty required custom fields returned" || bad "unexpected ids=$ids"
printf '%s' "$out" | jq -e '.[] | select(.id=="140282") | .optionsSource=="field_options" and .options[0].value=="bug"' >/dev/null \
  && ok "empty field-list options supplemented from field options API" || bad "field options not supplemented: $out"
printf '%s' "$out" | jq -e '.[] | select(.id=="140284") | .optionsSource=="field_list" and .options[0].value=="manual"' >/dev/null \
  && ok "existing field-list options preserved" || bad "field-list options lost: $out"
[ "$(grep -c 'field options 140282' "$log" || true)" -eq 1 ] \
  && [ "$(grep -c 'field options 140284' "$log" || true)" -eq 0 ] \
  && ok "options API called only for empty option list" || bad "unexpected options calls: $(cat "$log")"

: > "$log"
fill_out="$(PATH="$tmp/bin:$PATH" bash "$root/bootstrap/aone-fields.sh" fill 9001 '140282=bug' '140284=manual' 2>&1)"; rc=$?
[ "$rc" -eq 0 ] && grep -q -- '--cfs 140282=bug --cfs 140284=manual' "$log" \
  && ok "fill forwards explicit repeated --cfs assignments" || bad "fill failed rc=$rc log=$(cat "$log") out=$fill_out"

PATH="$tmp/bin:$PATH" bash "$root/bootstrap/aone-fields.sh" fill 9001 '140282=' >/dev/null 2>&1; rc=$?
[ "$rc" -eq 2 ] && ok "fill rejects empty value" || bad "empty value should exit 2, got $rc"
PATH="$tmp/bin:$PATH" bash "$root/bootstrap/aone-fields.sh" missing 'bad/id' >/dev/null 2>&1; rc=$?
[ "$rc" -eq 2 ] && ok "missing rejects unsafe workitem id" || bad "unsafe id should exit 2, got $rc"

cat > "$tmp/workitem.json" <<'JSON'
{"title":"all required fields already present","fields":[
  {"identifier":"workitemType","value":"36","displayValue":"需求问题"},
  {"identifier":"space","value":"1086837","displayValue":"Terraform customer"},
  {"identifier":"140282","value":"defined","displayValue":"有OpenAPI，资源未定义，开放平台维护"}
]}
JSON
cat > "$tmp/fields.json" <<'JSON'
[
  {"identifier":"140282","displayName":"Terraform需求类型","isRequired":true,"sourceType":"team","format":"list","options":[]}
]
JSON
: > "$log"
preflight_out="$(PATH="$tmp/bin:$PATH" bash "$root/bootstrap/aone-fields.sh" preflight 9001 1086837 2>"$tmp/preflight.err")"; rc=$?
[ "$rc" -eq 0 ] && printf '%s' "$preflight_out" | jq -e '
  .status == "ok" and .errorType == null and .workitemId == "9001" and
  .project == "1086837" and .workitemType == "36" and
  .assignments == [] and .unresolved == [] and .readback == [] and .filled == false
' >/dev/null && [ "$(grep -c 'project workitem update' "$log" || true)" -eq 0 ] \
  && ok "preflight complete item is a stable no-op" \
  || bad "preflight complete item failed rc=$rc out=$preflight_out err=$(cat "$tmp/preflight.err") log=$(cat "$log")"

rm -f "$A1_TEST_STATE.get_count" "$A1_TEST_STATE.updated"
: > "$log"
preflight_out="$(PATH="$tmp/bin:$PATH" A1_DRIFT_BEFORE_UPDATE=project \
  bash "$root/bootstrap/aone-fields.sh" preflight 9001 1086837 2>/dev/null)"; rc=$?
[ "$rc" -eq 3 ] && printf '%s' "$preflight_out" | jq -e '
  .failureReason=="context_drift_noop_readback" and .filled==false
' >/dev/null && [ "$(grep -c 'workitem update' "$log" || true)" -eq 0 ] \
  && ok "complete no-op item still fences context drift" \
  || bad "no-op context drift was not fenced rc=$rc out=$preflight_out log=$(cat "$log")"

cat > "$tmp/workitem.json" <<'JSON'
{"title":"support alicloud_mongodb_sharding_instance named whitelist groups","fields":[
  {"identifier":"workitemType","value":"36","displayValue":"需求问题"},
  {"identifier":"space","value":"1086837","displayValue":"Terraform customer"},
  {"identifier":"140097","value":"","displayValue":""},
  {"identifier":"140282","value":"","displayValue":""},
  {"identifier":"107239","value":"","displayValue":""}
]}
JSON
cat > "$tmp/fields.json" <<'JSON'
[
  {"identifier":"140097","displayName":"涉及云产品","isRequired":true,"sourceType":"team","format":"list","options":[]},
  {"identifier":"140282","displayName":"Terraform需求类型","isRequired":true,"sourceType":"team","format":"list","options":[]},
  {"identifier":"107239","displayName":"归属产品","isRequired":true,"sourceType":"team","format":"list","options":[]}
]
JSON
: > "$log"
auto_out="$(PATH="$tmp/bin:$PATH" A1_AUTO_FIXTURES=1 bash "$root/bootstrap/aone-fields.sh" auto-fill 9001 2>"$tmp/auto.err")"; rc=$?
[ "$rc" -eq 0 ] && printf '%s' "$auto_out" | jq -e '.filled == true and (.assignments | length) == 3' >/dev/null \
  && ok "auto-fill resolves configured default and unique title matches" \
  || bad "auto-fill failed rc=$rc out=$auto_out err=$(cat "$tmp/auto.err")"
grep -q -- '--cfs 140097=mongodb' "$log" \
  && grep -q -- '--cfs 140282=defined' "$log" \
  && grep -q -- '--cfs 107239=891438' "$log" \
  && ok "auto-fill submits every required field in one update" \
  || bad "auto-fill assignments missing: $(cat "$log")"

cat > "$tmp/workitem.json" <<'JSON'
{"title":"Provider required fields regression","fields":[
  {"identifier":"workitemType","value":"36","displayValue":"需求问题"},
  {"identifier":"space","value":"1086837","displayValue":"Terraform customer"},
  {"identifier":"140097","value":"","displayValue":""},
  {"identifier":"140282","value":"","displayValue":""},
  {"identifier":"107239","value":"","displayValue":""}
]}
JSON
cat > "$tmp/fields.json" <<'JSON'
[
  {"identifier":"140097","displayName":"涉及云产品","isRequired":true,"sourceType":"team","format":"list","options":[]},
  {"identifier":"140282","displayName":"Terraform需求类型","isRequired":true,"sourceType":"team","format":"list","options":[]},
  {"identifier":"107239","displayName":"归属产品","isRequired":true,"sourceType":"team","format":"list","options":[]}
]
JSON
rm -f "$A1_TEST_STATE.updated"
: > "$log"
preflight_out="$(PATH="$tmp/bin:$PATH" A1_GENERIC_FIXTURES=1 A1_MUTATE_ON_UPDATE=1 \
  bash "$root/bootstrap/aone-fields.sh" preflight 84608993 1086837 2>"$tmp/generic.err")"; rc=$?
[ "$rc" -eq 0 ] && printf '%s' "$preflight_out" | jq -e '
  .status == "ok" and .filled == true and .readback == [] and
  ([.assignments[] | select(.id=="140097" and .value=="terraform" and .source=="validated_fallback")] | length) == 1 and
  ([.assignments[] | select(.id=="140282" and .value=="defined" and .source=="configured_default")] | length) == 1 and
  ([.assignments[] | select(.id=="107239" and .value=="906688" and .source=="validated_fallback")] | length) == 1
' >/dev/null \
  && [ "$(grep -c 'project workitem update' "$log" || true)" -eq 1 ] \
  && [ "$(grep -c 'project workitem field options' "$log" || true)" -eq 3 ] \
  && ok "preflight validates generic fallbacks, empty value identifiers, uppercase product options, one update and readback" \
  || bad "generic preflight failed rc=$rc out=$preflight_out err=$(cat "$tmp/generic.err") log=$(cat "$log")"

# An explicit value survives preflight; the completed item performs no second update.
: > "$log"
preflight_out="$(PATH="$tmp/bin:$PATH" A1_GENERIC_FIXTURES=1 A1_MUTATE_ON_UPDATE=1 \
  bash "$root/bootstrap/aone-fields.sh" preflight 84608993 1086837)"; rc=$?
[ "$rc" -eq 0 ] && [ "$(grep -c 'project workitem update' "$log" || true)" -eq 0 ] \
  && ok "preflight never overwrites explicit values" \
  || bad "explicit values were not a no-op rc=$rc out=$preflight_out log=$(cat "$log")"

# Project drift fails before field lookup/update.
: > "$log"
preflight_out="$(PATH="$tmp/bin:$PATH" A1_GENERIC_FIXTURES=1 \
  bash "$root/bootstrap/aone-fields.sh" preflight 84608993 528766)"; rc=$?
[ "$rc" -eq 3 ] && printf '%s' "$preflight_out" | jq -e '
  .status=="failed" and .errorType=="preflight_validation_failed" and
  .failureReason=="project_mismatch" and
  .unresolved[0].expectedProject=="528766" and .unresolved[0].actualProject=="1086837"
' >/dev/null \
  && [ "$(grep -c 'field list\\|workitem update' "$log" || true)" -eq 0 ] \
  && ok "preflight project mismatch fails closed before mutation" \
  || bad "project mismatch contract failed rc=$rc out=$preflight_out log=$(cat "$log")"

reset_generic_item() {
  rm -f "$A1_TEST_STATE.updated" "$A1_TEST_STATE.get_count"
  jq '(.fields // []) |= map(
    if ((.identifier // "") == "140097" or (.identifier // "") == "140282" or
        (.identifier // "") == "107239")
    then .value = "" | .displayValue = ""
    else .
    end)' "$tmp/workitem.json" > "$tmp/reset.json"
  mv "$tmp/reset.json" "$tmp/workitem.json"
  : > "$log"
}

reset_generic_item
preflight_out="$(PATH="$tmp/bin:$PATH" A1_GENERIC_FIXTURES=1 A1_OPTIONS_FAIL_FIELD=140097 \
  bash "$root/bootstrap/aone-fields.sh" preflight 84608993 1086837 2>/dev/null)"; rc=$?
[ "$rc" -eq 3 ] && printf '%s' "$preflight_out" | jq -e '
  .failureReason=="unresolved_required_fields" and
  any(.unresolved[]; .id=="140097" and .reason=="options_lookup_error") and
  .filled==false
' >/dev/null && [ "$(grep -c 'workitem update' "$log" || true)" -eq 0 ] \
  && ok "options lookup failure never blind-fills" \
  || bad "options lookup failure did not fail closed rc=$rc out=$preflight_out log=$(cat "$log")"

reset_generic_item
preflight_out="$(PATH="$tmp/bin:$PATH" A1_GENERIC_FIXTURES=1 A1_UPDATE_FAIL=1 \
  bash "$root/bootstrap/aone-fields.sh" preflight 84608993 1086837 2>/dev/null)"; rc=$?
[ "$rc" -eq 3 ] && printf '%s' "$preflight_out" | jq -e '
  .failureReason=="update_error" and .filled==false and (.assignments|length)==3
' >/dev/null && [ "$(grep -c 'workitem update' "$log" || true)" -eq 1 ] \
  && ok "update failure returns stable validation failure" \
  || bad "update failure contract failed rc=$rc out=$preflight_out log=$(cat "$log")"

reset_generic_item
preflight_out="$(PATH="$tmp/bin:$PATH" A1_GENERIC_FIXTURES=1 A1_MUTATE_ON_UPDATE=1 A1_READBACK_FAIL=1 \
  bash "$root/bootstrap/aone-fields.sh" preflight 84608993 1086837 2>/dev/null)"; rc=$?
[ "$rc" -eq 3 ] && printf '%s' "$preflight_out" | jq -e '
  .failureReason=="readback_error" and .filled==true and .readback==[]
' >/dev/null && [ "$(grep -c 'workitem update' "$log" || true)" -eq 1 ] \
  && ok "readback lookup failure fails closed after one update" \
  || bad "readback lookup failure contract failed rc=$rc out=$preflight_out log=$(cat "$log")"

reset_generic_item
preflight_out="$(PATH="$tmp/bin:$PATH" A1_GENERIC_FIXTURES=1 \
  bash "$root/bootstrap/aone-fields.sh" preflight 84608993 1086837 2>/dev/null)"; rc=$?
[ "$rc" -eq 3 ] && printf '%s' "$preflight_out" | jq -e '
  .failureReason=="readback_not_empty" and .filled==true and (.readback|length)==3
' >/dev/null && [ "$(grep -c 'workitem update' "$log" || true)" -eq 1 ] \
  && ok "non-empty readback fails closed" \
  || bad "non-empty readback contract failed rc=$rc out=$preflight_out log=$(cat "$log")"

reset_generic_item
preflight_out="$(PATH="$tmp/bin:$PATH" A1_GENERIC_FIXTURES=1 A1_MUTATE_ON_UPDATE=1 \
  A1_DRIFT_BEFORE_UPDATE=project \
  bash "$root/bootstrap/aone-fields.sh" preflight 84608993 1086837 2>/dev/null)"; rc=$?
[ "$rc" -eq 3 ] && printf '%s' "$preflight_out" | jq -e '
  .failureReason=="context_drift_before_update" and .filled==false
' >/dev/null && [ "$(grep -c 'workitem update' "$log" || true)" -eq 0 ] \
  && ok "project drift immediately before update fails closed" \
  || bad "pre-update project drift was not fenced rc=$rc out=$preflight_out log=$(cat "$log")"

reset_generic_item
preflight_out="$(PATH="$tmp/bin:$PATH" A1_GENERIC_FIXTURES=1 A1_MUTATE_ON_UPDATE=1 \
  A1_DRIFT_AFTER_UPDATE=type \
  bash "$root/bootstrap/aone-fields.sh" preflight 84608993 1086837 2>/dev/null)"; rc=$?
[ "$rc" -eq 3 ] && printf '%s' "$preflight_out" | jq -e '
  .failureReason=="context_drift_after_readback" and .filled==true
' >/dev/null && [ "$(grep -c 'workitem update' "$log" || true)" -eq 1 ] \
  && ok "workitem type drift after update fails closed" \
  || bad "post-update type drift was not fenced rc=$rc out=$preflight_out log=$(cat "$log")"

reset_generic_item
preflight_out="$(PATH="$tmp/bin:$PATH" A1_GENERIC_FIXTURES=1 A1_MUTATE_ON_UPDATE=1 \
  A1_CONFLICT_AFTER_UPDATE=140282 \
  bash "$root/bootstrap/aone-fields.sh" preflight 84608993 1086837 2>/dev/null)"; rc=$?
[ "$rc" -eq 3 ] && printf '%s' "$preflight_out" | jq -e '
  .failureReason=="assignment_conflict_after_readback" and .filled==true and
  any(.unresolved[]; .id=="140282" and .expected=="defined" and .actual=="manual")
' >/dev/null && [ "$(grep -c 'workitem update' "$log" || true)" -eq 1 ] \
  && ok "concurrent legal non-empty assignment change fails closed" \
  || bad "assignment conflict was not fenced rc=$rc out=$preflight_out log=$(cat "$log")"

# No configured policy and multiple candidates remains unresolved.
cat > "$tmp/workitem.json" <<'JSON'
{"title":"Provider generic issue","fields":[
  {"identifier":"workitemType","value":"36","displayValue":"需求问题"},
  {"identifier":"space","value":"528766","displayValue":"provider"},
  {"identifier":"140097","value":"","displayValue":""}
]}
JSON
cat > "$tmp/fields.json" <<'JSON'
[
  {"identifier":"140097","displayName":"涉及云产品","isRequired":true,"sourceType":"team","format":"list","options":[]}
]
JSON
rm -f "$A1_TEST_STATE.updated" "$A1_TEST_STATE.get_count"
: > "$log"
preflight_out="$(PATH="$tmp/bin:$PATH" A1_AUTO_FIXTURES=1 \
  bash "$root/bootstrap/aone-fields.sh" preflight 9002 528766 2>/dev/null)"; rc=$?
[ "$rc" -eq 3 ] && printf '%s' "$preflight_out" | jq -e '
  .failureReason=="unresolved_required_fields" and
  .unresolved[0].reason=="no_title_match" and .filled==false
' >/dev/null && [ "$(grep -c 'workitem update' "$log" || true)" -eq 0 ] \
  && ok "ambiguous required field remains unresolved" \
  || bad "unresolved contract failed rc=$rc out=$preflight_out log=$(cat "$log")"

inspect_out="$(PATH="$tmp/bin:$PATH" A1_AUTO_FIXTURES=1 \
  bash "$root/bootstrap/aone-fields.sh" inspect 9002 528766 2>/dev/null)"; rc=$?
[ "$rc" -eq 0 ] && printf '%s' "$inspect_out" | jq -e '
  .status=="repair_required" and .workitemId=="9002" and .project=="528766" and
  (.missing | length)==1 and .assignments==[] and
  .unresolved[0].id=="140097" and (.unresolved[0].options | length)==2
' >/dev/null && [ "$(grep -c 'workitem update' "$log" || true)" -eq 0 ] \
  && ok "inspect returns deterministic and unresolved repair context without mutation" \
  || bad "inspect contract failed rc=$rc out=$inspect_out log=$(cat "$log")"
inspect_canonical="$(printf '%s' "$inspect_out" | jq -S -c '{
  missing:(.missing // []),assignments:(.assignments // []),
  unresolved:(.unresolved // [])
}')"
inspect_digest="$(printf '%s' "$inspect_canonical" | shasum -a 256 \
  | awk '{print substr($1,1,24)}')"
inspect_type="$(printf '%s' "$inspect_out" | jq -r '.workitemType')"
inspect_revision="$(printf '%s' "$inspect_out" | jq -r '.revision')"

: > "$log"
PATH="$tmp/bin:$PATH" A1_AUTO_FIXTURES=1 A1_MUTATE_ON_UPDATE=1 \
  bash "$root/bootstrap/aone-fields.sh" apply 9002 528766 \
  "$inspect_type" "$inspect_revision" 000000000000000000000000 \
  '140097=mongodb' >/dev/null 2>&1; rc=$?
[ "$rc" -eq 3 ] && [ "$(grep -c 'workitem update' "$log" || true)" -eq 0 ] \
  && ok "apply snapshot drift fails closed before mutation" \
  || bad "snapshot drift was not fenced rc=$rc log=$(cat "$log")"

: > "$log"
PATH="$tmp/bin:$PATH" A1_AUTO_FIXTURES=1 A1_MUTATE_ON_UPDATE=1 \
  bash "$root/bootstrap/aone-fields.sh" apply 9002 528766 \
  "$inspect_type" "$inspect_revision" "$inspect_digest" \
  '140097=not-a-candidate' >/dev/null 2>&1; rc=$?
[ "$rc" -eq 4 ] && [ "$(grep -c 'workitem update' "$log" || true)" -eq 0 ] \
  && ok "apply rejects values outside the fresh candidate set before mutation" \
  || bad "illegal apply was not rejected rc=$rc log=$(cat "$log")"

: > "$log"
PATH="$tmp/bin:$PATH" A1_AUTO_FIXTURES=1 A1_MUTATE_ON_UPDATE=1 \
  A1_CONFLICT_AFTER_UPDATE=140097 \
  bash "$root/bootstrap/aone-fields.sh" apply 9002 528766 \
  "$inspect_type" "$inspect_revision" "$inspect_digest" \
  '140097=mongodb' >"$tmp/apply-mismatch.out" 2>/dev/null; rc=$?
[ "$rc" -eq 3 ] && printf '%s' "$(cat "$tmp/apply-mismatch.out")" | jq -e '
  .errorType=="field_apply_readback_mismatch" and
  .readback==[{"id":"140097","value":"redis"}] and
  .mismatches==[{"id":"140097","expected":"mongodb","actual":"redis"}]
' >/dev/null && [ "$(grep -c 'workitem update' "$log" || true)" -eq 1 ] \
  && ok "apply reads and rejects conflicting canonical values" \
  || bad "canonical mismatch was not detected rc=$rc out=$(cat "$tmp/apply-mismatch.out") log=$(cat "$log")"

cat > "$tmp/workitem.json" <<'JSON'
{"title":"Provider generic issue","fields":[
  {"identifier":"workitemType","value":"36","displayValue":"需求问题"},
  {"identifier":"space","value":"528766","displayValue":"provider"},
  {"identifier":"140097","value":"","displayValue":""}
]}
JSON
rm -f "$A1_TEST_STATE.updated" "$A1_TEST_STATE.get_count"
: > "$log"
apply_out="$(PATH="$tmp/bin:$PATH" A1_AUTO_FIXTURES=1 A1_MUTATE_ON_UPDATE=1 \
  bash "$root/bootstrap/aone-fields.sh" apply 9002 528766 \
  "$inspect_type" "$inspect_revision" "$inspect_digest" \
  '140097=mongodb' 2>"$tmp/apply.err")"; rc=$?
[ "$rc" -eq 0 ] && printf '%s' "$apply_out" | jq -e '
  .status=="ready" and .filled==true and .missing==[] and .unresolved==[] and
  .assignments==[{"id":"140097","value":"mongodb"}] and
  .readback==[{"id":"140097","value":"mongodb"}]
' >/dev/null && [ "$(grep -c 'workitem update' "$log" || true)" -eq 1 ] \
  && ok "apply validates, updates once, and returns an empty readback" \
  || bad "legal apply failed rc=$rc out=$apply_out err=$(cat "$tmp/apply.err") log=$(cat "$log")"

# --- configured default outranks a unique title match -----------------------
# tf_provider pins 107239 归属产品 in config/pools.json. The generic 107239
# option set also contains "云数据库 MongoDB 版", which a MongoDB-flavoured
# title matches uniquely — the exact shape that made a real Terraform ticket
# resolve to an unrelated product before the pinned default was given priority.
cat > "$tmp/workitem.json" <<'JSON'
{"title":"alicloud mongodb sharding instance apply failure","fields":[
  {"identifier":"workitemType","value":"36","displayValue":"需求问题"},
  {"identifier":"space","value":"528766","displayValue":"provider"},
  {"identifier":"107239","value":"","displayValue":""}
]}
JSON
cat > "$tmp/fields.json" <<'JSON'
[
  {"identifier":"107239","displayName":"归属产品","isRequired":true,"sourceType":"team","format":"plugin","options":[]}
]
JSON
rm -f "$A1_TEST_STATE.updated" "$A1_TEST_STATE.get_count"
: > "$log"
inspect_out="$(PATH="$tmp/bin:$PATH" A1_GENERIC_FIXTURES=1 \
  bash "$root/bootstrap/aone-fields.sh" inspect 9002 528766 2>/dev/null)"; rc=$?
[ "$rc" -eq 0 ] && printf '%s' "$inspect_out" | jq -e '
  .status=="repair_required" and .unresolved==[] and
  (.assignments | length)==1 and
  .assignments[0].id=="107239" and .assignments[0].value=="906688" and
  .assignments[0].source=="configured_default"
' >/dev/null && [ "$(grep -c 'workitem update' "$log" || true)" -eq 0 ] \
  && ok "pinned pool default beats a unique title match" \
  || bad "configured default precedence failed rc=$rc out=$inspect_out log=$(cat "$log")"

# --- field metadata cache ---------------------------------------------------
# Same inputs twice: the second pass must serve field list/options from cache
# and issue no further metadata reads, while still reading the workitem itself
# live (apply's revision/digest fence depends on that never being cached).
cache_probe="$tmp/cache-probe"
rm -rf "$cache_probe"
rm -f "$A1_TEST_STATE.updated" "$A1_TEST_STATE.get_count"
: > "$log"
PATH="$tmp/bin:$PATH" A1_GENERIC_FIXTURES=1 JARVIS_CACHE_DIR="$cache_probe" \
  JARVIS_FIELD_META_TTL=900 bash "$root/bootstrap/aone-fields.sh" \
  missing 9002 >/dev/null 2>&1
cold_list="$(grep -c 'field list' "$log" || true)"
cold_opts="$(grep -c 'field options 107239' "$log" || true)"
: > "$log"
PATH="$tmp/bin:$PATH" A1_GENERIC_FIXTURES=1 JARVIS_CACHE_DIR="$cache_probe" \
  JARVIS_FIELD_META_TTL=900 bash "$root/bootstrap/aone-fields.sh" \
  missing 9002 >/dev/null 2>&1
warm_list="$(grep -c 'field list' "$log" || true)"
warm_opts="$(grep -c 'field options 107239' "$log" || true)"
warm_get="$(grep -c 'workitem get' "$log" || true)"
[ "$cold_list" -eq 1 ] && [ "$cold_opts" -eq 1 ] \
  && [ "$warm_list" -eq 0 ] && [ "$warm_opts" -eq 0 ] && [ "$warm_get" -eq 1 ] \
  && ok "warm metadata cache drops field reads but still reads the workitem live" \
  || bad "cache contract failed cold(list=$cold_list opts=$cold_opts) warm(list=$warm_list opts=$warm_opts get=$warm_get)"

# TTL=0 must fully bypass the cache even with a warm directory present.
: > "$log"
PATH="$tmp/bin:$PATH" A1_GENERIC_FIXTURES=1 JARVIS_CACHE_DIR="$cache_probe" \
  JARVIS_FIELD_META_TTL=0 bash "$root/bootstrap/aone-fields.sh" \
  missing 9002 >/dev/null 2>&1
[ "$(grep -c 'field list' "$log" || true)" -eq 1 ] \
  && ok "JARVIS_FIELD_META_TTL=0 bypasses a warm cache" \
  || bad "TTL=0 did not bypass cache: $(cat "$log")"

echo "Results: $PASS passed, $FAIL failed"
[ "$FAIL" -eq 0 ]
