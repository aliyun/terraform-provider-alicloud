import unittest

from document_scope import select


class DocumentationScopeTest(unittest.TestCase):
    def test_helper_change_runs_go_checks_without_document_tools(self):
        paths = ['alicloud/helper/retry.go', 'alicloud/helper/retry_test.go']
        self.assertEqual(select(paths, False), {'docs': False, 'provider': True, 'framework': True})

    def test_readme_only_has_no_applicable_jobs(self):
        self.assertFalse(any(select(['README.md'], False).values()))

    def test_removed_document_still_selects_doc_and_consistency_checks(self):
        self.assertEqual(select(['website/docs/r/removed.html.markdown'], False),
                         {'docs': True, 'provider': True, 'framework': False})

    def test_workflow_change_and_full_runs_select_all(self):
        self.assertTrue(all(select(['.aoneci/document-check.yml'], False).values()))
        self.assertTrue(all(select(['README.md'], True).values()))
