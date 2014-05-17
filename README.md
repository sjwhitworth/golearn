GoLearn
=======

<img src="http://talks.golang.org/2013/advconc/gopherhat.jpg" width=125><br>
[![GoDoc](https://godoc.org/github.com/sjwhitworth/golearn?status.png)](https://godoc.org/github.com/sjwhitworth/golearn)<br>

GoLearn is a 'batteries included' machine learning library for Go. **Simplicity**, paired with customisability, is the goal.
We are in active development, and would love comments from users out in the wild. Drop us a line on Twitter.

twitter: [@golearn_ml](http://www.twitter.com/golearn_ml)

Install
=======

```
go get github.com/sjwhitworth/golearn
cd src/github.com/sjwhitworth/golearn
go get ./...
```

Getting Started
=======

Data are loaded in as Instances. You can then perform matrix like operations on them, and pass them to estimators.
We implement the scikit-learn interface of Fit/Predict, so you can easily swap out estimators for trial and error.

```
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
confusionMat := evaluation.GetConfusionMatrix(testData, predictions)
fmt.Println(evaluation.GetSummary(confusionMat))
```

```
Iris-virginica	28	2	  56	0.9333	0.9333  0.9333
Iris-setosa	    29	0	  59	1.0000  1.0000	1.0000
Iris-versicolor	27	2	  57  0.9310	0.9310  0.9310
Overall accuracy: 0.9545
```

Examples
========

GoLearn comes with practical examples. Dive in and see what is going on.

```
cd examples/
go run knnclassifier_iris.go
go run instances.go
```

Join the team
=============

Please send me a mail at stephen dot whitworth at hailocab dot com.
