// Package memo memoizing a functions output for cache utilization after first go
package memo

// Memo caches the results of a function
type Memo struct {
	f     Func
	cache map[string]result
}

// Func is the type of a function to memoize.
type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

// New Memo
func New(f Func) *Memo {
	return &Memo{
		f:     f,
		cache: make(map[string]result),
	}
}

// Get ...
// NOTE: not concurrency-safe!
func (memo *Memo) Get(key string) (interface{}, error) {
	res, ok := memo.cache[key] // if does not exist, returns error (aka "ok" is false therefore satisfying the init of function in if-block)
	if !ok {
		res.value, res.err = memo.f(key)
		memo.cache[key] = res
	}
	return res.value, res.err
}
