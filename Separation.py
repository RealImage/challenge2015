__author__ = 'PrabaKarthi'

import requests, json, multiprocessing, time, sys

check_list = []

'''
def worker(url):
    r1 = requests.get("http://data.moviebuff.com/" + url)
    if r1.status_code == 200:
            data = json.loads(r1.text)
            check_list.append({"url": url, "data": data})


pool=multiprocessing.pool(processes=4)
'''

counter = 0


class status_control:
    status_code = 200


def Separation(actor_from, actor_to, degree, count, history, check_list):
    global counter
    if count > degree:
        return False, count, history, check_list

    check = [x for x in check_list if x["url"] in actor_from]

    if len(check) > 0:
        return False, count, history, check_list
    else:
        counter += 1
        r1 = requests.get("http://data.moviebuff.com/" + actor_from)
        if r1.status_code == 200:
            data = json.loads(r1.text)
        check = [x for x in check_list if x["url"] in actor_to]
        if len(check) > 0:
            data1 = check[0]["data"]
            r2 = status_control
            r2.status_code = 200
        else:
            r2 = requests.get("http://data.moviebuff.com/" + actor_to)
            if r2.status_code == 200:
                data1 = json.loads(r2.text)
                check_list.append({"url": actor_to, "data": data1})
            else:
                return False, count, history, check_list
        if r1.status_code != 200 or r2.status_code != 200:
            return False, count, history, check_list

    if data["type"] == "Person" and data1["type"] == "Person":

        movie = [x for x in data["movies"] if x["url"] in [x["url"] for x in data1["movies"]]]

        if len(movie) > 0:
            history.append(movie[0])
            return True, count, history, check_list

        if len(data["movies"]) > len(data1["movies"]):
            check_list.append({"url": actor_to, "data": data1})
            data = data1
            actor_from, actor_to = actor_to, actor_from

    if data["type"] == "Person":
        movie = [x for x in data["movies"] if x["url"] == actor_to]

        if len(movie) > 0:
            history.append(movie[0])
            return True, count, history, check_list

        for movie in data["movies"]:
            count += 1
            history.append(movie)
            if actor_to == movie["url"]:
                return True, count, history, check_list

            a, b, c, d = Separation(movie["url"], actor_to, degree, count, history, check_list)
            check_list = d
            if a:
                return a, b, c, check_list
            else:
                history.remove(movie)
                count -= 1
    elif data["type"] == "Movie":
        crew = [x for x in data["crew"] if x["url"] == actor_to]

        if len(crew) > 0:
            history.append(crew[0])

            return True, count, history, check_list

        cast = [x for x in data["cast"] if x["url"] == actor_to]

        if len(cast) > 0:
            history.append(cast[0])
            return True, count, history, check_list

        merge = data["crew"] + data["cast"]

        merge = {v['url']: v for v in merge}.values()

        for l in merge:
            count += 1
            history.append(l)
            a, b, c, d = Separation(l["url"], actor_to, degree, count, history, check_list)
            check_list = d
            if a:
                return a, b, c, check_list
            else:
                history.remove(l)
                count -= 1

    return False, count, history, check_list


if __name__ == "__main__":
    fr = raw_input()
    to = raw_input()
    result, degree, connections, movie_caches = Separation(fr, to, 3, 0, [], check_list)
    print("Degree of Separation - {0}".format(degree))
    if result:
        for connection in connections:
            print(connection)
    else:
        print("not found")