package cache

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
)

type redisStore struct {
	pool *redis.Pool
}

// NewRedis 创建 Redis 缓存连接池，连接失败会在启动时的健康检查中暴露。
func NewRedis(addr string, password string) (Store, error) {
	return &redisStore{pool: &redis.Pool{
		MaxIdle:     8,
		IdleTimeout: 5 * time.Minute,
		Dial: func() (redis.Conn, error) {
			options := []redis.DialOption{redis.DialConnectTimeout(3 * time.Second)}
			if password != "" {
				options = append(options, redis.DialPassword(password))
			}
			return redis.Dial("tcp", addr, options...)
		},
	}}, nil
}

func (s *redisStore) Get(ctx context.Context, key string) ([]byte, bool, error) {
	conn, err := s.pool.GetContext(ctx)
	if err != nil {
		return nil, false, err
	}
	defer conn.Close()
	value, err := redis.Bytes(conn.Do("GET", key))
	if err == redis.ErrNil {
		return nil, false, nil
	}
	return value, err == nil, err
}

func (s *redisStore) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	conn, err := s.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	if ttl > 0 {
		_, err = conn.Do("SET", key, value, "EX", int(ttl.Seconds()))
	} else {
		_, err = conn.Do("SET", key, value)
	}
	return err
}

func (s *redisStore) Delete(ctx context.Context, key string) error {
	conn, err := s.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Do("DEL", key)
	return err
}

func (s *redisStore) Health(ctx context.Context) string {
	conn, err := s.pool.GetContext(ctx)
	if err != nil {
		return "down"
	}
	defer conn.Close()
	if _, err := conn.Do("PING"); err != nil {
		return "down"
	}
	return "ok"
}

func (s *redisStore) Close() error { return s.pool.Close() }
