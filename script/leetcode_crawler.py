import requests
import json
import time
import jsonpath
from typing import List


class LcCrawler():
    def __init__(self, csrftoken: str):
        self.cookie = 'csrftoken=%s;' % (csrftoken,)

    def fetch_problem_list(self):

        def fetch(offset: int, limit: int) -> List:
            data = self.fetch_question_list_by_range(offset, limit)

            hasMore = jsonpath.jsonpath(data, '$.data.problemsetQuestionList.hasMore')[0]
            problems = jsonpath.jsonpath(data, '$.data.problemsetQuestionList.questions')[0]
            if hasMore:
                problems.extend(fetch(offset + limit, limit))
                return problems
            return problems if problems else []

        return fetch(0, 50)

    def fetch_question_list_by_range(self, offset: int, limit: int) -> dict:
        data = {
            'query': '''
                query problemsetQuestionList($categorySlug: String, $limit: Int, $skip: Int, $filters: QuestionListFilterInput) {
                  problemsetQuestionList(
                    categorySlug: $categorySlug
                    limit: $limit
                    skip: $skip
                    filters: $filters
                  ) {
                    hasMore
                    total
                    questions {
                      acRate
                      difficulty
                      freqBar
                      frontendQuestionId
                      isFavor
                      paidOnly
                      solutionNum
                      status
                      title
                      titleCn
                      titleSlug
                      topicTags {
                        name
                        nameTranslated
                        id
                        slug
                      }
                      extra {
                        hasVideoSolution
                        topCompanyTags {
                          imgUrl
                          slug
                          numSubscribed
                        }
                      }
                    }
                  }
                }
            ''',
            "variables": {
                "categorySlug": "all-code-essentials",
                "skip": offset,
                "limit": limit,
                "filters": {}
            },
            "operationName": "problemsetQuestionList",
        }
        encoded_data = json.dumps(data).encode('utf-8')
        retry_count = 5
        r = None
        for trial in range(1, retry_count + 1):
            r = requests.post(
                'https://leetcode.cn/graphql/',
                data=encoded_data,
                headers={
                    'Content-Type': 'application/json',
                    'Cookie': self.cookie,
                    'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36',
                },
            )
            if r.status_code == 200:
                break
            if trial < retry_count:
                print(
                    'Status %d got when fetch problems, will retry %d second(s) later...' % (r.status_code, trial ** 2))
                time.sleep(trial ** 2)
        if r.status_code != 200:
            raise RuntimeError('Fail to fetch problems! status: %d, data: %s' % (r.status_code, r.json()))
        return r.json()

    def fetch_problem_content(self, titleSlug: str) -> dict:
        dataCn = {
            "query": '''
                query questionTranslations($titleSlug: String!) {
                  question(titleSlug: $titleSlug) {
                    translatedTitle
                    translatedContent
                  }
            }
            ''',
            "variables": {
                "titleSlug": titleSlug
            },
            "operationName": "questionTranslations"
        }

        dataEn = {
            "query": '''
                query questionContent($titleSlug: String!) {
                  question(titleSlug: $titleSlug) {
                    content
                    editorType
                    mysqlSchemas
                    dataSchemas
                  }
                }''',
            "variables": {
                "titleSlug": titleSlug
            },
            "operationName": "questionContent"
        }
        # cn
        data = self.request(dataCn)
        contentCN = jsonpath.jsonpath(data, '$.data.question.translatedContent')[0]

        # en
        data = self.request(dataEn)
        titleEn = jsonpath.jsonpath(data, '$.data.question.content')[0]
        return {"cn": contentCN, "en": titleEn}

    def request(self, data):
        encoded_data = json.dumps(data).encode('utf-8')
        retry_count = 5
        r = None
        for trial in range(1, retry_count + 1):
            r = requests.post(
                'https://leetcode.cn/graphql/',
                data=encoded_data,
                headers={
                    'Content-Type': 'application/json',
                    'Cookie': self.cookie,
                    'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36',
                },
            )
            if r.status_code == 200:
                break
            if trial < retry_count:
                print(
                    'Status %d got when fetch problems, will retry %d second(s) later...' % (r.status_code, trial ** 2))
                time.sleep(trial ** 2)
        if r.status_code != 200:
            raise RuntimeError('Fail to fetch problems! status: %d, data: %s' % (r.status_code, r.json()))
        return r.json()

    def run(self):
        problems = {}
        for problem in self.fetch_problem_list():
            problems[problem['frontendQuestionId']] = problem

        i = 0
        for k, v in problems.items():
            print(f">>>>> progress: {i}/{len(problems)}", end="")
            content = self.fetch_problem_content(v['titleSlug'])
            problems[k]["content"] = content
            i+=1
        with open('./problems.json', 'w') as f:
            json.dump(problems, f, ensure_ascii=False)


if __name__ == "__main__":
    crawler = LcCrawler('zWST30if4WN4zAvPdsRNTblcAnEhxrndRmNPuVkMeRurf6IJo9D4rmbBh1hQQh8u')
    crawler.run()
