//go:build !windows && !linux

package main

func handleServiceMode(_ []string) (bool, error) {
	return false, nil
}
