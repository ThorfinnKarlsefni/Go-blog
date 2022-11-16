package requests

import (
	"goblog/app/models/article"

	"github.com/thedevsaddam/govalidator"
)

func ValidateArticleForm(data article.Article) map[string][]string {
	rules := govalidator.MapData{
		"title": []string{"required", "min_cn:3", "max_cn:40"},
		"body":  []string{"required", "min_cn:30"},
	}

	messages := govalidator.MapData{
		"title": []string{
			"required:标题为必填项",
			"min_cn:标题长度需大于3",
			"max_cn:标题长度需小于40",
		},
		"body": []string{
			"requried:标题为必填项",
			"min_cn:长度需要大于 10",
		},
	}

	opts := govalidator.Options{
		Data:          &data,
		Rules:         rules,
		TagIdentifier: "valid",
		Messages:      messages,
	}

	return govalidator.New(opts).ValidateStruct()
}
