package requests

import (
	"goblog/app/models/article"

	"github.com/thedevsaddam/govalidator"
)

func ValidateArticleForm(data article.Article) map[string][]string {
	rules := govalidator.MapData{
		"title": []string{"required", "min:3", "max:40"},
		"body":  []string{"required", "min:30"},
	}

	messages := govalidator.MapData{
		"title": []string{
			"required:标题为必填项",
			"min:标题长度需大于3",
			"max:标题长度需小于40",
		},
		"body": []string{
			"requried:标题为必填项",
			"min:长度需要大于 10",
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
