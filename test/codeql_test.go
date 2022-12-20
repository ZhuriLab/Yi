package test

import (
	"Yi/pkg/logging"
	"Yi/pkg/runner"
	"Yi/pkg/utils"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"os"
	"testing"
)

/**
  @author: yhy
  @since: 2022/12/13
  @desc: //TODO
**/

func TestCodeql(t *testing.T) {
	runner.DirNames.ResDir = "./test/"
	runner.QLFiles = &runner.QLFile{
		GoQL: []string{"/Users/yhy/CodeQL/codeql/go/ql/src/Security/CWE-020"},
	}
	for fileName, res := range runner.Analyze("/Users/yhy/CodeQL/database/go/grafana/v8.2.6", "grafana", "go", nil) {
		results := jsoniter.Get([]byte(res), "runs", 0, "results")

		if results.Size() > 0 {
			maps := make(map[string]bool)
			for i := 0; i < results.Size(); i++ {
				location := "{"

				msg := "ruleId: " + results.Get(i).Get("ruleId").ToString() + "\t "
				if results.Get(i).Get("relatedLocations").Size() > 0 {
					for j := 0; j < results.Get(i).Get("relatedLocations").Size(); j++ {
						msg += "relatedLocations: " + results.Get(i).Get("relatedLocations").Get(j).Get("physicalLocation", "artifactLocation", "uri").ToString() + "\t startLine: " + results.Get(i).Get("relatedLocations").Get(j).Get("physicalLocation", "region", "startLine").ToString() + "\t | "

						line := results.Get(i).Get("relatedLocations").Get(j).Get("physicalLocation", "region", "startLine").ToString()
						location += fmt.Sprintf("\"%s#L%s\":\"%s\",", results.Get(i).Get("relatedLocations").Get(j).Get("physicalLocation", "artifactLocation", "uri").ToString(), line, line)
					}

				} else if results.Get(i).Get("locations").Size() > 0 {
					for j := 0; j < results.Get(i).Get("locations").Size(); j++ {
						msg += "locations: " + results.Get(i).Get("locations").Get(j).Get("physicalLocation", "artifactLocation", "uri").ToString() + "\t startLine: " + results.Get(i).Get("locations").Get(j).Get("physicalLocation", "region", "startLine").ToString() + "\t | "

						line := results.Get(i).Get("locations").Get(j).Get("physicalLocation", "region", "startLine").ToString()
						location += fmt.Sprintf("\"%s#L%s\":\"%s\",", results.Get(i).Get("locations").Get(j).Get("physicalLocation", "artifactLocation", "uri").ToString(), line, line)
					}

				}

				if _, ok := maps[location]; ok {
					continue
				}

				location += "}"

				logging.Logger.Infof("(%s) Found: %s", fileName, msg)
			}
		} else {
			err := os.Remove(fileName) //删除文件

			if err != nil {
				logging.Logger.Infof("file remove Error! %s\n", err)
			}
		}
	}
}

func TestRead(t *testing.T) {

	lines := utils.LoadFile("/Users/yhy/Desktop/2022-12-14/superedge/1671002613.json")

	var result string

	for _, i := range lines {
		result += i
	}

	results := jsoniter.Get([]byte(result), "runs", 0, "results")

	if results.Size() > 0 {
		maps := make(map[string]bool)

		for i := 0; i < results.Size(); i++ {
			location := "{"

			msg := "ruleId: " + results.Get(i).Get("ruleId").ToString() + "\t "
			if results.Get(i).Get("locations").Size() > 0 {
				for j := 0; j < results.Get(i).Get("locations").Size(); j++ {
					msg += "locations: " + results.Get(i).Get("locations").Get(j).Get("physicalLocation", "artifactLocation", "uri").ToString() + "\t startLine: " + results.Get(i).Get("locations").Get(j).Get("physicalLocation", "region", "startLine").ToString() + "\t | "

					line := results.Get(i).Get("locations").Get(j).Get("physicalLocation", "region", "startLine").ToString()
					location += fmt.Sprintf("\"%s#L%s\":\"%s\",", results.Get(i).Get("locations").Get(j).Get("physicalLocation", "artifactLocation", "uri").ToString(), line, line)
				}
			}
			if results.Get(i).Get("relatedLocations").Size() > 0 {
				for j := 0; j < results.Get(i).Get("relatedLocations").Size(); j++ {
					msg += "relatedLocations: " + results.Get(i).Get("relatedLocations").Get(j).Get("physicalLocation", "artifactLocation", "uri").ToString() + "\t startLine: " + results.Get(i).Get("relatedLocations").Get(j).Get("physicalLocation", "region", "startLine").ToString() + "\t | "

					line := results.Get(i).Get("relatedLocations").Get(j).Get("physicalLocation", "region", "startLine").ToString()
					location += fmt.Sprintf("\"%s#L%s\":\"%s\",", results.Get(i).Get("relatedLocations").Get(j).Get("physicalLocation", "artifactLocation", "uri").ToString(), line, line)
				}

			}

			codeFlows := results.Get(i).Get("codeFlows").ToString()

			fmt.Println("------------------")
			fmt.Println(codeFlows)
			location += "}"

			if _, ok := maps[location]; ok {
				continue
			}

			fmt.Println(location)

			logging.Logger.Infof("(%s) Found: %s", "test", msg)
		}
	}

}
