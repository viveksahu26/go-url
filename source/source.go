package source

import "regexp"

func IsGit(in string) bool {
	return regexp.MustCompile("^(http|https)://").MatchString(in)
}
