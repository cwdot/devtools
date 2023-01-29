package opbridge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/pkg/errors"

	"github.com/cwdot/go-stdlib/wood"
)

func List(tag string) ([]OpEntry, error) {
	//op item list --tags "1px/mac4" --format json
	cmd := exec.Command("op", "item", "list", "--tags", fmt.Sprintf(`"%s"`, tag), "--format", "json")

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
