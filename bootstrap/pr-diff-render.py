#!/usr/bin/env python3
"""
pr-diff-render.py — render a GitHub PR diff (gh pr diff output) into the same
polished, localized (Chinese) HTML page as provider-gen-diff-render.py, but for
the PR-head-vs-master view (no generator baseline; for hand-written resources
that are not modeled in CloudSpec).

Reuses the generic unified-diff rendering (split + unified views, dual theme,
no-wrap split columns) by importing provider-gen-diff-render.py as a module.

USAGE
  python3 pr-diff-render.py <pr-diff.txt> <output.html> [--pr N] [--title T]
  gh pr diff <N> -R <repo> | python3 pr-diff-render.py /dev/stdin out.html --pr <N>
"""
import importlib.util
import os
import re
import sys

# Import the generic diff rendering (CSS, parse_hunks, split/unified rows &
# renderers, esc, ANSI_RE) from the sibling gen-diff renderer. Its main() is
# guarded by __name__ so importing it does not execute it.
_HERE = os.path.dirname(os.path.abspath(__file__))
_spec = importlib.util.spec_from_file_location(
    "_gen_diff_render",
    os.path.join(_HERE, "provider-gen-diff-render.py"))
g = importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(g)


def esc(s):
    return g.esc(s)


# ── classify file paths into section tags ─────────────────────────────────────
def classify(path):
    if path.startswith('alicloud/resource_alicloud_'):
        return 'resource'
    if path.startswith('alicloud/data_source_alicloud_'):
        return 'datasource'
    if path.startswith('alicloud/service_'):
        return 'service'
    if path.endswith('_test.go'):
        return 'test'
    if path.startswith('website/docs/r/'):
        return 'doc (r)'
    if path.startswith('website/docs/d/'):
        return 'doc (d)'
    if path == 'alicloud/provider.go':
        return 'provider'
    return 'other'


TAG_CN = {
    'resource': '资源代码',
    'datasource': '数据源代码',
    'service': '服务封装',
    'test': '测试用例',
    'doc (r)': '资源文档',
    'doc (d)': '数据源文档',
    'provider': 'provider 注册',
    'other': '其它',
}


def section_tag(label):
    key = label.split(':', 1)[0].strip()
    return TAG_CN.get(key, key)


def section_file(label):
    return g.section_file(label)


# ── parse a multi-file unified diff (gh pr diff output) ───────────────────────
def parse_pr_diff(text):
    lines = [g.ANSI_RE.sub('', l) for l in text.split('\n')]
    sections = []
    i, n = 0, len(lines)
    while i < n:
        if lines[i].startswith('diff --git'):
            start = i
            i += 1
            while i < n and not lines[i].startswith('diff --git'):
                i += 1
            body = lines[start:i]
            path = ''
            for l in body:
                if l.startswith('+++ b/'):
                    path = l[6:]
                    break
                if l.startswith('+++ '):
                    path = l[4:].strip()
                    break
            if not path:  # fall back to the diff --git b/<path> token
                m = re.match(r'^diff --git a/\S+ b/(.+)$', body[0])
                if m:
                    path = m.group(1)
            tag = classify(path)
            sections.append({'label': f'{tag}: {path}', 'body': body, 'path': path})
        else:
            i += 1
    return sections


# ── render one file panel (reuses generic split/unified renderers) ───────────
def render_panel(section):
    label = section['label']
    tag = section_tag(label)
    fname = section_file(label)
    pre, hunks = g.parse_hunks(section['body'])
    no_diff = not hunks and not any(l.startswith('@@') for l in section['body'])
    if no_diff:
        body_html = '<div class="empty">无 hunk（仅二进制 / 元数据变更）。</div>'
    else:
        u = g.render_unified(g.unified_rows(hunks))
        s = g.render_split(g.split_rows(hunks))
        body_html = (
            f'<div class="view-split">{s}</div>'
            f'<div class="view-unified" hidden><pre class="diff-pre unified">{u}</pre></div>'
        )
    return (
        f'<section class="diff-panel"><div class="diff-head">'
        f'<span class="tag">{esc(tag)}</span><span class="path">{esc(fname)}</span></div>'
        f'<div class="diff-body">{body_html}</div></section>'
    )


def main():
    argv = sys.argv[1:]
    pr_no = ''
    title = ''
    files_extra = []
    src = out = ''
    pos = []
    while argv:
        a = argv.pop(0)
        if a == '--pr':
            pr_no = argv.pop(0)
        elif a == '--title':
            title = argv.pop(0)
        elif a == '--meta':
            files_extra.append(argv.pop(0))
        elif not src:
            src = a
        elif not out:
            out = a
        else:
            pos.append(a)
    if not src or not out:
        sys.stderr.write('usage: pr-diff-render.py <pr-diff.txt> <output.html> [--pr N] [--title T] [--meta "k: v"]\n')
        sys.exit(2)
    with open(src, encoding='utf-8', errors='replace') as f:
        text = f.read()
    sections = parse_pr_diff(text)

    added = removed = 0
    for s in sections:
        for l in s['body']:
            if l.startswith('+') and not l.startswith('+++'):
                added += 1
            elif l.startswith('-') and not l.startswith('---'):
                removed += 1

    panels = ''.join(render_panel(s) for s in sections)

    # extra meta lines passed on the CLI (--meta "head: <branch>" ...)
    extra = []
    for ml in files_extra:
        m = re.match(r'\s*([^:]+?)\s*:\s*(.*)', ml)
        if m:
            extra.append((m.group(1).strip(), m.group(2).strip()))

    meta_rows = [
        ('PR', pr_no or '—'),
        ('base', 'master'),
        ('文件数', str(len(sections))),
        ('新增 / 删除', f'+{added} / −{removed}'),
    ] + extra
    meta_html = ''.join(f'<div class="mk">{esc(k)}</div><div class="mv">{esc(v)}</div>' for k, v in meta_rows)
    meta_d = dict(extra)

    h1 = esc(title) if title else (
        esc(os.path.basename(next((s['path'] for s in sections if classify(s['path']) == 'resource'), sections[0]['path'] if sections else '')).replace('resource_alicloud_', '').replace('.go', ''))
        if sections else 'PR diff')
    pr_pop = f'#{esc(pr_no)}' if pr_no else ''
    head = meta_d.get('head 分支', '')
    sub = (f'PR {esc(pr_no)}（head 分支 <code>{esc(head)}</code>）相对 '
           f'master 的完整改动。该资源为<b>纯手写</b>（RealtimeCompute 系列，走 '
           f'<code>foasconsole</code> API，不在 generator / CloudSpec 建模范围），'
           f'无生成器基线可对比，故本页展示 PR 整体新增代码，而非「手写 vs 生成」。'
           f'并排视图左侧为改动前（master），右侧为改动后（PR head）。')

    PAGE = f"""<title>PR {esc(pr_no)} 改动总览</title>
<link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@400;500;600&family=IBM+Plex+Sans:wght@400;500;600;700&family=Noto+Sans+SC:wght@400;500;600&display=swap">
<style>{g.CSS}</style>
<div class="wrap">
  <header class="page">
    <p class="eyebrow">PR diff · head vs master · 纯手写资源</p>
    <h1>{h1}<span class="pop">{pr_pop}</span></h1>
    <p class="sub">{sub}</p>
    <div class="meta-grid">{meta_html}</div>
    <div class="bar">
      <div class="legend">
        <span><span class="sw del"></span>改动前（master）</span>
        <span><span class="sw add"></span>改动后（PR head）</span>
        <span><span class="sw hunk"></span>hunk 头</span>
      </div>
      <div class="toggle" role="group" aria-label="diff 视图切换">
        <button id="t-split" class="active" type="button">并排</button>
        <button id="t-unified" type="button">统一</button>
      </div>
    </div>
  </header>
  <div class="stats">
    <div class="tile"><div class="num">{len(sections)}</div><div class="lab">文件区段</div></div>
    <div class="tile"><div class="num add">+{added}</div><div class="lab">新增行</div></div>
    <div class="tile"><div class="num del">−{removed}</div><div class="lab">删除行</div></div>
    <div class="tile"><div class="num">{sum(1 for s in sections if classify(s['path']) in ('resource','datasource','service','test'))}</div><div class="lab">手写代码文件</div></div>
  </div>
  {panels}
  <div class="footer">数据源：<code>gh pr diff {esc(pr_no)}</code>（head 分支 vs master）。该 resource 未在 CloudSpec 建模，<code>provider-gen-diff.sh</code> 对其 degrade 为 checklist-only（无生成基线），故此处改用 PR 整体 diff 渲染。代码单行不换行，长行横向滚动。</div>
</div>
<script>
(function(){{
  var b=document.body, s=document.getElementById('t-split'), u=document.getElementById('t-unified');
  function set(view){{
    b.classList.remove('view-split','view-unified'); b.classList.add('view-'+view);
    s.classList.toggle('active',view==='split'); u.classList.toggle('active',view==='unified');
    try{{localStorage.setItem('gen-diff-view',view);}}catch(e){{}}
  }}
  s.addEventListener('click',function(){{set('split');}});
  u.addEventListener('click',function(){{set('unified');}});
  var saved;try{{saved=localStorage.getItem('gen-diff-view');}}catch(e){{}}
  if(saved==='unified'){{set('unified');}}
}})();
</script>"""
    with open(out, 'w', encoding='utf-8') as f:
        f.write(PAGE)
    print(f'rendered {out} (+{added}/-{removed}, {len(sections)} sections)')


if __name__ == '__main__':
    main()
