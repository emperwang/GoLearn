package trending

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestGrepNumber(t *testing.T) {

	str := "12345 stars of today"

	pattern, _ := regexp.Compile("[0-9]{1,}")

	// search substring
	substr := pattern.FindString(str)

	if substr != "12345" {
		t.Errorf("get error value: %s", substr)
	}
	t.Logf(substr)
}

func TestColorOutput(t *testing.T) {
	// cyan
	cyan := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("00FFFF"))

	fmt.Println(cyan.Render("123"))
}
