package vnet

var endpoints = make(map[string][]*Endpoint)

//AddEndpoint - registers an REST endpoint
func AddEndpoint(ep *Endpoint) {
	eps, found := endpoints[ep.Category]
	if !found {
		eps = make([]*Endpoint, 0, 100)
		endpoints[ep.Category] = eps
	}
	eps = append(eps, ep)
}

//AddEndpoints - registers multiple REST endpoints
func AddEndpoints(eps ...*Endpoint) {
	for _, ep := range eps {
		AddEndpoint(ep)
	}
}

func Serve(port int) (err error) {
	return err
}
