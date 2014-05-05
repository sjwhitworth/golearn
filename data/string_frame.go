package data

type StringFrame struct {
	labels [][]string
}

func (self *StringFrame) Add(row []string) {
	self.labels = append(self.labels, row)
}

func (self *StringFrame) At(i, j int) string {
	return self.labels[i][j]
}

func (self *StringFrame) Row(i int) []string {
	return self.labels[i]
}

func (self *StringFrame) Col(j int) []string {
	col := make([]string, 0)
	for i := 0; i < self.NRow(); i++ {
		col = append(col, self.labels[i][j])
	}

	return col
}

func (self *StringFrame) NRow() int {
	return len(self.labels)
}

func (self *StringFrame) NCol() int {
	return len(self.labels[0])
}
