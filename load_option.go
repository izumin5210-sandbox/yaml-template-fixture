package tfx

type Params map[string]interface{}

type LoadOption func(Params)

func createParams(opts []LoadOption) Params {
	params := Params{}
	for _, f := range opts {
		f(params)
	}
	return params
}

func WithParams(params Params) LoadOption {
	return func(prevParams Params) {
		for k, v := range params {
			prevParams[k] = v
		}
	}
}

func WithLoop(name string, loopCount int) LoadOption {
	return func(params Params) {
		res := make([]int, loopCount)
		for i := 0; i < loopCount; i++ {
			res[i] = i + 1
		}
		params[name] = res
	}
}
