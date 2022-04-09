// Copyright 2022 foo Authors
// SPDX-License-Identifier: Apache-2.0

// Package foo provides the Bonzai command branch of the same name.
package foo

import (
	"log"

	"github.com/rwxrob/bonzai/comp"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/config"
	foo "github.com/rwxrob/foo/pkg"
	"github.com/rwxrob/help"
)

var Cmd = &Z.Cmd{

	Name:      `foo`,
	Summary:   `just a sample foo command`,
	Usage:     `[B|bar|own|h|help]`,
	Version:   `v2.2.8`,
	Copyright: `Copyright 2021 Robert S Muhlestein`,
	License:   `Apache-2.0`,
	Commands:  []*Z.Cmd{help.Cmd, config.Cmd, Bar, own, pkgfoo},

	Description: `
		The foo commands do foo stuff. You can start the description here
		and wrap it to look nice and it will just work. Descriptions are
		written in BonzaiMark™, a simplification of CommonMark that that
		mostly follows Go documentation guidelines`,

	Other: []Z.Section{
		{`foo`, `something about foo`},
		{`another`, `something about another command`},
	},

	// no Call since has Commands, if had Call would only call if
	// commands didn't match
}

// Commands can be grouped into the same file or separately, public or
// private. Public let's others compose specific subcommands (foo.Bar),
// private just keeps it composed and only available within this Bonzai
// command.

// Aliases are not commands but will be replaced by their target names.

var Bar = &Z.Cmd{
	Name:     `bar`,
	Aliases:  []string{"B", "notbar"}, // to make a point
	Commands: []*Z.Cmd{help.Cmd, file},

	// Call first-class functions can be highly detailed, refer to an
	// existing function someplace else, or can call high-level package
	// library functions. Developers are encouraged to consider well where
	// they maintain the core logic of their applications. Often, it will
	// not be here within the Z.Cmd definition. One use case for
	// decoupled first-class Call functions is when creating multiple
	// binaries for different target languages. In such cases this
	// Z.Cmd definition is essentially just a wrapper for
	// documentation and other language-specific embedded assets.

	Call: func(_ *Z.Cmd, _ ...string) error { // note conventional _
		log.Printf("would bar stuff")
		return nil
	},
}

// Different completion methods are be set including the expected
// standard ones from bash and other shells. Not that these completion
// methods only work if the shell supports completion (including
// the Bonzai Shell, which can be set as the default Cmd to provide rich
// shell interactivity where normally no shell is available, such as in
// FROM SCRATCH containers that use a Bonzai tree as the core binary).

var file = &Z.Cmd{
	Name:      `file`,
	Commands:  []*Z.Cmd{help.Cmd},
	Completer: comp.File,
	Call: func(x *Z.Cmd, args ...string) error {
		if len(args) == 0 {
			return x.UsageError()
		}
		log.Printf("would show file information about %v", args[0])
		return nil
	},
}

// When combining a high-level package library with a Bonzai command it
// is customary to create a pkg directory to avoid cyclical package
// import dependencies.

var pkgfoo = &Z.Cmd{
	Name: `pkgfoo`,
	Call: func(_ *Z.Cmd, _ ...string) error {
		foo.Foo()
		return nil
	},
}
