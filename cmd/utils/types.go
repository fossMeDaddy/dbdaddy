package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	VersionTagRegexp = regexp.MustCompile(`^[0-9A-Za-z-\.]+$`)
)

// DO NOT initialize this struct directly
type Version struct {
	Major int
	Minor int
	Patch int
	Tag   string
}

func (v *Version) String() string {
	vStr := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if len(v.Tag) > 0 {
		vStr += fmt.Sprintf("-%s", v.Tag)
	}

	return vStr
}

func NewVersion(vStr string) (Version, error) {
	v := Version{}

	tagSplit := strings.SplitN(vStr, "-", 2)
	if len(tagSplit) == 2 {
		tag := tagSplit[1]

		if !VersionTagRegexp.MatchString(tag) {
			return v, fmt.Errorf("'%s' must be of the format ^[0-9A-Za-z-]+$", tag)
		}

		v.Tag = tag
	} else if len(tagSplit) == 0 {
		return v, fmt.Errorf("received empty version string")
	}

	vSplit := strings.Split(tagSplit[0], ".")
	if len(vSplit) != 3 {
		return v, fmt.Errorf("invalid dot-separated segments, need: 3 (read semver spec for once)")
	}

	if major, err := strconv.Atoi(vSplit[0]); err != nil {
		fmt.Println("major version parse error")
		return v, err
	} else {
		v.Major = major
	}

	if minor, err := strconv.Atoi(vSplit[1]); err != nil {
		fmt.Println("minor version parse error")
		return v, err
	} else {
		v.Minor = minor
	}

	if patch, err := strconv.Atoi(vSplit[2]); err != nil {
		fmt.Println("patch version parse error")
		return v, err
	} else {
		v.Patch = patch
	}

	return v, nil
}
