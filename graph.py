#!/usr/bin/env python3

class Graph():
  def __init__(self):
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
    if node not in self.nodes:
      self.nodes.append(node)
      self.graph[node.value] = node

  def get_node(self, value):
    return self.graph[value]

  def build_graph(self, from_person, to_person):
    persons = deque([{'url': from_person}])
    found = False
    while persons:
      person = persons.popleft()
      person_node = self.get_node(person['url'])
      if person_node is not None:
        if person_node.visited:
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
      movies = deque(data['movies'])

      # Creating person node and adding to graph if needed
      if person_node is None:
        person_node = Node(person['url'], "Person")
        self.add_node(person_node)

      person_node.visited = True

      while movies:
        movie = movies.popleft()
        movie_node = self.get_node(movie['url'])

        # Create the movie node if needed
        if movie_node is None:
          movie_node = Node(movie['url'], "Movie")
          self.add_node(movie_node)
        if movie_node.visited:
          continue

        # Attach the person node with the movie node and vice versa
        person_node.add_edge(movie_node, movie)
        movie_node.add_edge(person_node, person)
        
        # Loading the movies's data information
        # with open("data/{0}.json".format(movie['url'])) as data_file:
          # data = json.load(data_file)
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
        # actors = deque([i['url'] for i in actors])

        for actor in actors:
          actor_node = self.get_node(actor['url'])
          # Create actor node if necessary
          if actor_node is None:
            actor_node = Node(actor['url'], "Person")
            self.add_node(actor_node)

          movie_node.add_edge(actor_node, actor)
          actor_node.add_edge(movie_node, movie)

          if actor['url'] == to_person:
            found = True
            print("Found")
            break
          persons.append(actor)
        movie_node.visited = True

        if found:
          print("Found 118")
          break
      
      print("Done for {0}".format(person['url']))
      if found:
        print("Found 121")
        break