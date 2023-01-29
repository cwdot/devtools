package newconf

import (
	"fmt"
	"os"
	"regexp"
	"text/template"

	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
	"gitter/internal/config"

	"github.com/cwdot/go-stdlib/wood"
)

func Do(g *git.Repository) error {
	headRef, err := g.Head()
	if err != nil {
		wood.Fatal(err)
	}
	name := headRef.Name().Short()

	b := `        - name: {{ .Name }}
          description: {{ .Description }}
{{- if .Links.Jira }}
          links:
            jira: {{ .Links.Jira }}
{{- end }}
`

	td := config.Branch{
		Name:        name,
		Description: "",
		Links: config.BranchLinks{
			Pr:   "",
			Jira: findJira(name),
		},
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

func findJira(branchName string) string {
	if branchName == "" {
		return ""
	}

	r, err := regexp.Compile("([A-Za-z0-9]+-[0-9]+)")
	if err != nil {
		panic(err)
	}

	matches := r.FindStringSubmatch(branchName)
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}
