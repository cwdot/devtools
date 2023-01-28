package opbridge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

func List(tag string) ([]OpEntry, error) {
	//op item list --tags "1px/mac4" --format json
	cmd := exec.Command("op", "item", "list", "--tags", fmt.Sprintf(`"%s"`, tag), "--format", "json")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	var entries []OpEntry
	err = json.Unmarshal(out.Bytes(), &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

// Inject replaces placeholders in the template and writes the result into the output file
func Inject(template string, output string) error {
	cmd := exec.Command("op", "inject", "-f", "-i", template, "-o", output)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Println(out.String())
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
