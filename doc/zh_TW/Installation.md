GoLearn 除了少數依賴於 C 語言以外，皆使用 Go 語言基本函式庫。 使用時，[你需要安裝 Go 1.4 或以上的版本](http://golang.org/doc/install)。

## 安裝

### 系統依賴
* 你需要先安裝好一個適用的編譯器 (在終端機上執行 `g++` 試試)
* GoLearn 使用的 [Gonum BLAS](https://github.com/gonum/blas) 需要事先安裝 OpenBLAS 或是其他類似功能的函式庫。 請依照你的系統跟著[安裝教學](https://github.com/gonum/blas#installation)操作，其中可能會需要安裝一個 C 語言函式庫。

### 安裝 GoLearn 內部依賴的函式庫
安裝完 go 以及系統依賴相關的要件以後，輸入
```bash
go get -t -u -v github.com/sjwhitworth/golearn
```

### 完成安裝
輸入以下指令完成安裝
```bash
cd $GOPATH/src/github.com/sjwhitworth/golearn
go get -t -u -v ./...
```

## 常見問題

在 Linux 或是 Mac OS X 上，你可以使用終端機來確認 Go 是否正確的安裝。 [編譯範例程式試試](http://golang.org/doc/install)來確認是否都安裝成功了。

**你的 `go` 資料夾必須存在於你的家目錄中且是可以被寫入的。**
    如果沒有，請輸入 `cd && mkdir go` 來建立

**你的 `GOPATH` 以及 `PATH` 變數必須正確的設定。**
* 輸入 `echo $GOROOT` 以及 `echo $GOPATH` 來確認。
* `$GOPATH` 變數須包含 `$GOROOT` 以及一個 `bin/` 資料夾。 舉例來說，如果你的 `$GOROOT` 設定在 `/home/sen/go` ，則你的 `$GOPATH` 就需要設定在 `/home/sen/go/bin` 。
* 加入 `export GOROOT=$HOME/go` 以及 `export PATH=$PATH:$GOROOT/bin` 到 Bash 設定檔中以確保所有的變數都被正確的設定。

## 支援情況
<table>
<tr>
<td>作業系統</td><td>Mac OS X 10.8 <br /> Ubuntu 14.04 <br /> OpenSUSE 13.1</td></tr>
<tr><td>Go 版本</td><td>1.2</td></tr>
<tr><td>GoLearn 版本</td><td>當前版本</td></tr>
<tr><td>支援近況 </td><td>當前版本</td>
<tr><td>下一個版本</td><td>參照升級版本</td></tr>
</table>

### Mac OS X 10.8
* 透過 HomeBrew 安裝 BLAS 尚未確認是否可行。