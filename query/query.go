package query

// Query used for querying resource database
type Query map[string]interface{}

// New query base constructor
func New() Query {
	return Query{}
}

// Add adds a value to a key
func (q Query) Add(key string, value interface{}) Query {
	q[key] = value
	return q
}

const (
	orderByKey = "__order_by__"
	paramsKey  = "__params__"
)

// OrderBy adds order by
func (q Query) OrderBy(key string, ascending bool) Query {
	q[orderByKey] = Ordering{key, ascending}
	return q
}

// GetFields return fields
func (q Query) GetFields() Query {
	return FilterKey(q, orderByKey, paramsKey)
}

// GetOrderBy returns key order by, ascending, and ok
// last boolean will return false if there are no keys there
func (q Query) GetOrderBy() (key string, ascending bool, ok bool) {
	var buff Ordering
	buff, ok = q[orderByKey].(Ordering)
	if ok {
		ascending = buff.Ascending
		key = buff.Key
	}
	return
}

// GetParams retrieves parameters, and ok
// ok means it exists otherwise returns false
// if not existing Parameters will be nil
func (q Query) GetParams() (params Parameters) {
	if q != nil {
		params, _ = q[paramsKey].(Parameters)
	}
	return
}

// AddParam add a parameter
func (q Query) AddParam(key string, value interface{}) Query {
	params := q.GetParams()
	if params == nil {
		params = make(Parameters)
	}
	params.Add(key, value)
	q.Add(paramsKey, params)
	return q
}

// Parameters parameters type for dataquery
type Parameters map[string]interface{}

// Has returns true if key exists
func (p Parameters) Has(key string) (ok bool) {
	if p != nil {
		_, ok = p[key]
	}
	return
}

// GetBool returns a parameter key as a boolean
func (p Parameters) GetBool(key string) (val bool) {
	if p != nil {
		val, _ = p[key].(bool)
	}
	return
}

// GetString returns a parameter key as a boolean
func (p Parameters) GetString(key string) (val string) {
	if p != nil {
		val, _ = p[key].(string)
	}
	return
}

// Add adds a parameter
func (p Parameters) Add(key string, val interface{}) Parameters {
	if p == nil {
		p = make(Parameters)
	}
	p[key] = val
	return p
}

// Ordering sorting class
type Ordering struct {
	Key       string
	Ascending bool
}

// FilterKey will filter the keys of the map and return a new map instance
func FilterKey(data map[string]interface{}, keys ...string) map[string]interface{} {
	new := make(map[string]interface{}, len(data))
	for k, v := range data {
		new[k] = v
	}
	for _, f := range keys {
		delete(new, f)
	}
	return new
}
