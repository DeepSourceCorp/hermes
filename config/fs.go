package config

import "os"

type readFileFn func(name string) ([]byte, error)

type statFn func(name string) (os.FileInfo, error)
