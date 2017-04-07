GoLearn
=======

<img src="http://talks.golang.org/2013/advconc/gopherhat.jpg" width=125><br>

欢迎阅读 GoLearn 中文文档. GoLearn 是一个 "开箱即用" 的机器学习库，使用并为 Go 语言编写。 machine learning library written for and in Go. [GoLearn Github](https://github.com/sjwhitworth/golearn).


[安装](Installation.md) |
[读入数据](Instances.md) [注意: 非最新版] | 
[过滤](Filtering.md) | 
分类 ([KNN](Classification/KNN.md) | [Trees](Classification/Trees.md) | [liblinear](Classification/liblinear.md)) | [Regression](Classification/Regression.md) 

## 快速开始
* [使用 csv 文件](CSVFiles.md)
* [数据分为训练和测试集](TrainTestSplit.md) (未编写)

## 代码摘录和示例
* [数据排序](Instances.md)
* [读入CSV文件](CSVFiles.md)
* [合并直方图](Filtering.md)
* [离散数据合并](Filtering.md)
* [设置一个 `FloatAttribute` 精度](FloatAttributePrecision.md)
* [添加特征](AddingAttributes.md)
* [检索属性值](AttributeSpecifications.md)
* [实现一个定制的 `DataGrid`](CustomDataGrids.md)

## Future 
[贡献](Contributing.md)

### The wish-list
* Expectation maximisation
* Native guided tree structures
* Support for time series processing
* Support for disk-backed `DenseInstances`
* SoftMax neural networks
* Deep-learning primitives (recursive neural networks)
* Recurrent neural networks
* Image manipulation
* Relational and pointer Attributes
* Support for sparse binary spaces
* Support for arbitrary insertion and deletion of Attributes