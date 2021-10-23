package micro_consul

import (
	"go-micro.dev/v4/config/source"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

type consul struct {
	url           string
	token         string
	key           string
	lastWaitIndex uint64
	opts          source.Options
}

func (c *consul) Read() (*source.ChangeSet, error) {
	config := api.DefaultConfig()
	config.Address = c.url
	config.Token = c.token
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	kv := client.KV()

	pair, meta, err := kv.Get(c.key, nil)
	if err != nil {
		return nil, err
	}

	if pair == nil {
		return nil, errors.New("Consul key not found")
	}

	c.lastWaitIndex = meta.LastIndex

	cs := &source.ChangeSet{
		Format:    "json",
		Source:    c.String(),
		Timestamp: time.Now(),
		Data:      pair.Value,
	}
	cs.Checksum = cs.Sum()

	return cs, nil
}

func (c *consul) String() string {
	return "consul"
}

func (c *consul) Watch() (source.Watcher, error) {
	return newWatcher(c)
}

func (c *consul) Write(cs *source.ChangeSet) error {
	return nil
}

func NewSource(opts ...source.Option) source.Source {
	options := source.NewOptions(opts...)
	config := consulConfig{}
	f, ok := options.Context.Value(consulConfig{}).(consulConfig)
	if ok {
		config = f
	}
	return &consul{opts: options, key: config.key, url: config.url, token: config.token}
}
