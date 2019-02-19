#!/usr/bin/env python3
from models.graph import Graph

if __name__ == '__main__':

  from_person = 'ellen-barkin'
  to_person = 'robert-de-niro'

  graph = Graph()
  graph.build_graph(from_person, to_person)
  graph.start_bfs(from_person, to_person)
  graph.get_the_shortest_connection(from_person, to_person)

  graph.print_path(from_person)
  # graph.reverse_path()

