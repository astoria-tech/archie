package msgs

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

type message struct {
	Input  []string
	Output []string
	re     []*regexp.Regexp
	usable bool
}

func (m *message) output() (string, error) {
	if len(m.Output) < 1 {
		return "", errors.New("No output messages")

	}
	return fmt.Sprintf(m.Output[rand.Intn(len(m.Output))]), nil
}

//Messages collection of messages
type Messages map[string]*message

var (
	messages *Messages
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

//LoadMessages preloads regex for matching
func LoadMessages(msgs *Messages) {
	messages = msgs
	for mType, m := range *messages {

		// No need to compile for unknown
		if mType == "unknown" {
			continue
		}

		// Don't compile if we don't have at least 1 input
		if len(m.Input) < 1 {
			fmt.Printf("Error: Message \"%s\" does not have input messages\n", mType)
			continue
		}

		// Don't compile if we don't have at leat 1 output
		if len(m.Output) < 1 {
			fmt.Printf("Error: Message \"%s\" does not have output messages\n", mType)
			continue
		}

		for _, i := range m.Input {
			pattern := "(?i)^" + strings.Replace(i, " ", `[[:space:][:punct:]]+`, -1) + "[[:space:][:punct:]]*$"
			re, err := regexp.Compile(pattern)
			if err != nil {
				fmt.Printf("Error: Message \"%s\" had RE compile error: %s\n", mType, err)
			}
			m.re = append(m.re, re)
		}
		m.usable = true
	}
}

//Response generates a reply to a message
func Response(msg string) (string, error) {
	for mType, m := range *messages {

		// Unknown and unusable shouldb e skipped
		if mType == "unknown" || !m.usable {
			continue
		}

		for _, re := range m.re {
			msg = strings.TrimSpace(msg)
			if re.MatchString(msg) {
				retOut, err := m.output()
				if err != nil {
					e := fmt.Errorf("getting \"%s\" output: %s", mType, err)
					return "", e
				}
				return retOut, nil
			}
		}
	}

	// ASSERT: No usable msgs found, return unknown
	retOut, err := (*messages)["unknown"].output()
	if err != nil {
		e := fmt.Errorf("getting \"unknown\" output: %s", err)
		return "", e
	}
	return retOut, nil
}
