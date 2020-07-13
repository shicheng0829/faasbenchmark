import os
import json
import csv
import oss2
import pandas as pd


def getpath(dir):
    files = os.listdir(dir)
    files = sorted(files)
    for i in range(len(files)-1, 0, -1):
        if "results" in files[i]:
            return files[i]

def convertToHtml(result,title):
    df = pd.DataFrame(columns=title)
    for row in result:
        df = df.append(pd.DataFrame([row], columns=title))
    return df.to_html(index=False)




path = getpath("./")
files = os.listdir(path)
result = []
with open("result.csv", "w") as csvfile:
    writer = csv.writer(csvfile)
    writer.writerow(["testNmae", "invocationOverhead", "duration", "responseTime", "reusedRate", "failedRate"])


for file in files:
    with open(os.path.join(path, file, "aliyun", "result.json")) as load_f:
        load_dict = json.load(load_f)
        list = load_dict["functions"][0]['results']
        sumofInvocationOverhead = 0
        sumofDuration = 0
        sumofResponseTime = 0
        failed = 0
        reused = 0
        for dic in list:
            if dic['failed'] is False:
                sumofInvocationOverhead += dic['invocationOverhead']
                sumofDuration += dic['duration']
                sumofResponseTime += dic['responseTime']
                if dic['reused'] is True:
                    reused += 1
            else:
                failed += 1
        non_failed_num = len(list) - failed
        if non_failed_num == 0:
            invocationOverhead = 0
            duration = 0
            responseTime = 0
            reusedRate = 0
            failedRate = 100
        else:
            invocationOverhead = sumofInvocationOverhead / non_failed_num
            duration = sumofDuration / non_failed_num
            responseTime = sumofResponseTime / non_failed_num
            reusedRate = reused / non_failed_num*100
            failedRate = failed / len(list)*100
        with open("result.csv", "a") as csvfile:
            writer = csv.writer(csvfile)
            writer.writerow([file, invocationOverhead, duration, responseTime, reusedRate, failedRate])
            result.append([file, invocationOverhead, duration, responseTime, reusedRate, failedRate])

with open("result.html","w") as htmlfile:
    title = ["testNmae", "invocationOverhead", "duration", "responseTime", "reusedRate", "failedRate"]
    htmlfile.write(convertToHtml(result, title))
access_key_id = os.getenv("ACCESS_KEY_ID")
access_key_secret = os.getenv("ACCESS_KEY_SECRET")
oss_endpoint = os.getenv("OSS_ENDPOINT")
bucket_name = os.getenv("BUCKET_NAME")

auth = oss2.Auth(access_key_id, access_key_secret)
bucket = oss2.Bucket(auth, oss_endpoint, bucket_name)
# write history result csv
bucket.put_object_from_file(os.path.join(path, "result.csv"), "result.csv")
# write history log file
bucket.put_object_from_file(os.path.join(path, "log"), "log")
# update current result html
bucket.put_object_from_file("result.html", "result.html")
# update current result csv
bucket.put_object_from_file("result.csv", "result.csv")
# update current result log
bucket.put_object_from_file("log", "log")

os.remove("./result.csv")
os.remove("./log")
os.remove("./result.html")
