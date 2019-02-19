#!/usr/bin/env python3
from models.graph import Graph

if __name__ == '__main__':

  # print("Enter the first person name:")
  # from_person = raw_input().strip()
  # print("Enter the second person name:")
  # to_person = raw_input().strip()
  from_person = "deepti-naval"
  to_person = "girish-kulkarni"

  graph = Graph()
  graph.build_graph(from_person, to_person)
  graph.start_bfs(from_person, to_person)
  graph.get_the_shortest_connection(from_person, to_person)

  graph.print_path(from_person)

