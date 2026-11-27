"""Select Aone Documentation Checks jobs after validating the MR comparison."""
import os
from pathlib import Path

from check_scope import changed_paths, provider_needed


def select(paths, full):
    shared = full or any(
        p in {'go.mod', 'go.sum', '.go-version', 'GNUmakefile', 'scripts/basic-check.sh'}
        or p.startswith(('.aoneci/', 'scripts/ci/', 'scripts/release/', 'scripts/internal/',
                         'scripts/collection-order/'))
        for p in paths
    )
    docs = shared or any(
        p.startswith(('website/docs/', 'scripts/document/', '.github/markdownlint/'))
        or p == '.markdown-link-check.json' for p in paths
    )
    return {
        'docs': docs,
        'provider': shared or provider_needed(paths, 'consistency'),
        'framework': shared or any(p.endswith('.go') for p in paths),
    }


if __name__ == '__main__':
    scope = select(changed_paths(), os.environ['DOCUMENT_CHECK_FULL'] == 'true')
    labels = {'docs': 'Markdown, formatting and content', 'provider': 'basic and consistency checks',
              'framework': 'Go test framework import guard'}
    lines = [f"{labels[group]}: {'selected' if needed else 'not applicable'}" for group, needed in scope.items()]
    if not any(scope.values()):
        lines.append('No applicable checks for this change. Documentation Checks passes.')
    for group, needed in scope.items():
        (Path(os.environ['AONE_CI_OUTPUTS']) / group).write_text('yes' if needed else 'no')
    print('\n'.join(lines))
    if os.environ.get('AONE_CI_SUMMARY_MD'):
        with open(os.environ['AONE_CI_SUMMARY_MD'], 'a') as summary:
            summary.writelines(f'- {line}\n' for line in lines)
