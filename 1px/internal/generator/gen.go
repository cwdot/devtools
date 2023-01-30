package generator

import (
	"fmt"
	"os"
	"time"

	"1px/internal/opbridge"
	"github.com/pkg/errors"

	"github.com/cwdot/go-stdlib/wood"
)

func Write(entries []Entry, output string) error {
	f, err := os.CreateTemp("", "prefix")
	if err != nil {
		return errors.Wrap(err, "failed to open temp")
	}
	defer os.Remove(f.Name())

	_, err = f.WriteString(fmt.Sprintf("# %s\n", time.Now().Format(time.RFC3339)))
	if err != nil {
		return errors.Wrap(err, "failed to write header")
	}

	for _, entry := range entries {
		if entry.Comment != "" {
			_, err = f.WriteString(fmt.Sprintf("# %s\n", entry.Comment))
			if err != nil {
				return errors.Wrap(err, "failed to write entry")
			}
		}
		_, err = f.WriteString(fmt.Sprintf("%s=%s\n", entry.Key, entry.Value))
		if err != nil {
			return errors.Wrap(err, "failed to write entry")
		}

		wood.Debugf("%s=%s\n", entry.Key, entry.Value)
	}

	err = opbridge.Inject(f.Name(), output)
	if err != nil {
		return errors.Wrap(err, "failed to inject credentials")
	}

	if err := f.Close(); err != nil {
		return errors.Wrap(err, "failed to close file")
	}

	wood.Infof("Exported %d credentials to %s", len(entries), output)

	return nil
}

type Entry struct {
	Key     string
	Value   string
	Comment string
}
