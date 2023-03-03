Attribute specifications (`AttributeSpec`) descrevem a localização de armazenamento dos dados digitados por um dado `Attribute`, e permitem que os atributos sejam reordenados e filtrados sem gerar comportamentos inesperados.

**Trecho de código: resolvendo a especificação de um atributo individual**

```go
package main

import (
	"fmt"
	"github.com/sjwhitworth/golearn/base"
)

func main() {

	// Read the CSV file
	data, err := base.ParseCSVToInstances("iris_headers.csv", true)
	// Error check
	if err != nil {
		panic(fmt.Sprintf("Couldn't load CSV file (error %s)", err))
	}

	// Resolve an Attribute
	as, err := data.GetAttribute(base.NewFloatAttribute("Sepal length"))
	if err != nil {
		panic(fmt.Sprintf("Couldn't resolve AttributeSpec (error %s)", err))
	}

	// Print first column
	asArr := []base.AttributeSpec{as}
	data.MapOverRows(asArr, func(row [][]byte, i int) (bool, error) {
		fmt.Println(base.UnpackBytesToFloat(row[0]))
		return true, nil
	})

}
```



Por conveniência, `base` também fornece [`ResolveAttributes`](https://godoc.org/github.com/sjwhitworth/golearn/base#ResolveAttributes), que retorna partes de `AttributeSpecs` para os atributos especificados.

**Trecho de código: resolvendo as especificações de todos os non-class atributos**

```go
package main

import (
	"fmt"
	"github.com/sjwhitworth/golearn/base"
)

func main() {

	// Read the CSV file
	data, err := base.ParseCSVToInstances("iris_headers.csv", true)
	// Error check
	if err != nil {
		panic(fmt.Sprintf("Couldn't load CSV file (error %s)", err))
	}

	// Resolve all non-class Attributes
	asArr := base.ResolveAttributes(data, base.NonClassAttributes(data))
	// (ResolveAllAttributes gets AttributeSpecs for every attribute)

	// Print non-class data
	data.MapOverRows(asArr, func(row [][]byte, i int) (bool, error) {
		for _, b := range row {
			fmt.Print(base.UnpackBytesToFloat(b), " ")
		}
		fmt.Println()
		return true, nil
	})

}
```
