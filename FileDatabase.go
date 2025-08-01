package filefinder

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type File struct {
	Name       string
	AbsDirPath string //必须以文件分隔符结尾
}

type SearchRule struct {
	RuleName        string
	DirRules        []string         //目录的约束条件，条件之间的关系是或。如果需要表达且应当令起一个搜索配置
	FileNameRegexps []*regexp.Regexp //目标文件名的正则表达式
}

type FileDB struct {
	FileIndex map[string][]int //map[文件名][]文件索引
	Files     []*File
	IsWindows bool
}

func (fdb FileDB) Search(rules []*SearchRule, funcs ...func(regex string, path string) string) map[string]map[string][]string {
	result := map[string]map[string][]string{}
	for _, rule := range rules {
		r, _ := fdb.SearchOne(rule, funcs...)
		result[rule.RuleName] = r
	}

	// 结果去重
	for ruleName, regexResults := range result {
		for regexStr, paths := range regexResults {
			clearResult := []string{}
			midResult := map[string]bool{}
			for _, path := range paths {
				if !midResult[path] {
					midResult[path] = true
					clearResult = append(clearResult, path)
				}
			}
			result[ruleName][regexStr] = clearResult
		}
	}
	return result
}

func (fdb FileDB) SearchOne(rule *SearchRule, funcs ...func(regex string, path string) string) (map[string][]string, error) {
	results := make(map[string][]string)
	targetIndexes := make(map[string][]int)

	//如果有多个正则匹配到同样的文件名，结果将出现重复路径
	for _, exp := range rule.FileNameRegexps {
		for fileName, indexes := range fdb.FileIndex {
			if exp.Match([]byte(fileName)) {
				targetIndexes[exp.String()] = append(targetIndexes[exp.String()], indexes...)
			}
		}
	}

	for expStr, indexes := range targetIndexes {
		var result []string
		for _, index := range indexes {
			f := fdb.Files[index]

			if len(rule.DirRules) > 0 {
				for _, r := range rule.DirRules {
					if DirCheck(r, f.AbsDirPath, fdb.IsWindows) {
						result = append(result, f.AbsDirPath+f.Name)
						break
					}
				}
			} else {
				path := f.AbsDirPath + f.Name
				for _, fn := range funcs {
					path = fn(expStr, path)
				}
				if path != "" {
					result = append(result, path)
				}
			}
		}
		results[expStr] = result
	}

	return results, nil

}

// Append 添加文件
func (fdb *FileDB) Append(path string) {
	dir, name := filepath.Split(path)
	fdb.Files = append(fdb.Files, &File{
		Name:       name,
		AbsDirPath: dir,
	})
	fdb.FileIndex[name] = append(fdb.FileIndex[name], len(fdb.Files)-1)
}

// ChangeOSType 修改FileDB路径解析格式
func (fdb *FileDB) ChangeOSType(isWindows bool) {
	fdb.IsWindows = isWindows
}

// NewFileDB root置空可以获得一个空FileDB
func NewFileDB(root string) (*FileDB, error) {
	fdb := &FileDB{
		FileIndex: map[string][]int{},
	}
	if os.PathSeparator == '\\' {
		fdb.IsWindows = true
	}
	if root == "" {
		return fdb, nil
	}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		fdb.Append(path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return fdb, nil
}

func DirCheck(DirRule, AbsFilePath string, isWindows bool) bool {
	var yes = true
	skip := false
	var ruleParts, pathParts []string
	if isWindows {
		pathParts = strings.Split(AbsFilePath, "\\")
		pathParts[0] = ""
	} else {
		pathParts = strings.Split(AbsFilePath, "/")
	}

	ruleParts = strings.Split(DirRule, "/")

	if len(ruleParts) <= 0 {
		return false
	}

	if ruleParts[0] != "" { //规则不是根目录
		skip = true
	}

	if ruleParts[len(ruleParts)-1] != "" { //结尾可以不直接跟文件名
		ruleParts = append(ruleParts, "...", "")
	}

	ruleIndex := 0

	for _, part := range pathParts {
		//fmt.Println(part, " => ", ruleParts[ruleIndex], "SKIP <", skip, ">")
		if !skip {
			if ruleParts[ruleIndex] == "..." {
				skip = true
				ruleIndex += 1
			} else {
				if ruleParts[ruleIndex] != "*" && part != ruleParts[ruleIndex] { //不符合
					yes = false
				} else { //符合
					ruleIndex += 1
				}
			}
		} else {
			if part == ruleParts[ruleIndex] { //符合
				skip = false
				ruleIndex += 1
			}
		}
	}

	return yes && !skip
}
