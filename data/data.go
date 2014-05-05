/* Data - consists of helper functions for parsing different data formats */

package data

import "github.com/gonum/matrix/mat64"

type DataFrame struct {
	Headers  []string
	Labels   []string
	Values   *mat64.Dense // We first focus on numeric values for now
	NRow     int
	NFeature int
	NLabel   int
}
