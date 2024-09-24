package libUtils

import (
	"fmt"
	"regexp"
	"strings"
)

func cleanSql(sqlStr string) string {
	sqlStr = strings.Trim(sqlStr, " ")
	sqlStr = strings.Trim(sqlStr, ";")
	sqlStr = strings.Trim(sqlStr, fmt.Sprintln())

	return sqlStr
}

func isSpecialComment(line string, re *regexp.Regexp) bool {
	line = strings.ToLower(line)
	line = strings.Trim(line, " ")
	line = strings.Trim(line, fmt.Sprintln())

	return re.MatchString(line)
}

func IsStmtBegin(line string) bool {
	re, _ := regexp.Compile(`^\-{3}\s?statement\s?begin$`)
	return isSpecialComment(line, re)
}

func IsStmtEnd(line string) bool {
	re, _ := regexp.Compile(`^\-{3}\s?statement\s?end$`)
	return isSpecialComment(line, re)
}

func GetSQLStmts(sqlStr string) []string {
	stmts := []string{}
	lines := strings.Split(sqlStr, "\n")

	res := ""
	inStmt := false
	for _, line := range lines {
		line = line + fmt.Sprintln()

		cleanLine := cleanSql(line)
		if len(cleanLine) == 0 {
			continue
		}

		if IsStmtBegin(line) {
			inStmt = true
			continue
		} else if IsStmtEnd(line) {
			inStmt = false
			if len(res) > 0 {
				stmts = append(stmts, res)
			}
			res = ""
			continue
		} else if strings.HasPrefix(cleanLine, "--") {
			continue
		}

		if inStmt {
			res += line
		} else {
			for _, c := range line {
				s := string(c)

				if s == ";" {
					if len(res) > 0 {
						stmts = append(stmts, cleanSql(res))
					}
					res = ""
					continue
				}

				res += s
			}
		}
	}

	return stmts
}

// requires ";" followed by \n to recognize different statements
// func GetSQLStmts_DEPR(sqlStr string) []string {
// 	stmts := []string{}
// 	lines := ""
// 	line := ""
// 	inStatement := false
// 	for _, c := range sqlStr {
// 		s := string(c)

// 		if s == fmt.Sprintln() {
// 			if len(line) > 0 {
// 				if line == constants.SqlStmtBegin {
// 					inStatement = true
// 					line = ""
// 					continue
// 				} else if line == constants.SqlStmtEnd {
// 					inStatement = false
// 					cleanLines := cleanSql(lines)
// 					if len(cleanLines) > 0 {
// 						stmts = append(stmts, cleanSql(lines))
// 					}
// 					lines = ""
// 					line = ""
// 					continue
// 				} else if strings.HasPrefix(line, "--") {
// 					line = ""
// 					continue
// 				}

// 				lines += line + s

// 				if strings.HasSuffix(line, ";") && inStatement == false {
// 					stmts = append(stmts, cleanSql(lines))
// 					lines = ""
// 					line = ""
// 					continue
// 				}
// 			}

// 			line = ""
// 			continue
// 		}

// 		line += s
// 	}

// 	return stmts
// }
