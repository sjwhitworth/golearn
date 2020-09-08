/*

	This package implements decision trees.

	ID3DecisionTree:
		Builds a decision tree using the ID3 algorithm
			by picking the Attribute which maximises
			Information Gain at each node.

		Attributes must be CategoricalAttributes at
			present, so discretise beforehand (see
			filters)

	CART (Classification and Regression Trees):
		Builds a binary decision tree using the CART algorithm
		using a greedy approach to find the best split at each node.

		Can be used for regression and classficiation.
		Attributes have to be FloatAttributes even for classification.
		Hence, convert to Integer Labels before hand for Classficiation.

	RandomTree:
		Builds a decision tree using the ID3 algorithm
			by picking the Attribute amongst those
			randomly selected that maximises Information
			Gain

		Attributes must be CategoricalAttributes at
			present, so discretise beforehand (see
			filters)

	IsolationForest:
		Unsupervised learning model for outlier detection.

		Builds a tree by randomly picking an attribute and splitting value.

		Attributes must be FloatAttributes.
		All Class Attributes will be treated as Normal Feature Attributes,
			So remove any Class Attributes you don't want during training beforehand.
*/

package trees
