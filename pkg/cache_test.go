package cache_test

import (
	"fmt"
	"os"

	cache "github.com/rwxrob/cache/pkg"
	"github.com/rwxrob/fs/file"
)

func ExampleCacheMap() {
	m := cache.New()
	m.Id = `foo`
	m.Dir = `testdata`
	m.File = `cache`
	fmt.Println(m.Path())
	fmt.Println(m.DirPath())
	// Output:
	// testdata/foo/cache
	// testdata/foo
}

func ExampleCacheMap_Init() {

	m := cache.New()
	m.Id = `foo`
	m.Dir = `testdata`
	m.File = `cache`

	defer func() { os.RemoveAll(m.DirPath()) }()

	m.Init()
	fmt.Println(file.Exists(`testdata/foo/cache`))

	// Output:
	// true
}

func ExampleCacheMap_Set() {

	m := cache.New()
	m.Id = `foo`
	m.Dir = `testdata`
	m.File = `cache`

	defer func() { os.RemoveAll(m.DirPath()) }()

	m.Init()
	if err := m.Set("some", "thing\nhere"); err != nil {
		fmt.Println(err)
	}
	byt, _ := os.ReadFile(`testdata/foo/cache`)
	fmt.Println(string(byt) == `some=thing\nhere`+"\n")

	// Output:
	// true
}

func ExampleCacheMap_Get() {

	m := cache.New()
	m.Id = `foo`
	m.Dir = `testdata`
	m.File = `cache`

	defer func() { os.RemoveAll(m.DirPath()) }()

	m.Init()
	if err := m.Set("some", "thing\nhere"); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%q\n", m.Get(`some`))

	// Output:
	// "thing\nhere"
}

func ExampleCacheMap_UnmarshalText() {

	in := `
some=thing here
another=one over here
`

	m := cache.New()
	m.UnmarshalText([]byte(in))
	fmt.Println(len(m.M))
	fmt.Println(m.M["some"])
	fmt.Println(m.M["another"])

	// Output:
	// 2
	// thing here
	// one over here
}

func ExampleCacheMap_MarshalText() {

	m := cache.New()
	m.M["some"] = "thing here"
	m.M["another"] = "one\rhere\nbut all good"

	byt, err := m.MarshalText()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(byt))

	// Ordered Output:
	// some:thing here
	// another:one\rhere\nbut all good
}
