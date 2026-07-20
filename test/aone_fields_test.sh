#!/usr/bin/env bash
# Hermetic tests for bootstrap/aone-fields.sh. All workitem/field data is synthetic.

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root="$(cd "$test_dir/.." && pwd)"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
mkdir -p "$tmp/bin"
log="$tmp/a1.log"

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
  cat "$A1_TEST_WI"; exit 0
fi
if [ "$1 $2 $3 $4" = "project workitem field list" ]; then
  cat "$A1_TEST_FIELDS"; exit 0
fi
if [ "$1 $2 $3 $4" = "project workitem field options" ]; then
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
  exit 0
fi
exit 9
STUB
chmod +x "$tmp/bin/a1"

export A1_TEST_LOG="$log" A1_TEST_WI="$tmp/workitem.json" A1_TEST_FIELDS="$tmp/fields.json"
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
  && grep -q -- '--cfs 140282=有OpenAPI，资源未定义，开放平台维护' "$log" \
  && grep -q -- '--cfs 107239=891438' "$log" \
  && ok "auto-fill submits every required field in one update" \
  || bad "auto-fill assignments missing: $(cat "$log")"

echo "Results: $PASS passed, $FAIL failed"
[ "$FAIL" -eq 0 ]
