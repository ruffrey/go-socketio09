package socketio09

import "reflect"

// HandlerCaller calls the function `Func` with optional arguments
type HandlerCaller struct {
	Func        reflect.Value
	Args        reflect.Type
	ArgsPresent bool
	Out         bool
}

/*
NewHandlerCaller parses function passed by using reflection, and stores its representation
for further call on message or ack. The callback handler is validated for conformity to the
expected handler format.
*/
func NewHandlerCaller(fn interface{}) (*HandlerCaller, error) {
	fnValOf := reflect.ValueOf(fn)
	if fnValOf.Kind() != reflect.Func {
		return nil, ErrorCallerShouldBeTypeFunc
	}

	fType := fnValOf.Type()
	if fType.NumOut() > 1 {
		return nil, ErrorCallerFunctionReturnsTooMuch
	}

	currentCaller := &HandlerCaller{
		Func: fnValOf,
		Out:  fType.NumOut() == 1,
	}
	if fType.NumIn() == 1 {
		currentCaller.Args = nil
		currentCaller.ArgsPresent = false
		return currentCaller, nil
	}
	if fType.NumIn() == 2 {
		currentCaller.Args = fType.In(1)
		currentCaller.ArgsPresent = true
		return currentCaller, nil
	}

	return nil, ErrorCallerShouldHaveTwoArgs
}

/*
getArgs returns params for a handler function
*/
func (c *HandlerCaller) getArgs() interface{} {
	return reflect.New(c.Args).Interface()
}

/*
callFunc calls a handler with arguments
*/
func (c *HandlerCaller) callFunc(h *SocketIOConnection, args interface{}) []reflect.Value {
	//nil is untyped, so use the default empty value of correct type
	if args == nil {
		args = c.getArgs()
	}

	a := []reflect.Value{reflect.ValueOf(h), reflect.ValueOf(args).Elem()}
	if !c.ArgsPresent {
		a = a[0:1]
	}

	return c.Func.Call(a)
}
