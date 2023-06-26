package main

type PairSet struct {
	m map[string]string
}

func New() *PairSet {
	r := &PairSet{}
	r.m = make(map[string]string)
	return r
}

func (st *PairSet) Insert(key string, value string) {
	st.m[key] = value
}

func (st *PairSet) IsPresent(key string) bool {
	_, ok := st.m[key]
	return ok
}

func (st *PairSet) Get(key string) string {
	return st.m[key]
}

func (st *PairSet) Delete(key string) {
	delete(st.m, key)
}

func (st *PairSet) GetAllKeys() []string {
	var result []string
	for k := range st.m {
		result = append(result, k)
	}
	return result
}
