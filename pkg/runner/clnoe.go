package runner

import (
	"Yi/pkg/logging"
	"bytes"
	"os/exec"
)

/**
  @author: yhy
  @since: 2022/10/27
  @desc: //TODO
**/

func GitClone(gurl string, name string) error {
	cmd := exec.Command("git", "clone", gurl, "--depth=1", DirNames.GithubDir+name)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误
	err := cmd.Run()
	_, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if err != nil {
		logging.Logger.Errorf("GitClone(%s) cmd.Run() failed with %s --  %s\n", gurl, err, errStr)

		return err
	}
	return nil
}
