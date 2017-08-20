package contexts

import "context"

// ContextUserKey user key
type ContextUserKey struct{}

// ContextParamsKey params key
type ContextParamsKey struct{}

// ContextPathKey path key
type ContextPathKey struct{}

// ContextArgumentsKey path key
type ContextArgumentsKey struct{}

var (
	// UserKey static key for UserObject in context
	UserKey ContextUserKey
	// ParamsKey static key for Query parameters in request
	ParamsKey ContextParamsKey
	// PathKey path string key for context
	PathKey ContextPathKey
	// ArgsKey key for arguments in URL
	ArgsKey ContextArgumentsKey
)

// SetPath set a string value to be used as a path string to a context
func SetPath(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, PathKey, value)
}

// GetPath gets a path argument if any
func GetPath(ctx context.Context) (path string) {
	if ctx != nil {
		if val, ok := ctx.Value(PathKey).(string); ok {
			path = val
		}
	}
	return
}

// SetParams sets a map[string]string to the context
// used mainly for url arguments like key=value
func SetParams(ctx context.Context, value map[string]string) context.Context {
	return context.WithValue(ctx, ParamsKey, value)
}

// GetParams gets a param as map[string]string if exists
// used mainly for url arguments like key=value
func GetParams(ctx context.Context) (params map[string]string) {
	if ctx != nil {
		if val, ok := ctx.Value(ParamsKey).(map[string]string); ok {
			params = val
		}
	}
	return
}

// SetUser sets a interface{} value to the user key
func SetUser(ctx context.Context, value interface{}) context.Context {
	return context.WithValue(ctx, UserKey, value)
}

// GetUser gets the user from the user key
func GetUser(ctx context.Context) (user interface{}) {
	if ctx != nil {
		user = ctx.Value(UserKey)
	}
	return
}

// SetArgs set a map to a a args key
// used mainly for url defined parameters like name or id for RESTful APIs
func SetArgs(ctx context.Context, value map[string]string) context.Context {
	return context.WithValue(ctx, ArgsKey, value)
}

// GetArgs get the arguments and return as a map[string]string
// used mainly for url defined parameters like name or id for RESTful APIs
func GetArgs(ctx context.Context) (value map[string]string) {
	if ctx != nil {
		if val, ok := ctx.Value(ArgsKey).(map[string]string); ok {
			value = val
		}
	}
	return
}
