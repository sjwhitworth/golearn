GoLearn
=======

<img src="http://talks.golang.org/2013/advconc/gopherhat.jpg" width=125><br>
[![GoDoc](https://godoc.org/github.com/sjwhitworth/golearn?status.png)](https://godoc.org/github.com/sjwhitworth/golearn)
[![Build Status](https://travis-ci.org/sjwhitworth/golearn.png?branch=master)](https://travis-ci.org/sjwhitworth/golearn)<br>

歡迎閱讀 GoLearn 中文文檔。 GoLearn 是一個 "開箱即用" 的 Go 語言機器學習函式庫。

[安装](Installation.md) |
[讀入數據](Instances.md) [注意: 非最新版] | 
[過濾](Filtering.md) | 
分類 ([KNN](Classification/KNN.md) | [Trees](Classification/Trees.md) | [liblinear](Classification/liblinear.md)) | [回歸分析](Classification/Regression.md) 

## 快速開始
* [使用 csv 文件](CSVFiles.md)
* [切割數據為訓練資料以及測試資料](TrainTestSplit.md)

## 程式碼節錄與範例
* [數據排序](Instances.md)
* [讀入CSV文件](CSVFiles.md)
* [直方圖合併 Histogram binning](Filtering.md)
* [離散數據合併 Chi-Merge binning](Filtering.md)
* [設定一個 `FloatAttribute` 物件的精度](FloatAttributePrecision.md)
* [增加特徵](AddingAttributes.md)
* [檢視特徵值](AttributeSpecifications.md)
* [製作一個特製的 `DataGrid`](CustomDataGrids.md)

## 發展
[貢獻](Contributing.md)

### 願望清單

* 最大化期望
* 本地引導樹結構
* 支援時間序列處理
* 支援備份 `DenseInstances` 
* SoftMax 神經網路
* 原始深度學習（遞迴神經網路）
* 循環神經網路
* 圖像處理
* 支援相關屬性和指標属性
* 支援稀疏矩陣
* 支援任意插入及刪除屬性
