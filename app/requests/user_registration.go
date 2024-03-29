package requests

import (
	"github.com/thedevsaddam/govalidator"

	"goblog/app/models/user"
)

func ValidateRegistrationForm(data user.User) map[string][]string {

	rules := govalidator.MapData{
		"name":             []string{"required", "alpha_num", "between:3,20", "not_exists:users,name"},
		"email":            []string{"required", "min:4", "max:30", "email", "not_exists:users,email"},
		"password":         []string{"required", "min:6"},
		"password_confirm": []string{"required"},
	}

	messages := govalidator.MapData{
		"name": []string{
			"required:用户名必须填写",
			"alpha_num:格式错误 只允许数字和英文",
			"between:用户名长度需要 3~20 之间",
		},
		"email": []string{
			"required:Email 为必填项",
			"min:Email 长度需要大于40",
			"max:Email 长度小于30",
			"email:Email 格式不正确 请提供有效邮箱",
		},

		"password": []string{
			"required:密码为必填项",
			"min:长度需大于6",
		},
		"password_confirm": []string{
			"required:确认密码框为必填项",
		},
	}

	opts := govalidator.Options{
		Data:          &data,
		Rules:         rules,
		TagIdentifier: "valid", // 模型中的 Struct 标签标识符
		Messages:      messages,
	}

	errs := govalidator.New(opts).ValidateStruct()

	if data.Password != data.PasswordConfirm {
		errs["password_confirm"] = append(errs["password_confirm"], "两次输入密码不匹配!")
	}

	return errs
}
