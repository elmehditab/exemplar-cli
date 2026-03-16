package review

type Engine struct{}

func (e Engine) Run(req ReviewRequest) (ReviewResult, error) {

	pipeline := Pipeline{}

	result, err := pipeline.Run(req)

	if err != nil {
		return ReviewResult{}, err
	}

	return result, nil
}
