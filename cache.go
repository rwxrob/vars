// Copyright 2022 Robert Muhlestein.
// SPDX-License-Identifier: Apache-2.0

// Package cache provides the Bonzai command branch of the same name.
package cache

import (
	"fmt"
	"os"
	"strings"

	Z "github.com/rwxrob/bonzai/z"
	cache "github.com/rwxrob/cache/pkg"
	"github.com/rwxrob/help"
	"github.com/rwxrob/term"
)

var _cache cache.Map

func init() {
	dir, _ := os.UserCacheDir()
	_cache = cache.New()
	_cache.Id = Z.ExeName
	_cache.Dir = dir
	_cache.File = `vars`
	Z.Vars = _cache
}

var Cmd = &Z.Cmd{

	Name:      `var`,
	Summary:   `cache variables in {{ execachedir "vars"}}`,
	Version:   `v0.0.1`,
	Copyright: `Copyright 2021 Robert S Muhlestein`,
	License:   `Apache-2.0`,
	Commands:  []*Z.Cmd{help.Cmd, _init, set, get, _file, data, edit},

	Description: `
		The *{{.Name}}* command provides a cross-platform, persistent
		alternative to environment/system variables. The subcommands are
		designed to be safe and convenient.

		Implementation

		Variables are cached as key=val (property) pairs, one to a line, in
		the {{ execachedir "vars" }} file.

		Key names are automatically prefixed with the Cmd.Path ('{{ .Path
		}}' in this case) which changed depending on where this cache (var)
		branch is composed into the Bonzai command tree.

		Keys must not include an equal sign (=) which is the only line
		delimiter.

		Carriage returns (\r) and line returns (\n) are escaped
		and each line is terminated with a line return (\n).`}

var set = &Z.Cmd{
	Name:    `set`,
	Summary: `safely sets (persists) a cached variable`,
	Description: `
		The *{{.Name}}* command writes the changes to the specified cached
		variable in a way that is reasonably safe for system-wide concurrent
		writes by checking the file for any changes since last right and
		refusing to overwrite if so (much like editing from a Vim session).

		The exact process is as follows:

    1. Save the current time in nanoseconds
    2. Load and parse {{ execachefile "vars" }} into cache.Map
		3. Change the specified value
		4. Check file for changes since saved time, error if changed
		5. Marshal cache.Map and atomically write to file

		`,

	MinArgs: 2,

	Call: func(x *Z.Cmd, args ...string) error {
		path := "."
		if x.Caller != nil {
			path = x.Caller.Path()
		}
		if path != "." {
			path += "."
		}
		val := strings.Join(args[1:], " ")
		return _cache.Set(path+args[0], val)
	},
}

var get = &Z.Cmd{
	Name:    `get`,
	Summary: `gets a cached variable and prints with new line`,
	Description: `
		The *{{.Name}}* command retrieves a cached variable from the cache
		file ({{execachedir "vars"}}) and prints it with a new line to
		standard output. Prints a blank line if not set.`,

	MinArgs: 1,
	//MaxArgs: 1,

	Call: func(x *Z.Cmd, args ...string) error {
		path := "."
		if x.Caller != nil {
			path = x.Caller.Path()
		}
		if path != "." {
			path += "."
		}
		fmt.Println(_cache.Get(path + args[0]))
		return nil
	},
}

var _file = &Z.Cmd{
	Name:     `file`,
	Aliases:  []string{"f"},
	Summary:  `outputs full path to the cached variables file`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(x *Z.Cmd, _ ...string) error {
		fmt.Println(_cache.Path())
		return nil
	},
}

var _init = &Z.Cmd{
	Name:     `init`,
	Aliases:  []string{"i"},
	Summary:  `(re)initializes current variable cache`,
	Commands: []*Z.Cmd{help.Cmd},
	ReqVars:  true, // but fulfills at init() above
	Call: func(x *Z.Cmd, _ ...string) error {
		if term.IsInteractive() {
			r := term.Prompt(`Really initialize %v? (y/N) `, _cache.DirPath())
			if r != "y" {
				return nil
			}
		}
		return Z.Vars.Init()
	},
}

var data = &Z.Cmd{
	Name:    `data`,
	Aliases: []string{"d"},
	Summary: `outputs contents of the cached variables file`,

	Description: `
			The *data* command prints the entire, unobfuscated contents of the
			cached variables file.

			WARNING: Since cached variables regularly includes secrets
			(tokens, keys, etc.) be aware that anyone able to view your screen
			could compromise your security when using this command in front of
			them (presentations, streaming, etc.).`,

	Commands: []*Z.Cmd{help.Cmd},
	Call: func(x *Z.Cmd, _ ...string) error {
		fmt.Print(_cache.Data())
		return nil
	},
}

var edit = &Z.Cmd{
	Name:     `edit`,
	Summary:  `edit variables file ({{execachedir "vars"}}) `,
	Aliases:  []string{"e"},
	Commands: []*Z.Cmd{help.Cmd},

	Description: `
		The *{{.Name}}* command will the configuration file for editing in
		the currently configured editor (in order or priority):

		* $VISUAL
		* $EDITOR
		* vi
		* vim
		* nano

		The edit command hands over control of the currently running process
		to the editor. `,

	Call: func(x *Z.Cmd, _ ...string) error { return _cache.Edit() },
}
