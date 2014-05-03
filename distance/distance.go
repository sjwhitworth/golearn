package distance

import (
  "math"
)

func EucledianDistance(p1 []float64, p2 []float64) float64 {
  sum := float64(0)
  for i := 0; i < cap(p1); i++ {
    sum += math.Pow(p1[i] - p2[i], 2)
  }
  return math.Sqrt(sum);
}

func cranberraDistanceStep(num float64, denom float64) float64 {
  if num == float64(0) && denom == float64(0) {
    return float64(0)
  } else {
    return num/denom
  }
}

func CranberraDistance(p1 []float64, p2 []float64) float64 {
  sum := float64(0)

  for i := 0; i < cap(p1); i++ {
    num := math.Abs(p1[i] - p2[i])
    denom := math.Abs(p1[i]) + math.Abs(p2[i])
    sum += cranberraDistanceStep(num, denom)
  }
  return sum;
}

func ChebyshevDistance(p1 []float64, p2 []float64) float64 {
  max := float64(0)

  for i := 0; i < cap(p1); i++ {
    max = math.Max(max, math.Abs(p1[i] - p2[i]))
  }

  return max
}


func ManhattanDistance(p1 []float64, p2 []float64) float64 {
  sum := float64(0)

  for i := 0; i < cap(p1); i++ {
    sum += math.Abs(p1[i] - p2[i])
  }

  return sum
}
