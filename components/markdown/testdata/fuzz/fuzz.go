package fuzzing

import "github.com/rprtr258/tea/components/markdown"

func Fuzz(data []byte) int {
	_, err := markdown.RenderBytes(data, markdown.DarkStyle)
	if err != nil {
		return 0
	}
	return 1
}
