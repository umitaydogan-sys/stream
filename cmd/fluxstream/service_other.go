//go:build !windows

package main

func handleServiceMode(_ []string) (bool, error) {
	return false, nil
}
