// Copyright 2022 foo Authors
// SPDX-License-Identifier: Apache-2.0

package foo

// Go treats all files as if they are, more or less, in the same large
// file. Create separate files to help you and others find the code you
// need quickly.

import (
	"log"

	Z "github.com/rwxrob/bonzai/z"
)

var own = &Z.Cmd{
	Name: `own`,
	Call: func(caller *Z.Cmd, none ...string) error {
		log.Print("I'm in my own file.")
		return nil
	},
}
