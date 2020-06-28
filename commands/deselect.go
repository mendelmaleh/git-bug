package commands

import (
	"github.com/spf13/cobra"

	"github.com/MichaelMure/git-bug/cache"
	"github.com/MichaelMure/git-bug/commands/select"
	"github.com/MichaelMure/git-bug/util/interrupt"
)

func newDeselectCommand() *cobra.Command {
	env := newEnv()

	cmd := &cobra.Command{
		Use:   "deselect",
		Short: "Clear the implicitly selected bug.",
		Example: `git bug select 2f15
git bug comment
git bug status
git bug deselect
`,
		PreRunE: loadRepo(env),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeselect(env)
		},
	}

	return cmd
}

func runDeselect(env *Env) error {
	backend, err := cache.NewRepoCache(env.repo)
	if err != nil {
		return err
	}
	defer backend.Close()
	interrupt.RegisterCleaner(backend.Close)

	err = _select.Clear(backend)
	if err != nil {
		return err
	}

	return nil
}
