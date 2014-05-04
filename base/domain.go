// Package base provides base interfaces for GoLearn objects to implement.
// It also provides a raw base for those objects.
package base

import (
	mat64 "github.com/gonum/matrix/mat64"
)

// An object that can ingest some data and train on it.
type Estimator interface {
	Fit()
	Summarise()
}

// An object that provides predictions.
type Predictor interface {
	Predict()
}

// An supervised learning object, that is possible of scoring accuracy against a test set.
type Model interface {
	Score()
}

// @todo: Implement BaseEstimator setters and getters.
type BaseEstimator struct {
	Estimator
	Data *mat64.Dense
}
