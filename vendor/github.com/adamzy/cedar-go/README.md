# cedar-go [![GoDoc](https://godoc.org/github.com/adamzy/cedar-go?status.svg)](https://godoc.org/github.com/adamzy/cedar-go)

Package `cedar-go` implementes double-array trie.

It is a [Golang](https://golang.org/) port of [cedar](http://www.tkl.iis.u-tokyo.ac.jp/~ynaga/cedar) which is written in C++ by Naoki Yoshinaga. `cedar-go` currently implements the `reduced` verion of cedar. 
This package is not thread safe if there is one goroutine doing insertions or deletions. 

## Install
```
go get github.com/adamzy/cedar-go
```

## Usage
```go
package main

import (
	"fmt"

	"github.com/adamzy/cedar-go"
)

func main() {
	// create a new cedar trie.
	trie := cedar.New()

	// a helper function to print the id-key-value triple given trie node id
	printIdKeyValue := func(id int) {
		// the key of node `id`.
		key, _ := trie.Key(id)
		// the value of node `id`.
		value, _ := trie.Value(id)
		fmt.Printf("%d\t%s:%v\n", id, key, value)
	}

	// Insert key-value pairs.
    // The order of insertion is not important.
	trie.Insert([]byte("How many"), 0)
	trie.Insert([]byte("How many loved"), 1)
	trie.Insert([]byte("How many loved your moments"), 2)
	trie.Insert([]byte("How many loved your moments of glad grace"), 3)
	trie.Insert([]byte("姑苏"), 4)
	trie.Insert([]byte("姑苏城外"), 5)
	trie.Insert([]byte("姑苏城外寒山寺"), 6)

	// Get the associated value of a key directly.
	value, _ := trie.Get([]byte("How many loved your moments of glad grace"))
	fmt.Println(value)

	// Or, jump to the node first,
	id, _ := trie.Jump([]byte("How many loved your moments"), 0)
	// then get the key and the value
	printIdKeyValue(id)

	fmt.Println("\nPrefixMatch\nid\tkey:value")
	for _, id := range trie.PrefixMatch([]byte("How many loved your moments of glad grace"), 0) {
		printIdKeyValue(id)
	}

	fmt.Println("\nPrefixPredict\nid\tkey:value")
	for _, id := range trie.PrefixPredict([]byte("姑苏"), 0) {
		printIdKeyValue(id)
	}
}
```
will produce
```
3
281	How many loved your moments:2

PrefixMatch
id	key:value
262	How many:0
268	How many loved:1
281	How many loved your moments:2
296	How many loved your moments of glad grace:3

PrefixPredict
id	key:value
303	姑苏:4
309	姑苏城外:5
318	姑苏城外寒山寺:6
```
