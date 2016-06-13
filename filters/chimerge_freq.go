package filters

import (
	"fmt"
)

// FrequencyTableEntry is a struct holding a value and a map of frequency
type FrequencyTableEntry struct {
	Value     float64
	Frequency map[string]int
}

func (t *FrequencyTableEntry) String() string {
	return fmt.Sprintf("%.2f %+v", t.Value, t.Frequency)
}
