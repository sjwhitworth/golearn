package filters

import (
	"fmt"
)

type FrequencyTableEntry struct {
	Value     float64
	Frequency map[string]int
}

func (t *FrequencyTableEntry) String() string {
	return fmt.Sprintf("%.2f %s", t.Value, t.Frequency)
}
