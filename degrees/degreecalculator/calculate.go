package degreecalculator

import (
	"degrees/moviebuffclient"
	"fmt"
	"sync"
	"sync/atomic"
)

type void struct{}

type node struct {
	parent     *node
	parentRole string
	person
}

type person struct {
	name  string
	url   string
	role  string
	movie string
}

type Result struct {
	Level int
	Node  *node
	Err   error
}

var pregL, mregL sync.Mutex

func Calculate(p1, p2 string) (int, *node, error) {
	ch := make(chan Result)
	// maintains the list of visited urls in the search process
	registry := make(map[string]void)
	registry[p1] = void{}
	movieRegistry := make(map[string]void)
	parentNode := &node{person: person{url: p1}}
	list := []*node{parentNode}
	go traverse(ch, 1, p2, list, registry, movieRegistry)
	res := <-ch
	return res.Level, res.Node, res.Err
}

func traverse(ch chan Result, level int, destinationUrl string, list []*node, registry map[string]void, movieRegistry map[string]void) {

	if len(list) == 0 {
		ch <- Result{-1, nil, nil}
		return
	}

	terminate := new(atomic.Bool)
	terminate.Store(false)

	var nextLevelList []*node

	wg := sync.WaitGroup{}
	maxGoroutines := make(chan struct{}, 150)
	defer close(maxGoroutines)

	// fetch direct assocaited persons of all the person in the current level
	for _, p := range list {
		if terminate.Load() {
			return
		}
		wg.Add(1)
		maxGoroutines <- struct{}{}

		go func(p *node) {
			defer wg.Done()
			defer func() {
				<-maxGoroutines
			}()

			// fetch person info
			personInfo, err := moviebuffclient.GetPersonInfo(p.url)
			if err != nil {
				ch <- Result{-2, nil, err}
				return
			}

			if p.parent == nil {
				// update person info in the node
				p.name = personInfo.Name
			}

			for _, m := range personInfo.Movies {
				// check if the movie is already visited
				mregL.Lock()
				if _, ok := movieRegistry[m.Url]; ok {
					mregL.Unlock()
					continue
				}
				// add the movie to the registry
				movieRegistry[m.Url] = void{}
				mregL.Unlock()

				// fetch movie info
				movieInfo, err := moviebuffclient.GetMovieInfo(m.Url)
				if err != nil {
					ch <- Result{-2, nil, err}
					return
				}

				parentRole := m.Role

				for _, c := range movieInfo.Cast {
					// generate a new node
					newNode := &node{
						p,
						parentRole,
						person{
							name:  c.Name,
							url:   c.Url,
							role:  c.Role,
							movie: movieInfo.Name,
						},
					}

					// check if the destination url is reached
					if c.Url == destinationUrl {
						if terminate.Load() {
							return
						}
						ch <- Result{level, newNode, nil} // complete the function
						terminate.Store(true)
						return
					}
					// check if the person is already visited
					pregL.Lock()
					if _, ok := registry[c.Url]; !ok {
						// add the person to the registry
						registry[c.Url] = void{}
						// add the person to the next level
						nextLevelList = append(nextLevelList, newNode)
						pregL.Unlock()
						continue
					}
					pregL.Unlock()

				}
			}
		}(p)
	}
	wg.Wait()

	traverse(ch, level+1, destinationUrl, nextLevelList, registry, movieRegistry)
}

func PrintRespose(level int, n *node, err error) {
	if checkForNonTrivialBehavior(level, err) {
		return
	}
	var entries []string
	fmt.Println("Degrees of seperation: ", level)
	count := level
	cur := n
	for cur.parent != nil {
		entries = append(entries, fmt.Sprintf(`%d. Movie: %s
%s: %s
%s: %s%s`, count, cur.movie, cur.parentRole, cur.parent.name, cur.role, cur.name, "\n"))
		cur = cur.parent
		count--
	}
	for i := len(entries) - 1; i >= 0; i-- {
		fmt.Println(entries[i])
	}
}

func PrintResponseInReverse(level int, n *node, err error) {
	if checkForNonTrivialBehavior(level, err) {
		return
	}
	fmt.Println("Degrees of seperation: ", level)
	count := 1
	cur := n
	for cur.parent != nil {
		fmt.Printf(`%d. Movie: %s
%s: %s
%s: %s%s`, count, cur.movie, cur.role, cur.name, cur.parentRole, cur.parent.name, "\n")
		cur = cur.parent
		count++
	}
}

func checkForNonTrivialBehavior(level int, err error) bool {
	if err != nil {
		fmt.Println("Error in calculating degree of seperation:", err)
		return true
	}
	if level == -1 {
		fmt.Println("No degree of seperation found")
		return true
	}
	return false
}
