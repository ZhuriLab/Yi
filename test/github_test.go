package test

import (
	"Yi/pkg/runner"
	"Yi/pkg/utils"
	"fmt"
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
	runner.Option.Session = utils.NewSession("")
	fmt.Println(runner.GetLanguage("https://api.github.com/repos/agragregra/OptimizedHTML-4", "https://github.com/agragregra/OptimizedHTML-4"))

	asd := make(map[string]string)

	delete(asd, "11")

	fmt.Println("1 ", asd)

	asd["11"] = "12"

	delete(asd, "22")

	fmt.Println(asd)
	delete(asd, "11")

	fmt.Println(asd)
}
