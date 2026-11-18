#!/usr/bin/env python3
"""Explicit successful skip for the temporarily disabled automatic AI review."""
import json
import os
from pathlib import Path


def main():
    message = 'AI 评审已关闭，本次未执行；按配置默认通过，不代表 AI 评审已完成。'
    report = {'status': 'skipped', 'conclusion': 'success', 'reason': 'ai-review-disabled',
              'summary': message, 'headSha': os.environ.get('ACUBE_HEAD', ''),
              'mr': os.environ.get('ACUBE_MR', ''), 'runId': os.environ.get('ACUBE_CI_RUN', '')}
    Path(os.environ.get('AONE_CI_SUMMARY_MD', 'review-summary.md')).write_text(
        '# ACube AI 评审 · 已关闭\n\n' + message + '\n', encoding='utf-8')
    Path('review-report.json').write_text(json.dumps(report, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
    print(message, flush=True)
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
