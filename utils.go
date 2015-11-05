package main

import "os"

func StringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func CurrentDir() string {
	path, _ := os.Getwd()

	return path
}
