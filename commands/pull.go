package commands

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/MichaelMure/git-bug/cache"
	"github.com/MichaelMure/git-bug/entity"
	"github.com/MichaelMure/git-bug/util/interrupt"
)

func newPullCommand() *cobra.Command {
	env := newEnv()

	cmd := &cobra.Command{
		Use:     "pull [<remote>]",
		Short:   "Pull bugs update from a git remote.",
		PreRunE: loadRepo(env),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPull(env, args)
		},
	}

	return cmd
}

func runPull(env *Env, args []string) error {
	if len(args) > 1 {
		return errors.New("Only pulling from one remote at a time is supported")
	}

	remote := "origin"
	if len(args) == 1 {
		remote = args[0]
	}

	backend, err := cache.NewRepoCache(env.repo)
	if err != nil {
		return err
	}
	defer backend.Close()
	interrupt.RegisterCleaner(backend.Close)

	env.out.Println("Fetching remote ...")

	stdout, err := backend.Fetch(remote)
	if err != nil {
		return err
	}

	env.out.Println(stdout)

	env.out.Println("Merging data ...")

	for result := range backend.MergeAll(remote) {
		if result.Err != nil {
			env.err.Println(result.Err)
		}

		if result.Status != entity.MergeStatusNothing {
			env.out.Printf("%s: %s\n", result.Id.Human(), result)
		}
	}

	return nil
}
