package micro_consul

import (
	"context"
	"go-micro.dev/v4/config/source"
)

type consulConfig struct{
	key   string
	url   string
	token string
}

func WithConfig(url string, key string, token string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		config := consulConfig{
			key: key,
			url: url,
			token: token,
		}
		o.Context = context.WithValue(o.Context, consulConfig{}, config)
	}
}