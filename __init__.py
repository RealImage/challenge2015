#!/usr/bin/env python3
import sys
from models import Graph

''' 
Example: 

  python __init__.py deepti-naval girish-kulkarni

'''

if __name__ == '__main__':

  from_person = sys.argv[1]
  to_person = sys.argv[2]

  if from_person == to_person:
    sys.exit('Connection cannot be established to the same person')

  graph = Graph()
  graph.build_graph(from_person, to_person)
  graph.start_bfs(from_person, to_person)
  graph.get_the_shortest_connection(from_person, to_person)

  graph.print_path(from_person)

