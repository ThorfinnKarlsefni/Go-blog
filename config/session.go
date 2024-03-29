package config

import "goblog/pkg/config"

func init() {
	config.Add("session", config.StrMap{
		"default": config.Env("SESSION_DRIVER", "cookie"),

		"session_name": config.Env("SESSION_NAME", "goblog-session"),
	})
}
