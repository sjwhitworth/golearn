package trees

import base "github.com/sjwhitworth/golearn/base"

type DecisionTree struct {
	base.BaseEstimator
}

type Branch struct {
	LeftBranch  Branch
	RightBranch Branch
}
