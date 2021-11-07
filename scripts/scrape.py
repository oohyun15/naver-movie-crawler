import math

import requests
from bs4 import BeautifulSoup
import csv


class NaveMovie:
    def __init__(self, code):
        self.code = code
        self.total_size = 0
        self.total_page = self.get_total_page()
        self.score_dict = {}

    def get_total_page(self):
        url = "https://movie.naver.com/movie/bi/mi/pointWriteFormList.nhn?code=%d&type=after&onlyActualPointYn=N&onlySpoilerPointYn=N&order=lowest" % (
            self.code)
        response = requests.get(url, headers={
            "accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
            "accept-language": "ko-KR,ko;q=0.9",
            "sec-fetch-dest": "document",
            "sec-fetch-mode": "navigate",
            "sec-fetch-site": "same-site",
            "sec-fetch-user": "?1",
            "upgrade-insecure-requests": "1",
            "referrer": "https://movie.naver.com/movie/bi/mi/pointWriteFormList.nhn?code=%d&type=after&onlyActualPointYn=N&onlySpoilerPointYn=N&order=lowest" % self.code,
            'cookie': "NNB=4EAGKT2MDDNV6"
        })
        if response.status_code == 200:
            soup = BeautifulSoup(response.text, 'html.parser')
            total_size = int(soup.find("strong", "total").find("em").text.strip().replace(",", ""))
            self.total_size = total_size
            size_per_page = len(soup.find("div", "score_result").find_all("li"))
            total_page = math.ceil(total_size / size_per_page)
            return total_page
        else:
            raise

    def parse(self, page):
        url = "https://movie.naver.com/movie/bi/mi/pointWriteFormList.nhn?code=%d&type=after&onlyActualPointYn=N&onlySpoilerPointYn=N&order=lowest&page=%d" % (
            self.code, page)
        response = requests.get(url, headers={
            "accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
            "accept-language": "ko-KR,ko;q=0.9",
            "sec-fetch-dest": "document",
            "sec-fetch-mode": "navigate",
            "sec-fetch-site": "same-site",
            "sec-fetch-user": "?1",
            "upgrade-insecure-requests": "1",
            "referrer": "https://movie.naver.com/movie/bi/mi/pointWriteFormList.nhn?code=%d&type=after&onlyActualPointYn=N&onlySpoilerPointYn=N&order=lowest" % self.code,
            'cookie': "NNB=4EAGKT2MDDNV6"
        })
        if response.status_code == 200:
            soup = BeautifulSoup(response.text, 'html.parser')
            scores = soup.find("div", "score_result").find_all("li")
            for index, score in enumerate(scores):
                score_value = int(score.find("em").text)
                if score_value not in self.score_dict:
                    self.score_dict[score_value] = 0
                self.score_dict[score_value] += 1
        else:
            raise

movie_list = [
    {"name": "타짜", "code": 57723},
    {"name": "겨울왕국", "code": 100931},
    {"name": "극한직업", "code": 167651},
    {"name": "기생충", "code": 161967},
    {"name": "소울", "code": 184517},
]
with open('./result3.csv', 'w', newline='') as csvfile:
    fieldnames = ['title', 'score', 'count']
    writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
    for movie in movie_list:
        naver_movie = NaveMovie(movie["code"])
        print(movie['name'])
        for i in range(naver_movie.total_page):
            naver_movie.parse(i + 1)
            # break if 10 in naver_movie.score_dict
            print(f'({i+1}/{naver_movie.total_page})\t{naver_movie.score_dict}')
        for i in range(10):
            score = i + 1
            len_score = naver_movie.score_dict[score]
            writer.writerow({'title': movie["name"], 'score': score, 'count': len_score})
        writer.writerow({'title': movie["name"], 'score': "total_size", 'count': naver_movie.total_size})
