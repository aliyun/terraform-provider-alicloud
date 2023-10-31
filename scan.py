# -- coding: utf-8 --**
import os
import re
import datetime
import subprocess
import time
import urllib.parse
import json
import requests


def des():
    '''


    本地got-hook插件
    敏感信息扫描,本地运行部分。
    这部分是放在云端，供使用者本地的脚本拉取。
    commit 时执行

    :return:
    '''


# 测试地址
# root_url = "http://127.0.0.1:8000/"


# 线上地址
root_url = "https://sea.aliyun-inc.com/"

def print_time():
    """
    获取当前时间，返回一个字符串，格式为：%Y-%m-%d %H:%M:%S
    """
    now = datetime.datetime.now()
    return now.strftime("%Y-%m-%d %H:%M:%S")


# 判断系统是否存在某个命令
def check_command_exists(command):
    try:
        output = subprocess.check_output(f"which {command}", shell=True, stderr=subprocess.STDOUT)
        return True
    except subprocess.CalledProcessError:
        return False


# 获取当前路径
def get_current_path():
    return os.getcwd()


def get_hook_path():
    # 执行 Git 命令获取钩子文件的路径
    process = subprocess.Popen(['git', 'config', 'core.hooksPath'], stdout=subprocess.PIPE)
    output, _ = process.communicate()
    hook_path = output.decode('utf-8').strip()
    # print(hook_path)
    return hook_path


owner_email = ""
owner_num1 = ""
no_file_path = get_hook_path() + "/my_no.txt"
if os.path.exists(no_file_path):
    with open(no_file_path, 'r') as f:
        content = f.readlines()
        non_empty_lines = [line.strip() for line in content if line.strip()]
        owner_num1 = non_empty_lines[0]

# 判断一下是否存在系统命令
elif check_command_exists("security"):
    # 获取本机属于哪个员工/获取邮箱
    cre_res = os.popen('security dump-keychain')
    cre_res_email = cre_res.read()
    email_list = re.findall(r"[A-Za-z]+[A-Za-z0-9.\-+_]+@alibaba-inc.com", cre_res_email)
    email_list = list(set(email_list))
    for email in email_list:
        if email != "alilang@alibaba-inc.com" and email != "cloudshell@alibaba-inc.com":
            owner_email = email

    # 获取本机的研发工号

    wb_owner_num_list = re.findall(r'0x00000001 <blob>="(.*)"', cre_res_email)
    # print(wb_owner_num_list)
    owner_num_list = list(set(wb_owner_num_list))
    # print(owner_num_list)

    for owner_unm in owner_num_list:
        wb_owner_num = re.findall(r'[W][B][0-9]{3,8}', owner_unm)
        # print(wb_owner_num)
        owner_num = re.findall(r'[0-9]{3,6}', owner_unm)
        # print(owner_num)
        if wb_owner_num == [] and owner_num == []:
            pass
        elif wb_owner_num == [] and owner_num != []:
            owner_num1 = owner_num[0]
            # print(owner_num)
        else:
            owner_num1 = wb_owner_num[0]
else:
    print(print_time() + "[!] 获取工号失败！")
# 获取项目地址
project_url_res = os.popen('git ls-remote --get-url')
project_url = project_url_res.read()

# 获取提交人信息
create_people_res = os.popen('git show -s --format=%an')
create_people = create_people_res.read()
local_file_name = 'diff_sec' + str(time.time()) + '.zip'
# 打包项目代码
os.popen('git diff --cached > diff_sec.txt && zip -q ' + local_file_name + ' diff_sec.txt && rm -rf diff_sec.txt')

time.sleep(2)

# 如果更改服务器 这里的地址需要同步修改
url = root_url + 'get_sign?tmp_name=' + local_file_name

req = requests.get(url)
json_data = req.json()
# print(json_data)
file_url = ""
if req.status_code == 200:
    # proxies = {
    #     "http": "http://127.0.0.1:8080",
    #     "https": "http://127.0.0.1:8080",
    # }
    oss_file_name = local_file_name
    try:
        # 上传文件
        with open(local_file_name, 'rb') as f:
            response = requests.put(json_data['upload_url'], data=f)
        if response.status_code == 200:
            # 下载地址传递给yy扫描
            yy_scan_url = "http://yy.aliyun-inc.com/api/check/AliyunCommonScan.json?source=fuchen&auth=tg6edtvaaer4dfjisdf6sdfkjefiuh3dfsYaGxdlsrkzLWJuqhF&priority=18&ruleId=43&repoPath=" + urllib.parse.quote(
                json_data['down_url'])
            yy_scan_res = requests.get(yy_scan_url)
            yy_res = yy_scan_res.json()
            # print(yy_scan_res.text)
            # yy扫描任务id
            scanTaskId = yy_res['scanTasks'][0]['scanTaskId']

            # 将yy任务id上传到平台
            post_url = root_url + "rec_commit?token=uyjejuhuhuwe787826384hkbnkjfbn" + "&project_url=" + \
                       project_url + "&create_people=" + create_people + "&scanTaskId=" + str(
                scanTaskId) + "&owner_email=" + owner_email + "&owner_num=" + owner_num1

            post_res = requests.get(post_url)

        else:
            print(print_time() + "[!]文件上传异常！联系抚尘")
    except:
        print(print_time() + "[!]文件上传异常！联系抚尘")

# 删除生成的diff压缩包
os.system('rm -rf ' + local_file_name)
