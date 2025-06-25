package domain

import (
	"fmt"
	"path"
	"strings"
)

type Path string

func (p Path) String() string {
	return string(p)
}

func (p Path) IsRoot() bool {
	return p == "/" || p == ""
}

func (p Path) IsValid() bool {
	str := string(p)
	if str == "" {
		return false
	}
	if !strings.HasPrefix(str, "/") {
		return false
	}
	return path.Clean(str) == str
}

func (p Path) GetParent() Path {
	if p.IsRoot() {
		return Path("/")
	}
	return Path(path.Dir(string(p)))
}

func (p Path) GetName() string {
	if p.IsRoot() {
		return ""
	}
	return path.Base(string(p))
}

func (p Path) Join(name string) Path {
	if p.IsRoot() {
		return Path("/" + name)
	}
	return Path(path.Join(string(p), name))
}

func (p Path) Clean() Path {
	return Path(path.Clean(string(p)))
}

func NewPath(str string) (Path, error) {
	if str == "" {
		str = "/"
	}

	if !strings.HasPrefix(str, "/") {
		str = "/" + str
	}

	cleaned := path.Clean(str)
	p := Path(cleaned)

	if !p.IsValid() {
		return "", fmt.Errorf("invalid path %v", str)
	}

	return p, nil
}
