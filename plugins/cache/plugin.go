package cache

import (
	kv "github.com/philippgille/gokv/redis"
	"go.uber.org/zap"
)

func New(log *zap.SugaredLogger, opts *kv.Options) (*Cache, error) {
	client, err := kv.NewClient(*opts)
	if err != nil {
		return nil, err
	}
	return &Cache{client: &client, options: opts, log: log}, nil
}

type Cache struct {
	client  *kv.Client
	options *kv.Options
	log     *zap.SugaredLogger
}

func (c *Cache) Set(k string, v interface{}) error {
	return c.client.Set(k, v)
}

func (c *Cache) Get(k string, v interface{}) (bool, error) {
	return c.client.Get(k, v)
}

func (c *Cache) Delete(k string) error {
	return c.client.Delete(k)
}

func (c *Cache) Fetch(k string, v interface{}, f func() (interface{}, error)) (bool, error) {
	ok, err := c.client.Get(k, v)
	// there was an error
	if err != nil {
		return ok, err
	}
	// the item was found
	if ok {
		//c.log.Infof("cache: hit: %s", k)
		return ok, nil
	}

	// get the value and set it
	v, err = f()
	if err != nil {
		return false, err
	}
	//c.log.Infof("cache: miss: %s", k)
	return false, c.client.Set(k, v)
}

func (c *Cache) Close() error {
	return c.client.Close()
}
