package vnet

var endpoints = make(map[string][]*Endpoint)

//AddEndpoint - registers an endpoint
func AddEndpoint(ep *Endpoint) {
	eps, found := endpoints[ep.Category]
	if !found {
		eps = make([]*Endpoint, 0, 100)
		endpoints[ep.Category] = eps
	}
	eps = append(eps, ep)
}
