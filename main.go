package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
)

type movieBuffType string

const (
	movieBuff = "http://data.moviebuff.com/"

	unknown             movieBuffType = "Unknown"
	person                            = "Person"
	movie                             = "Movie"
	defaultDiskCacheDir               = "./httpcache"
)

type entity struct {
	URL    string    `json:"url"`
	Name   string    `json:"name"`
	Type   string    `json:"type"`
	Movies []*detail `json:"movies"`
	Cast   []*detail `json:"cast"`
}

type detail struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Role string `json:"role"`
}

type node struct {
	Prev        *node
	Data        string
	Connections []*node
	From        int
}

var (
	lock    = &sync.RWMutex{}
	visited = map[string]struct{}{}
)

var head1, head2 *node

func main() {
	if len(os.Args) < 3 {
		panic("Not enough input")
	}

	node1 := os.Args[1]
	node2 := os.Args[2]

	head1 = &node{
		Data: node1,
		From: 1,
	}
	head2 = &node{
		Data: node2,
		From: 2,
	}

	degrees := make(chan int)
	jobs := make(chan *node)
	results := make(chan *node)

	for w := 0; w < 2; w++ {
		go worker(jobs, results, degrees)
	}

	go successors(head1, jobs)
	go successors(head2, jobs)

	go func() {
		for {
			select {
			case e := <-results:
				visits := []string{}
				if e.From == 1 {
					go intersects(head2, e, visits, 0, degrees)
					go successors(e, jobs)
				} else {
					go intersects(head1, e, visits, 0, degrees)
					go successors(e, jobs)
				}
			case <-degrees:
				return
			}
		}
	}()

	fmt.Println(<-degrees)
}

func intersects(head, currentNode *node, visits []string, degrees int, done chan<- int) {
	degrees++
	if head.Data == currentNode.Data {
		done <- degrees
		return
	}

	for _, c := range head.Connections {
		hasc := false
		for _, v := range visits {
			if v == c.Data {
				hasc = true
				break
			}
		}

		if !hasc {
			visits = append(visits, c.Data)

			go intersects(c, currentNode, visits, degrees, done)
		}
	}
}

func walk(n *node) {
	fmt.Println(n.Data)
	for _, c := range n.Connections {
		walk(c)
	}
}

func successors(n *node, jobs chan<- *node) {
	lock.RLock()
	if _, ok := visited[n.Data]; ok {
		lock.RUnlock()
		return
	}
	lock.RUnlock()

	jobs <- n
}

func worker(jobs <-chan *node, results chan<- *node, done <-chan int) {
	for {
		select {
		case j := <-jobs:
			e := fetch(j.Data)
			lock.Lock()
			visited[j.Data] = struct{}{}
			lock.Unlock()

			if e != nil {
				var data []*detail
				if len(e.Movies) > 0 {
					data = e.Movies
				} else {
					data = e.Cast
				}

				jc := []*node{}

				for _, d := range data {
					lock.RLock()
					if _, ok := visited[d.URL]; ok {
						lock.RUnlock()
						continue
					}
					lock.RUnlock()

					cn := &node{
						Prev: j,
						From: j.From,
						Data: d.URL,
						Connections: []*node{
							j,
						},
					}

					jc = append(jc, cn)
				}

				j.Connections = jc
				results <- j

				for _, cn := range jc {
					results <- cn
				}
			}

		case <-done:
			return
		}
	}
}

func fetch(node string) *entity {
	resp, err := cachingHTTPClient().Get(movieBuff + node)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	var e entity

	err = json.NewDecoder(resp.Body).Decode(&e)
	if err != nil {
		return nil
	}

	return &e
}

func cachingHTTPClient() *http.Client {
	return httpcache.NewTransport(diskcache.New(defaultDiskCacheDir)).Client()
}
