package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	client *redis.Client
}

func NewRedis() (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:        os.Getenv("REDIS"),
		DB:          0,
		DialTimeout: 100 * time.Millisecond,
		ReadTimeout: 100 * time.Millisecond,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) GetName(ctx context.Context, nconst string) (Person, error) {
	cmd := c.client.Get(ctx, nconst)

	cmdb, err := cmd.Bytes()
	if err != nil {
		return Person{}, err
	}

	b := bytes.NewReader(cmdb)

	var res Person

	if err := gob.NewDecoder(b).Decode(&res); err != nil {
		return Person{}, err
	}

	return res, nil
}

func (c *Client) SetName(ctx context.Context, n Person) error {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(n); err != nil {
		return err
	}

	return c.client.Set(ctx, n.Id, b.Bytes(), 25*time.Second).Err()
}
