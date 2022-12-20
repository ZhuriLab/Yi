package test

import (
	"fmt"
	"path/filepath"
	"testing"
)

/**
  @author: yhy
  @since: 2022/12/13
  @desc: //TODO
**/

func TestGithub(t *testing.T) {
	//req, _ := http.NewRequest("GET", "https://api.github.com/repos/airshipit/treasuremap/languages", nil)
	//req.Header.Set("Accept", "application/vnd.github.v3.text-match+json")
	//
	//req.Header.Set("User-Agent", uarand.GetRandom())
	//req.Close = true
	//
	//resp, err := utils.Httpx("").Do(req)
	//
	//if err != nil {
	//	logging.Logger.Errorln("GetLanguage client.Do(req) err:", err)
	//
	//}
	//defer resp.Body.Close()
	//body, _ := ioutil.ReadAll(resp.Body)
	//results := jsoniter.Get(body).Keys()
	//
	//fmt.Println(results)

	fmt.Println(filepath.Join("../", filepath.Clean("/"+".../../../../../../../etc/passwd")))
}
