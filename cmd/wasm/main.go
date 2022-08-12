//go:build js

// +build,js,!unittest

package main

func main() {
	Run()
	select {}
}
