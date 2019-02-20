#!/usr/bin/env python3
from collections import defaultdict, OrderedDict, deque
import json
import requests
from models.node import Node
from pprint import pprint

class Graph():
  def __init__(self):
    # A normal graph structure
    self.nodes = []
    self.graph = defaultdict(lambda: None)
    self.end = None
    self.start = None

  def set_start(self, value):
    self.start = self.graph[value]
    return self.graph[value]

  def set_end(self, value):
    self.end = self.graph[value]
    return self.graph[value]

  def add_node(self, node):
    # if node not in self.nodes:
    self.nodes.append(node)
    self.graph[node.value] = node

  def get_node(self, value):
    return self.graph[value]

  def build_graph(self, from_person, to_person):
    print("Building Graph from {0} to {1}\n".format(from_person, to_person))
    persons = deque([{'url': from_person}])
    found = False
    while persons:
      person = persons.popleft()
      person_node = self.get_node(person['url'])
      if person_node is not None:
        if person_node.visited:
          # Reduce the number of operations
          continue

      # Loading the person's data information
      try:
        with open("data/persons/{0}.json".format(person['url'])) as data_file:
          data = json.load(data_file)
      except:
        request_object = requests.get("http://data.moviebuff.com/{0}".format(person['url']))
        if request_object.status_code == 200:
          data = request_object.json(object_pairs_hook=OrderedDict)
          with open("data/persons/{0}.json".format(person['url']), 'w') as outfile:
            json.dump(data, outfile)
        else:
          continue
        
      # Gathering all his movies
      person = data
      movies = deque(data['movies'])

      # Creating person node and adding to graph if needed
      if person_node is None:
        person_node = Node(person['url'], "Person")
        self.add_node(person_node)

      # person_node.visited = True

      while movies:
        movie = movies.popleft()

        
        # Loading the movies's data information
        try:
          with open("data/movies/{0}.json".format(movie['url'])) as data_file:
            data = json.load(data_file)
        except:
          request_object = requests.get("http://data.moviebuff.com/{0}".format(movie['url']))
          if request_object.status_code == 200:
            data = request_object.json(object_pairs_hook=OrderedDict)
            with open("data/movies/{0}.json".format(movie['url']), 'w') as outfile:
              json.dump(data, outfile)
          else:
            continue

        actors = data['crew'] + data['cast']

        movie_node = self.get_node(movie['url'])

        # Create the movie node if needed
        if movie_node is None:
          movie_node = Node(movie['url'], "Movie")
          self.add_node(movie_node)
        if movie_node.visited:
          continue

        # Attach the person node with the movie node and vice versa
        person_node.add_edge(movie_node, data)
        movie_node.add_edge(person_node, None)

        for actor in actors:
          actor_node = self.get_node(actor['url'])
          # Create actor node if necessary
          if actor_node is None:
            actor_node = Node(actor['url'], "Person")
            self.add_node(actor_node)

          movie_node.add_edge(actor_node, actor)
          actor_node.add_edge(movie_node, data)

          if actor['url'] == to_person:
            found = True
            break
          persons.append(actor)
        movie_node.visited = True
        if found:
          break
      # print("Still searching..")
      if found:
        break

  def start_bfs(self, from_person, to_person):
    start_node = self.set_start(from_person)
    end_node = self.set_end(to_person)
    queue = deque([start_node])

    while queue:
      current_node = queue.popleft()
      current_node.searched = True

      for edge in current_node.edges:
        # Search and explore
        if edge.searched == False:
          queue.append(edge)
          edge.searched = True
          edge.parent = current_node

  def get_the_shortest_connection(self, from_person, to_person):
    print("Finding the shortest path to connect both.")
    start_node = self.set_start(from_person)
    current_node = self.set_end(to_person)
    self.path = []

    while current_node != start_node:
      # Traverse till the start node
      self.path.insert(0, current_node)
      current_node = current_node.parent
    self.path.insert(0, start_node)

  def get_person_movie_info(self, person_name, movie_name):
    request_object = requests.get("http://data.moviebuff.com/{0}".format(person_name))
    if request_object.status_code == 200:
      person_data = request_object.json()

    request_object = requests.get("http://data.moviebuff.com/{0}".format(movie_name))
    if request_object.status_code == 200:
      movie_data = request_object.json()

    actors = movie_data["cast"] + movie_data["crew"]

    for actor in actors:
      if actor["url"] == person_name:
        role = actor["role"]
        break
      role = ""
    
    return {"name": person_data["name"], "role": role}

  def load_print_path(self, from_person):
    # This is the naive approach to pull up metas to print the necessary info. (Can be improved too.)
    current_person = from_person
    current_movie = None
    responses = []
    response = {}
    found_p1, found_p2 = False, False
    for path in self.path:
      if path.category == "Person":
        if found_p1:
          response["person_2"] = path.value
          found_p2 = True
          response["current_meta"] = path.meta
          try:
            0/0
            response["person_1_meta"] = response["movie_meta"][response["person_1"]]
          except:
            response["person_1_meta"] = self.get_person_movie_info(response["person_1"], response["movie"])
          response["person_2_meta"] = response["movie_meta"][response["person_2"]]
          response["movie"] = response["current_meta"][response["movie"]]["name"]
          response.pop('current_meta', None)
          response.pop('movie_meta', None)
          responses.append(response)
        else:
          if "person_1" not in response:
            response["person_1"] = path.value
          found_p1 = True
      else:
        if found_p1 and found_p2:
          found_p2 = False
          response = {"person_1": response["person_2"], "person_1_meta": response["person_2_meta"]}
        response["movie"] = path.value
        response["movie_meta"] = path.meta

    return responses

  def print_path(self, from_person):
    responses = self.load_print_path(from_person)
    for response in responses:
      print("Movie: {0}".format(response["movie"]))
      person_1_meta = response["person_1_meta"]
      print("{0}: {1}".format(person_1_meta["role"], person_1_meta["name"]))
      person_2_meta = response["person_2_meta"]
      print("{0}: {1}".format(person_2_meta["role"], person_2_meta["name"]))
      print("\n")

