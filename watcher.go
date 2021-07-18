package micro_consul

import (
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/micro/go-micro/config/source"
)

type watcher struct {
	c *consul
	kv *api.KV
	exit chan bool
}

func newWatcher(c *consul) (source.Watcher, error) {
	config := api.DefaultConfig()
	config.Address = c.url
	config.Token = c.token
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	kv := client.KV()

	return &watcher{
		c: c,
		kv: kv,
		exit: make(chan bool),
	}, nil
}

func (w *watcher) Next() (*source.ChangeSet, error) {
	select {
	case <-w.exit:
		return nil, source.ErrWatcherStopped
	default:
	}

	options := api.QueryOptions{
		WaitIndex: w.c.lastWaitIndex,
	}
	pair, meta, err := w.kv.Get(w.c.key, &options)
	if err != nil {
		return nil, err
	}

	if pair == nil {
		return nil, source.ErrWatcherStopped
	}

	w.c.lastWaitIndex = meta.LastIndex

	cs := &source.ChangeSet{
		Format:    "json",
		Source:    w.c.String(),
		Timestamp: time.Now(),
		Data:      pair.Value,
	}
	cs.Checksum = cs.Sum()

	return cs, nil
}

func (w *watcher) Stop() error {
	return nil
}