package data

// Used for storing matrix-like strings. For example, if you have multiple labels associated with each instance, then you probably need a matrix-like struct to store whole your labels.
type StringFrame struct {
	labels [][]string
}

// Add a row of string to StringFrame
func (self *StringFrame) Add(row []string) {
	self.labels = append(self.labels, row)
}

// Get the string at (i, j)
func (self *StringFrame) At(i, j int) string {
	return self.labels[i][j]
}

// Get the i-th row of string
func (self *StringFrame) Row(i int) []string {
	return self.labels[i]
}

// Get the j-th column of string
func (self *StringFrame) Col(j int) []string {
	col := make([]string, 0)
	for i := 0; i < self.NRow(); i++ {
		col = append(col, self.labels[i][j])
	}

	return col
}

// Number of rows
func (self *StringFrame) NRow() int {
	return len(self.labels)
}

// Number of columns
func (self *StringFrame) NCol() int {
	return len(self.labels[0])
}
