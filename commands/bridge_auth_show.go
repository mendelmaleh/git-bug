package commands

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/MichaelMure/git-bug/bridge/core/auth"
	"github.com/MichaelMure/git-bug/cache"
	"github.com/MichaelMure/git-bug/util/interrupt"
)

func newBridgeAuthShow() *cobra.Command {
	env := newEnv()

	cmd := &cobra.Command{
		Use:     "show",
		Short:   "Display an authentication credential.",
		PreRunE: loadRepo(env),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBridgeAuthShow(env, args)
		},
		Args: cobra.ExactArgs(1),
	}

	return cmd
}

func runBridgeAuthShow(env *Env, args []string) error {
	backend, err := cache.NewRepoCache(env.repo)
	if err != nil {
		return err
	}
	defer backend.Close()
	interrupt.RegisterCleaner(backend.Close)

	cred, err := auth.LoadWithPrefix(env.repo, args[0])
	if err != nil {
		return err
	}

	env.out.Printf("Id: %s\n", cred.ID())
	env.out.Printf("Target: %s\n", cred.Target())
	env.out.Printf("Kind: %s\n", cred.Kind())
	env.out.Printf("Creation: %s\n", cred.CreateTime().Format(time.RFC822))

	switch cred := cred.(type) {
	case *auth.Token:
		env.out.Printf("Value: %s\n", cred.Value)
	}

	env.out.Println("Metadata:")

	meta := make([]string, 0, len(cred.Metadata()))
	for key, value := range cred.Metadata() {
		meta = append(meta, fmt.Sprintf("    %s --> %s\n", key, value))
	}
	sort.Strings(meta)

	env.out.Print(strings.Join(meta, ""))

	return nil
}
