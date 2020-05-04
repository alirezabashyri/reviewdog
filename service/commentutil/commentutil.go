package commentutil

import (
	"fmt"
	"log"

	"github.com/reviewdog/reviewdog"
)

// `path` to `position`(Lnum for new file) to comment `body`s
type PostedComments map[string]map[int][]string

// IsPosted returns true if a given comment has been posted in code review service already,
// otherwise returns false. It sees comments with same path, same position,
// and same body as same comments.
func (p PostedComments) IsPosted(c *reviewdog.Comment, lineNum int) bool {
	if _, ok := p[c.Path]; !ok {
		return false
	}
	bodies, ok := p[c.Path][lineNum]
	if !ok {
		return false
	}
	for _, body := range bodies {
		if body == CommentBody(c) {
			return true
		}
	}
	return false
}

// AddPostedComment adds a posted comment.
func (p PostedComments) AddPostedComment(path string, lineNum int, body string) {
	if _, ok := p[path]; !ok {
		p[path] = make(map[int][]string)
	}
	if _, ok := p[path][lineNum]; !ok {
		p[path][lineNum] = make([]string, 0)
	}
	p[path][lineNum] = append(p[path][lineNum], body)
}

// DebugLog outputs posted comments as log for debugging.
func (p PostedComments) DebugLog() {
	for filename, f := range p {
		for line := range f {
			log.Printf("[debug] posted: %s:%d", filename, line)
		}
	}
}

// BodyPrefix is prefix text of comment body.
const BodyPrefix = `<sub>reported by [reviewdog](https://github.com/reviewdog/reviewdog) :dog:</sub>`

// CommentBody creates comment body text.
func CommentBody(c *reviewdog.Comment) string {
	tool := ""
	if c.ToolName != "" {
		tool = fmt.Sprintf("**[%s]** ", c.ToolName)
	}
	return tool + BodyPrefix + "\n" + c.Body
}
