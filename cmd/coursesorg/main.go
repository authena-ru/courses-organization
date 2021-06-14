package main

import "github.com/authena-ru/courses-organization/internal/runner"

const configsDir = "configs"

func main() {
	runner.Start(configsDir)
}
