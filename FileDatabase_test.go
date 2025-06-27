package filefinder

import (
	"fmt"
	"regexp"
	"testing"
)

func TestDirCheck(t *testing.T) {
	println(DirCheck("a/b/", "D:\\a\\b\\", true))
}

func TestNewFileDB(t *testing.T) {
	db, err := NewFileDB("/Users/restr0/Projects/GolangProjects")
	if err != nil {
		println(err)
		return
	}
	//for _, f := range db.Files {
	//	fmt.Printf("%s (%s)\n", f.Name, f.AbsDirPath)
	//}
	one, err := db.SearchOne(&SearchRule{
		DirRules:        nil,
		FileNameRegexps: []*regexp.Regexp{regexp.MustCompile(".*\\.go$")},
	})
	if err != nil {
		fmt.Println(err)
	}
	for regexStr, paths := range one {
		fmt.Printf("Regex: %s\n", regexStr)
		for i, path := range paths {
			fmt.Printf("  %d. %s\n", i, path)
		}
	}

}
