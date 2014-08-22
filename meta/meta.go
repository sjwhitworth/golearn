/*

	Meta contains base.Classifier implementations which
		combine the outputs of others defined elsewhere.

	Bagging:
		Bootstraps samples of the original training set
		with a number of selected attributes, and uses
		that to train an ensemble of models. Predictions
		are generated via majority voting.
*/

package meta
