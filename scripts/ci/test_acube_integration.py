import importlib.util
import tempfile
from pathlib import Path
import unittest
from unittest.mock import patch

spec = importlib.util.spec_from_file_location('integration', Path(__file__).with_name('acube-integration.py'))
m = importlib.util.module_from_spec(spec); spec.loader.exec_module(m)

class IntegrationTest(unittest.TestCase):
    def request(self):
        return dict(mode='acctest',sourceId='a'*64,checkBaseSourceId='b'*64,commitSha='c'*40,
                    baseCommitSha='d'*40,variables={'TEST':'alicloud/service/rocketmq'},testCases=['TestAccOne','TestAccTwo'])
    def result(self):
        r=self.request()
        return dict(status='success', result=dict(r,exitCode=0,testDirectory=r['variables']['TEST'],tests={
            'cases':[{'name':n,'result':'PASS'} for n in r['testCases']]}))
    def test_missing_skipped_duplicate_and_foreign_results_fail(self):
        m.validate_result(self.result(), self.request())
        for defect in ['skip','missing','duplicate','source','directory']:
            r=self.result()
            if defect=='skip': r['result']['tests']['cases'][0]['result']='SKIP'
            if defect=='missing': r['result']['tests']['cases'].pop()
            if defect=='duplicate': r['result']['tests']['cases'][1]=r['result']['tests']['cases'][0]
            if defect=='source': r['result']['sourceId']='wrong'
            if defect=='directory': r['result']['testDirectory']='alicloud'
            with self.subTest(defect=defect),self.assertRaises(ValueError): m.validate_result(r,self.request())
    def test_plan_deduplicates_complete_package_cases_and_preserves_example(self):
        def git(*args):
            if args[0]=='diff': return 'alicloud/resource_alicloud_demo.go\nalicloud/resource_alicloud_demo_test.go\nwebsite/docs/r/demo.html.markdown'
            if args[0]=='ls-tree': return ''
            if args[0]=='show':
                return 'func TestAccOne(t *testing.T) {\nfunc TestAccTwo(t *testing.T) {' if args[-1].endswith('_test.go') else '```terraform\nresource "alicloud_demo" "test" {}\n```'
            raise AssertionError(args)
        with patch.object(m,'git',side_effect=git): p=m.plan('a'*40,'b'*40)
        self.assertEqual(p['packages'],{'alicloud':['TestAccOne','TestAccTwo']})
        self.assertEqual(len(p['examples']),1)
        self.assertIn('source = "aliyun/alicloud"',p['examples'][0]['mainTf'])
        self.assertEqual(len(p['examples'][0]['originalSha256']),64)
    def test_same_run_and_plan_recover_same_key_across_processes(self):
        request=self.request()
        frozen={'head':request['commitSha'],'base':request['baseCommitSha'],'packages':{'alicloud/service/rocketmq':['TestAccOne','TestAccTwo']},'examples':[]}
        keys=[]
        class Client:
            def call(_,method,path,body=None):
                if path=='/capabilities': return {'submissionIdempotency':'oss-admission-v1','hclActions':['example']}
                if path=='/sources/import': return {'sourceId':request['sourceId'] if body['commitSha']==request['commitSha'] else request['checkBaseSourceId']}
                if path=='': keys.append(body['idempotencyKey']); return {'taskId':123,'reused':len(keys)>1}
                if path=='/123/logs': return {}
                if path=='/123': return self.result()
                raise AssertionError(path)
        with tempfile.TemporaryDirectory() as root:
            for _ in range(2): m.execute(Client(),frozen,'acctest','ci-123',Path(root))
        self.assertEqual(keys[0],keys[1])
    def test_failed_task_does_not_resubmit(self):
        calls=[]
        class Client:
            def call(_,method,path,body=None):
                calls.append((method,path))
                if path=='/capabilities': return {'submissionIdempotency':'oss-admission-v1','hclActions':['example']}
                if path=='/sources/import': return {'sourceId':'a'*64}
                if path=='': return {'taskId':123}
                if path=='/123': return {'status':'failed'}
                if path=='/123/logs': return {}
        with tempfile.TemporaryDirectory() as root,self.assertRaises(ValueError):
            m.execute(Client(),{'head':'a'*40,'base':'b'*40,'packages':{'alicloud':['TestAccOne']},'examples':[]},'acctest','run-123',Path(root))
        self.assertEqual(calls.count(('POST','')),1)

if __name__=='__main__': unittest.main()
