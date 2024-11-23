package paho

import (
	"strings"

	"github.com/pkg/errors"
)

type Topic struct {
	raw      string
	resolved string

	placeholders map[string]string
	needsResolve bool
}

func (t *Topic) String() string {
	if !t.needsResolve {
		return t.resolved
	}

	t.resolved = t.raw[:]
	for k, v := range t.placeholders {
		if k == "~" { // special case
			t.resolved = strings.ReplaceAll(t.resolved, "~", v)
			continue
		}
		t.resolved = strings.ReplaceAll(t.resolved, "{{"+k+"}}", v)
	}

	t.needsResolve = false
	return t.resolved
}

//goland:noinspection ALL
func (t *Topic) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}

	t.raw = value
	t.needsResolve = true
	return nil
}

func (t *Topic) MarshalYAML() ([]byte, error) {
	return []byte(t.raw), nil
}

func (t *Topic) IsEmpty() bool {
	return t.raw == ""
}

func (t *Topic) SetPlaceholders(m map[string]string) *Topic {
	t.placeholders = m
	t.needsResolve = true
	return t
}

func (t *Topic) Validate() error {
	if strings.HasPrefix(t.String(), "~") {
		return errors.New("topic cannot start with ~")
	}
	if strings.Contains(t.resolved, "{{") {
		return errors.New("topic contains unresolved placeholders")
	}
	return nil
}

func (t *Topic) SetPlaceholder(s string, topic string) *Topic {
	if t.placeholders == nil {
		t.placeholders = make(map[string]string)
	}
	t.placeholders[s] = topic
	t.needsResolve = true
	return t
}

func (t *Topic) Copy() *Topic {
	return &Topic{raw: t.raw, resolved: t.resolved, placeholders: t.placeholders, needsResolve: t.needsResolve}
}

func NewTopic(topic string, ph map[string]string) *Topic {
	t := &Topic{raw: topic, resolved: topic}
	t.SetPlaceholders(ph)
	return t
}
