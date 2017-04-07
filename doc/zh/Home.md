GoLearn
=======

<img src="http://talks.golang.org/2013/advconc/gopherhat.jpg" width=125><br>

欢迎阅读 GoLearn 中文文档. GoLearn 是一个 "开箱即用" 的机器学习库，使用并为 Go 语言编写。 machine learning library written for and in Go. [GoLearn Github](https://github.com/sjwhitworth/golearn).


[安装](Installation.md) |
[读入数据](Instances.md) [注意: 非最新版] | 
[过滤](Filtering.md) | 
分类 ([KNN](Classification/KNN.md) | [Trees](Classification/Trees.md) | [liblinear](Classification/liblinear.md)) | [回归分析](Classification/Regression.md) 

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

## 发展 
[贡献](Contributing.md)

### 愿望清单

* 最大化期望
* 本地引导树结构
* 支持时间序列处理
* 支持 `DenseInstances` 用于磁盘备份
* SoftMax 神经网络
* 原始深度学习（递归神经网络）
* 循环神经网络
* 图像处理
* 关系和指针属性
* 支持稀疏的二进制空间
* 支持任意插入和删除属性