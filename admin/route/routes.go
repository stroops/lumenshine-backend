package route

import (
	"fmt"
	//mw "github.com/Soneso/lumenshine-backend/admin/middleware"
)

//Route is a route in the API
type Route struct {
	Prefix string
	Method string
	Path   string
	Name   string

	//HandlerFunc is a little bit hacky, because we need the route package also in the implementation
	//therefore we have to use the func as interface, because else, we would get circle references
	HandlerFunc interface{}

	//RequiredGroups are the groups the user must be in, in order to access the route
	//the groups are used as OR, means, the user must be in at least one group
	RequiredGroups []string
}

//Routes dictionary of name->route
type Routes map[string]*Route

//AppRoutes holds all routes of the app
var AppRoutes Routes

//Add adds a route to the list
func Add(r *Route) {
	if AppRoutes == nil {
		AppRoutes = make(map[string]*Route)
	}

	_, exists := AppRoutes[r.Name]
	if exists {
		panic(fmt.Errorf("route-name %s already exists", r.Name))
	}

	//check that path does not exist yet
	for _, route := range AppRoutes {
		if r.Prefix+r.Path == route.Prefix+route.Path {
			panic(fmt.Errorf("route %s already exists", r.Path))
		}
	}

	AppRoutes[r.Name] = r
}

//AddRoute adds a route to the list
func AddRoute(Method string, Path string, HandlerFunc interface{}, RequiredGroups []string, Name string, Prefix string) {
	Add(&Route{
		Prefix:         Prefix,
		Name:           Name,
		Method:         Method,
		Path:           Path,
		HandlerFunc:    HandlerFunc,
		RequiredGroups: RequiredGroups,
	})
}

//GetRoutesForPrefix returns all routes for a prefix
func GetRoutesForPrefix(p string) []*Route {
	var routes []*Route
	for _, r := range AppRoutes {
		if r.Prefix == p {
			routes = append(routes, r)
		}
	}
	return routes
}

//GetRouteForName returns the route for a name or nil
func GetRouteForName(n string) *Route {
	r, exists := AppRoutes[n]
	if exists {
		return r
	}

	return nil
}
