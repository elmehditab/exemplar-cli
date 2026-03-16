package app

import "github.com/mehditabet/exemplar-cli/internal/core/review"

func RunReview(req review.ReviewRequest) (review.ReviewResult, error) {
	engine := review.Engine{}

	result, err := engine.Run(req)
	if err != nil {
		return review.ReviewResult{}, err
	}

	return result, nil
}
