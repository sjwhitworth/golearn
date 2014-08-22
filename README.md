GoLearn
=======

<img src="http://talks.golang.org/2013/advconc/gopherhat.jpg" width=125><br>
[![GoDoc](https://godoc.org/github.com/sjwhitworth/golearn?status.png)](https://godoc.org/github.com/sjwhitworth/golearn)
[![Build Status](https://travis-ci.org/sjwhitworth/golearn.png?branch=master)](https://travis-ci.org/sjwhitworth/golearn)<br>

[![Support via Gittip](https://rawgithub.com/twolfson/gittip-badge/0.2.0/dist/gittip.png)](https://www.gittip.com/sjwhitworth/)

GoLearn is a 'batteries included' machine learning library for Go. **Simplicity**, paired with customisability, is the goal.
We are in active development, and would love comments from users out in the wild. Drop us a line on Twitter.

twitter: [@golearn_ml](http://www.twitter.com/golearn_ml)

Install
=======

See [here](https://github.com/sjwhitworth/golearn/wiki/Installation) for installation instructions.

Getting Started
=======

Data are loaded in as Instances. You can then perform matrix like operations on them, and pass them to estimators.
GoLearn implements the scikit-learn interface of Fit/Predict, so you can easily swap out estimators for trial and error.
GoLearn also includes helper functions for data, like cross validation, and train and test splitting.

```go
// Load in a dataset, with headers. Header attributes will be stored.
// Think of instances as a Data Frame structure in R or Pandas.
// You can also create instances from scratch.
data, err := base.ParseCSVToInstances("datasets/iris_headers.csv", true)

// Print a pleasant summary of your data.
fmt.Println(data)

// Split your dataframe into a training set, and a test set, with an 80/20 proportion.
trainTest := base.InstancesTrainTestSplit(rawData, 0.8)
trainData := trainTest[0]
testData := trainTest[1]

// Instantiate a new KNN classifier. Euclidean distance, with 2 neighbours.
cls := knn.NewKnnClassifier("euclidean", 2)

// Fit it on your training data.
cls.Fit(trainData)

// Get your predictions against test instances.
predictions := cls.Predict(testData)

// Print a confusion matrix with precision and recall metrics.
confusionMat, _ := evaluation.GetConfusionMatrix(testData, predictions)
fmt.Println(evaluation.GetSummary(confusionMat))
```

```
Iris-virginica	28	2	  56	0.9333	0.9333  0.9333
Iris-setosa	    29	0	  59	1.0000  1.0000	1.0000
Iris-versicolor	27	2	  57	0.9310	0.9310  0.9310
Overall accuracy: 0.9545
```

Examples
========

GoLearn comes with practical examples. Dive in and see what is going on.

```bash
cd $GOPATH/src/github.com/sjwhitworth/golearn/examples/knnclassifier
go run knnclassifier_iris.go
```
```bash
cd $GOPATH/src/github.com/sjwhitworth/golearn/examples/instances
go run instances.go
```
```bash
cd $GOPATH/src/github.com/sjwhitworth/golearn/examples/trees
go run trees.go
```

Join the team
=============

Please send me a mail at stephen dot whitworth at hailocab dot com.
