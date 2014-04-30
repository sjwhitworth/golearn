package base

// Specifies the base interfaces
type Estimator interface {
	Fit()
	Summarise()
}

type Predictor interface {
	Predict()
}

type Model interface {
	Score()
}
