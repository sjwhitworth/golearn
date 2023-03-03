```go
// Create a new, empty DenseInstances
newInst := base.NewDenseInstances()

// Create some Attributes 
attrs := make([]base.Attribute, 2)
attrs[0] = base.NewFloatAttribute("Arbitrary Float Quantity")
attrs[1] = new(base.CategoricalAttribute)
attrs[1].SetName("Class")

// Add the attributes
newSpecs := make([]base.AttributeSpec, len(attrs))
for i, a := range attrs {
    newSpecs[i] = newInst.AddAttribute(a)
}
```

**Trecho de código: Adicionando dois novos atributos a um `UpdatableDataGrid` em branco**

* Os atributos só podem ser adicionados a `UpdatableDataGrids` vazios (ou seja, antes de chamar `Extend` para alocar memória)..
