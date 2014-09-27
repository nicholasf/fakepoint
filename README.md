Fake Round Trip - a http.Client that will let you set up fake round trips for http testing.

##Install

`go get github.com/nicholasf/go


WIP.

Returns a client containing a fake round trip for mock http testing.


//1. a factory for producing fake clients
//- should take HTTP Method
//- expected body
//- URL (as a literal or a regex)
//- expected response status code
//- the expected URL parameters with corresponding assertions available via regex (can I get this to work for post vars)
//- perhaps convenience functions for working with JSON posts?
