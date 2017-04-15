package kdtree

import(
  "sort"
  "errors"
)

type node struct{
  feature   int
  value     []float64
  left      *node
  right     *node
}

// Tree is a kdtree.
type Tree struct{
  firstDiv  *node
}

// New return a Tree pointer.
func New()*Tree{
  return &Tree{}
}

// Build builds the kdtree with specific data.
func (t *Tree) Build(data [][]float64)err{
  if len(data)==0 {
    return errors.New("no input data")
  }
  size := len(data[0])
  for _, v := range data {
    if len(v) != size {
      return errors.New("amounts of features are not the same")
    }
  }

  t.firstDiv = t.buildHandle(data, 0)

  return nil
}

// buildHandle builds the kdtree recursively.
func (t *tree) buildHandle(data [][]float64, featureIndex int)*node{
  n := &node{feature:featureIndex}

  sort.Slice(data, func(i, j int)bool{
    return data[i][featureIndex]<data[j][featureIndex]
  })
  middle:= len(data)/2

  n.value = make([]float64, len(data[middle]))
  copy(n.value, data[middle])

  divPoint := middle
  for i:=middle+1 ; i<len(data) ; i++ {
    if data[i][featureIndex] == data[middle][featureIndex] {
      divPoint=i;
    }else{
      break
    }
  }

  if divPoint==1 {
    n.Left = &node{feature:-1}
    n.Left.value = make([]float64, len(data[0]))
    copy(n.Left.value, data[0])
  }else{
    n.Left = t.buildHandle(data[:divPoint])
  }

  if divPoint==(len(data)-2) {
    n.Right = &node{feature:-1}
    n.Right.value = make([]float64, len(data[divPoint+1]))
    copy(n.Right.value, data[divPoint+1])
  }else if divPoint!=(len(data)-1){
    n.Right = t.buildHandle(data[divPoint+1:])
  }

  return n
}

