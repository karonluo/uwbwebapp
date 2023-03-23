package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"uwbwebapp/conf"
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/cache"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"
)

func main() {
	conf.LoadWebConfig("../conf/WebConfig.json")
	dao.InitDatabase()
	cache.InitRedisDatabase()
	type TMP struct {
		Url         string `json:"url"`
		Method      string `json:"method"`
		DisplayName string `json:"display_name"`
	}
	var tmp []TMP
	var bcontent []byte
	bcontent, _ = os.ReadFile("../sysfuncpages.json")
	fmt.Println(string(bcontent))
	json.Unmarshal(bcontent, &tmp)
	var pageId string
	var err error
	err = biz.ClearSysFuncPages()
	if err == nil {
		for idx, t := range tmp {
			// fmt.Println(t.Url)
			var funcPage entities.SysFuncPage
			funcPage.DisplayName = t.DisplayName
			funcPage.OrderKey = idx
			funcPage.URLAddress = t.Url
			funcPage.ParentID = "top"
			funcPage.URLType = "INTERFACE"
			funcPage.URLMethod = strings.ToUpper(t.Method)
			pageId, err = biz.CreateSysFuncPage(&funcPage)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(pageId)
			}

		}
	}

}
