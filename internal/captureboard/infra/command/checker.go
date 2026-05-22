package command

import (
	"fmt"
	"os/exec"
	"strings"
)

type Checker struct{}

func NewChecker() *Checker {
	return &Checker{}
}

func (c *Checker) Ensure(names ...string) error {
	var missing []string
	for _, name := range names {
		if _, err := exec.LookPath(name); err != nil {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("%s 명령어를 찾을 수 없습니다", strings.Join(missing, ", "))
	}

	return nil
}
