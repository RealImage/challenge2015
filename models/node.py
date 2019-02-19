#!/usr/bin/env python3

class Node():
  def __init__(self, value, category):
    self.value = value
    self.edges = []
    self.searched = False
    self.visited = False
    self.parent = None
    self.category = category
    self.meta = {}

  def __str__(self):
    return self.value

  def add_edge(self, node, meta):
    if node == self:
      return

    if node not in self.edges:
      node.meta[self.value] = meta
      self.edges.append(node)