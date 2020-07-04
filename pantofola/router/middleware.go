package router

// Middleware is a chain component that definse a path execution mode
type Middleware func(r *Request, next *MiddlewareChain)

// structs

// MiddlewareChain is a linked list that makes possibile parallel executions
// without holding a state for every single connection we have on the sever
// we just move forword in the list and call all the middlewares
type MiddlewareChain struct {
	current Middleware
	next    *MiddlewareChain
	handler RequestHandler
}

//*********************************************************************************************************************

// Next runs the next middleware or the hander of the chain
// use this function to move forword in the chain
// or never call this to simply stop the chain execution
func (mc *MiddlewareChain) Next(req *Request) {
	if mc.current == nil {
		mc.handler(req)
	} else {
		mc.current(req, mc.next)
	}
}

// MakeMiddlewareChain creats a runnable chain of middlewares
// passing a nil middleware list will produce a chain that is the same as directly running the handler
func MakeMiddlewareChain(middlewares []Middleware, handler RequestHandler) *MiddlewareChain {
	if middlewares == nil {
		return &MiddlewareChain{handler: handler}
	}

	var first, last *MiddlewareChain
	for i := 0; i < len(middlewares); i++ {

		mc := &MiddlewareChain{current: middlewares[i], next: nil, handler: handler}

		if i == 0 {
			first = mc // the first must be saved
		} else {
			last.next = mc // update last chain item to poin the new one created
		}

		last = mc
	}

	// append handler to last
	last.next = &MiddlewareChain{handler: handler}

	return first // return head of list
}
