package showdoc

import (
	"errors"
	"fmt"
	"github.com/levigross/grequests"
)

var Host = "https://www.showdoc.com.cn" // ShowDoc地址，默认 "https://www.showdoc.com.cn"

// 认证凭证。登录showdoc，进入具体项目后，点击右上角的”项目设置”-“开放API”便可看到
var (
	ApiKey   = ""
	ApiToken = ""
)

// UpdateByApi ShowDoc提供的开放API接口，使用api_key+api_token进行认证。
// 参考官方文档 https://www.showdoc.com.cn/page/102098
//
// @param catName 可选参数。当页面文档处于目录下时，请传递目录名。当目录名不存在时，showdoc会自动创建此目录。需要创建多层目录的时候请用斜杆隔开，例如 “一层/二层/三层”。
// @param pageTitle 页面标题。请保证其唯一。（或者，当页面处于目录下时，请保证页面标题在该目录下唯一）。当页面标题不存在时，showdoc将会创建此页面。当页面标题存在时，将用page_content更新其内容
// @param sNumber 可选，页面序号。默认是99。数字越小，该页面越靠前
func UpdateByApi(catName, pageTitle, sNumber string, content string) error {
	url := fmt.Sprintf("%s/server/?s=/api/item/updateByApi", Host)
	data := map[string]string{
		"api_key":      ApiKey,
		"api_token":    ApiToken,
		"cat_name":     catName,
		"page_title":   pageTitle,
		"page_content": content, // 页面内容
		"s_number":     sNumber,
	}
	result := ErrResult{}
	if err := post(url, data, &result); err != nil {
		return err
	}
	return result.Error()
}

func post(url string, data map[string]string, result interface{}) error {
	resp, err := grequests.Post(url, &grequests.RequestOptions{
		Data: data,
	})
	if err != nil {
		return err
	}
	return resp.JSON(result)
}

type ErrResult struct {
	ErrorCode    int    `json:"error_code"` // 返回 0 表示成功
	ErrorMessage string `json:"error_message"`
}

func (p *ErrResult) Error() error {
	if p.ErrorCode != 0 {
		return errors.New(p.ErrorMessage)
	}
	return nil
}
