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