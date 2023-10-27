package newconf

import (
	"fmt"
	"os"
	"text/template"

	"github.com/cwdot/go-stdlib/wood"
	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"

	"gitter/internal/config"
	"gitter/internal/providers/jirap"
)

func Do(g *git.Repository) error {
	headRef, err := g.Head()
	if err != nil {
		wood.Fatal(err)
	}
	name := headRef.Name().Short()

	b := `        - name: {{ .Name }}
          description: {{ .Description }}
          pr: {{ .Links.Pr }}
          jira: {{ .Links.Jira }}
`

	td := config.Branch{
		Name:        name,
		Description: "",
		Pr:          "",
		Jira:        jirap.Extract("", name),
	}

	t, err := template.New("branch").Parse(b)
	if err != nil {
		return errors.Wrap(err, "failed to parse template")
	}

	err = t.Execute(os.Stdout, td)
	if err != nil {
		return errors.Wrap(err, "failed to execute template")
	}

	conf, err := config.DefaultConfigFile()
	if err != nil {
		return errors.Wrap(err, "failed to get config")
	}
	fmt.Println()
	fmt.Println(conf.Location)

	return nil
}
