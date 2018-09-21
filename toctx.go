package cancel

// CtxWrap represents an implementation of
// a context.Context.
//
// Usage:
//
//		db.QueryContext(cancel.CtxWrap{canc}, query, args)
//
type CtxWrap struct {
	Canceller
}

// Value always returns nil, and is only used to
// implement the context.Context interface.
func (CtxWrap) Value(key interface{}) interface{} {
	return nil
}
