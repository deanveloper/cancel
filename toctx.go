package cancel

import "context"

//
// Context wraps `canc` in an implementation of a `context.Context`.
//
// Usage:
//
//		db.QueryContext(cancel.Context(canc), query, args)
//
func Context(canc Canceller) context.Context {
	return ctxWrap{Canceller: canc}
}

type ctxWrap struct {
	Canceller
}

// Value always returns nil, and is only used to
// implement the context.Context interface.
func (ctxWrap) Value(key interface{}) interface{} {
	return nil
}
