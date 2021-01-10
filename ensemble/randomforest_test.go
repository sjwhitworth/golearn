package ensemble

import (
	"testing"

	"io/ioutil"
	"os"

	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/filters"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRandomForest(t *testing.T) {
	Convey("Given a valid CSV file", t, func() {
		inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
		So(err, ShouldBeNil)

		Convey("When Chi-Merge filtering the data", func() {
			filt := filters.NewChiMergeFilter(inst, 0.90)
			for _, a := range base.NonClassFloatAttributes(inst) {
				filt.AddAttribute(a)
			}
			filt.Train()
			instf := base.NewLazilyFilteredInstances(inst, filt)

			Convey("Splitting the data into test and training sets", func() {
				trainData, testData := base.InstancesTrainTestSplit(instf, 0.60)

				Convey("Fitting and predicting with a Random Forest", func() {
					rf := NewRandomForest(10, 3)
					err = rf.Fit(trainData)
					So(err, ShouldBeNil)

					predictions, err := rf.Predict(testData)
					So(err, ShouldBeNil)

					confusionMat, err := evaluation.GetConfusionMatrix(testData, predictions)
					So(err, ShouldBeNil)

					Convey("Predictions should be somewhat accurate", func() {
						So(evaluation.GetAccuracy(confusionMat), ShouldBeGreaterThan, 0.35)
					})
				})
			})
		})

		Convey("Fitting with a Random Forest with too many features compared to the data", func() {
			rf := NewRandomForest(10, len(base.NonClassAttributes(inst))+1)
			err = rf.Fit(inst)

			Convey("Should return an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestRandomForestSerialization(t *testing.T) {
	Convey("Given a valid CSV file", t, func() {
		inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
		So(err, ShouldBeNil)

		Convey("When Chi-Merge filtering the data", func() {
			filt := filters.NewChiMergeFilter(inst, 0.90)
			for _, a := range base.NonClassFloatAttributes(inst) {
				filt.AddAttribute(a)
			}
			filt.Train()
			instf := base.NewLazilyFilteredInstances(inst, filt)

			Convey("Splitting the data into test and training sets", func() {
				trainData, testData := base.InstancesTrainTestSplit(instf, 0.60)

				Convey("Fitting and predicting with a Random Forest", func() {
					rf := NewRandomForest(10, 3)
					err = rf.Fit(trainData)
					So(err, ShouldBeNil)

					oldPredictions, err := rf.Predict(testData)
					So(err, ShouldBeNil)

					Convey("Saving the model should work...", func() {
						f, err := ioutil.TempFile(os.TempDir(), "rf")
						So(err, ShouldBeNil)
						err = rf.Save(f.Name())
						defer func() {
							f.Close()
						}()
						So(err, ShouldBeNil)
						Convey("Loading the model should work...", func() {
							newRf := NewRandomForest(10, 3)
							err := newRf.Load(f.Name())
							So(err, ShouldBeNil)
							So(len(newRf.Model.Models), ShouldEqual, 10)
							Convey("Predictions should be the same...", func() {
								newPredictions, err := newRf.Predict(testData)
								So(err, ShouldBeNil)
								So(base.InstancesAreEqual(newPredictions, oldPredictions), ShouldBeTrue)
							})
						})
					})
				})
			})
		})
	})
}
