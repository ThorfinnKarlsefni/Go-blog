package requests

import (
	"errors"
	"fmt"
	"goblog/pkg/model"
	"strings"

	"github.com/thedevsaddam/govalidator"
)

func init() {
	govalidator.AddCustomRule("not_exists", func(field, rule, message string, value interface{}) error {
		rng := strings.Split(strings.TrimPrefix(rule, "not_exists:"), ",")

		tableName := rng[0]
		dbField := rng[1]
		val := value.(string)

		var count int64

		model.DB.Table(tableName).Where(dbField+"= ?", val).Count(&count)

		if count != 0 {
			if message != "" {
				return errors.New(message)
			}

			return fmt.Errorf("%v 已被占用", val)
		}

		return nil
	})
}