package opbridge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/cwdot/stdlib-go/wood"
)

func List(tag string) ([]OpEntry, error) {
	//op item list --tags "1px/mac4" --format json

	cmd := exec.Command("op", "item", "list", "--tags", fmt.Sprintf(`"%s"`, tag), "--long", "--format", "json")

	wood.Debugf("Invoked op item list: %s", cmd.String())

	var outs bytes.Buffer
	var errs bytes.Buffer
	cmd.Stdout = &outs
	cmd.Stderr = &errs

	err := cmd.Run()
	if err != nil {
		os.Stderr.Write(outs.Bytes())
		os.Stderr.Write(errs.Bytes())
		return nil, errors.Wrap(err, "list failed")
	}

	var entries []OpEntry
	err = json.Unmarshal(outs.Bytes(), &entries)
	if err != nil {
		return nil, err
	}

	wood.Debugf("Invoked op item list: %s", outs.String())
	wood.Infof("Found %d credentials matching criteria", len(entries))
	return entries, nil
}

func Item(id string, fieldNames ...string) (map[string]Field, error) {
	//op item XXX --fields "username,credentials" --format json

	cmd := exec.Command("op", "item", "get", id, "--fields", strings.Join(fieldNames, ","), "--format", "json")

	wood.Debugf("Invoked op item list: %s", cmd.String())

	var outs bytes.Buffer
	var errs bytes.Buffer
	cmd.Stdout = &outs
	cmd.Stderr = &errs

	err := cmd.Run()
	if err != nil {
		os.Stderr.Write(outs.Bytes())
		os.Stderr.Write(errs.Bytes())
		return nil, errors.Wrap(err, "list failed")
	}

	var fieldList []Field

	if len(fieldNames) == 1 {
		var field Field
		err = json.Unmarshal(outs.Bytes(), &field)
		if err != nil {
			return nil, err
		}
		fieldList = append(fieldList, field)
	} else {
		err = json.Unmarshal(outs.Bytes(), &fieldList)
		if err != nil {
			return nil, err
		}
	}

	fields := make(map[string]Field)
	for _, field := range fieldList {
		fields[field.Id] = field
	}

	return fields, nil
}

// Inject replaces placeholders in the template and writes the result into the output file
func Inject(template string, output string) error {
	cmd := exec.Command("op", "inject", "-f", "-i", template, "-o", output)

	var outs bytes.Buffer
	var errs bytes.Buffer
	cmd.Stdout = &outs
	cmd.Stderr = &errs

	err := cmd.Run()
	if err != nil {
		os.Stderr.Write(outs.Bytes())
		os.Stderr.Write(errs.Bytes())
		return errors.Wrap(err, "inject failed")
	}

	wood.Debugf("Invoked op inject: %s", outs.String())
	return nil
}

type OpEntry struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Tags    []string `json:"tags"`
	Version int      `json:"version"`
	Vault   struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"vault"`
	Category              string    `json:"category"`
	LastEditedBy          string    `json:"last_edited_by"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	AdditionalInformation string    `json:"additional_information"`
}

type Field struct {
	Id        string `json:"id"`
	Type      string `json:"type"`
	Purpose   string `json:"purpose"`
	Label     string `json:"label"`
	Value     string `json:"value"`
	Reference string `json:"reference"`
}
