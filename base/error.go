package base

import (
	"runtime/debug"
	"strings"
	"fmt"
	"os"
)

type GoLearnError struct {
	WrappedError error
	CurrentStack string
	Description string
}

func wrapLinesWithTabPrefix(s string) string {
	split := strings.Split(s, "\n")
	stack := make([]string, len(split))
	for i := 0; i < len(split); i++ {
		stack[i] = fmt.Sprintf("\t%s", split[i])
	}
	return strings.Join(stack, "\n")
}

func (g *GoLearnError) Error() string {

	if os.Getenv("GOLEARN_FULL_DEBUG") == "true" {
		return fmt.Sprintf("GoLearnError( %s\n\tCaptured at: %s\n)",
			wrapLinesWithTabPrefix(g.WrappedError.Error()), wrapLinesWithTabPrefix(g.CurrentStack))
	}

	if g.Description == "" {
		fmt.Sprintf("%s", g.WrappedError)
	}
	return fmt.Sprintf("GoLearnError( %s: %s )", g.Description, g.WrappedError)
}

func (g *GoLearnError) attachFormattedStack() {
	stackString := string(debug.Stack())
	stackFrames := strings.Split(stackString, "\n")

	stackFmt := make([]string, 0)
	for i := 3; i < len(stackFrames); i++ {
		if strings.Contains(stackFrames[i], "golearn") {
			stackFmt = append(stackFmt, stackFrames[i])
		}
	}
	stackOut := "<invalid>"
	if len(stackFmt) > 0 {
		stackOut = strings.Join(stackFmt, "\t\t\n")
	}
	g.CurrentStack = stackOut
}

func DescribeError(description string, err error) error {
	ret := &GoLearnError{}
	ret.WrappedError = err
	ret.attachFormattedStack()
	ret.Description = description
	return ret
}

func WrapError(err error) error {
	ret := &GoLearnError{}
	ret.WrappedError = err
	ret.attachFormattedStack()
	return ret
}