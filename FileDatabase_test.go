package filefinder

import (
	"fmt"
	"testing"
)

func TestDirCheck(t *testing.T) {
	println(DirCheck("a/b/", "D:\\a\\b\\", true))
}

func TestNewFileDB(t *testing.T) {
	db, err := NewFileDB("/Users/xuziyi/Desktop/TestPy")
	if err != nil {
		println(err)
		return
	}
	//for _, f := range db.Files {
	//	fmt.Printf("%s (%s)\n", f.Name, f.AbsDirPath)
	//}
	one, err := db.SearchOne(SearchRule{
		DirRules:        nil,
		FileNameRegexps: []string{".+\\.py"},
	})
	for i, name := range one {
		fmt.Printf("%d. %s\n", i, name)
	}

}
