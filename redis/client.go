package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	connTimeout = 1 * time.Second
	rwTimeout   = 1 * time.Second
	poolSize    = 100
	poolTimeout = 30 * time.Second
)

type Option struct {
	Addr string
	Pwd  string
	Db   int32
}

type Client struct {
	rdb *redis.Client
}

func New(opt *Option) (*Client, error) {
	// new client
	client := &Client{}
	client.rdb = redis.NewClient(&redis.Options{
		Addr:         opt.Addr,
		Password:     opt.Pwd,
		DialTimeout:  connTimeout,
		ReadTimeout:  rwTimeout,
		WriteTimeout: rwTimeout,
		PoolSize:     poolSize,
		PoolTimeout:  poolTimeout,
		DB:           int(opt.Db),
	})
	// check conn
	err := client.rdb.Ping(context.Background()).Err()
	return client, err
}

func MustNew(opt *Option) *Client {
	c, err := New(opt)
	if err != nil {
		panic(any(err))
	}
	return c
}

// String
func (r *Client) Get(ctx context.Context, key string) string {
	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		fmt.Println("redis get fail:", key, err, val)
	}
	return val
}

func (r *Client) GetInt(ctx context.Context, key string) (int, bool, error) {
	val, err := r.rdb.Get(ctx, key).Int()
	switch {
	case err == redis.Nil:
		return 0, false, nil
	case err != nil:
		return 0, false, err
	default:
		return val, true, nil
	}
}

func (r *Client) GetString(ctx context.Context, key string) (string, bool, error) {
	val, err := r.rdb.Get(ctx, key).Result()
	switch {
	case err == redis.Nil:
		return "", false, nil
	case err != nil:
		return "", false, err
	default:
		return val, true, nil
	}
}

func (r *Client) Set(ctx context.Context, key string, val interface{}) bool {
	err := r.rdb.Set(ctx, key, val, 0).Err()
	if err != nil {
		fmt.Println("redis set fail:", key, val, err)
		return false
	}
	return true
}

func (r *Client) Exists(ctx context.Context, key string) int64 {
	val, err := r.rdb.Exists(ctx, key).Result()
	if err != nil {
		fmt.Println("redis exist fail:", key, err)
	}
	return val
}

func (r *Client) SetEX(ctx context.Context, key string, val interface{}, expire int64) bool {
	err := r.rdb.Set(ctx, key, val, time.Duration(expire)*time.Second).Err()
	if err != nil {
		fmt.Println("redis setex fail:", key, val, expire, err)
		return false
	}
	return true
}

func (r *Client) SetEXNotChangeTTL(ctx context.Context, key string, val interface{}) (expire bool, err error) {
	ttl, err := r.rdb.TTL(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if ttl < 0 { // 键值已过期
		return true, nil
	}
	err = r.rdb.Set(ctx, key, val, ttl).Err()
	if err != nil {
		return false, err
	}
	return false, nil
}

func (r *Client) SetNX(ctx context.Context, key string, val interface{}, expire int64) (bool, error) {
	ok, err := r.rdb.SetNX(ctx, key, val, time.Duration(expire)*time.Second).Result()
	if err != nil {
		fmt.Println("redis setnx fail:", key, val, expire, err)
		return false, err
	}
	return ok, nil
}

func (r *Client) Expire(ctx context.Context, key string, expire int64) bool {
	err := r.rdb.Expire(ctx, key, time.Duration(expire)*time.Second).Err()
	if err != nil {
		fmt.Println("redis expire fail:", key, expire, err)
		return false
	}
	return true
}

func (r *Client) Del(ctx context.Context, keys ...string) bool {
	err := r.rdb.Del(ctx, keys...).Err()
	if err != nil {
		fmt.Println("redis del fail:", keys, err)
		return false
	}
	return true
}

func (r *Client) Incr(ctx context.Context, key string) int64 {
	val, err := r.rdb.Incr(ctx, key).Result()
	if err != nil {
		fmt.Println("redis incr fail:", key, err)
	}
	return val
}

func (r *Client) IncrBy(ctx context.Context, key string, value int64) int64 {
	val, err := r.rdb.IncrBy(ctx, key, value).Result()
	if err != nil {
		fmt.Println("redis incrby fail:", key, value, err)
	}
	return val
}

func (r *Client) Decr(ctx context.Context, key string) int64 {
	val, err := r.rdb.Decr(ctx, key).Result()
	if err != nil {
		fmt.Println("redis decr fail:", key, err)
	}
	return val
}

func (r *Client) MGet(ctx context.Context, keys ...string) map[string]string {
	m := make(map[string]string, len(keys))
	vals, err := r.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		fmt.Println("redis mget fail:", keys, err)
	} else {
		for k, val := range keys {
			m[val] = fmt.Sprintf("%v", vals[k])
		}
	}
	return m
}

func (r *Client) MSet(ctx context.Context, m map[string]interface{}) bool {
	err := r.rdb.MSet(ctx, m).Err()
	if err != nil {
		fmt.Println("mset key fail:", m, err)
		return false
	}
	return true
}

func (r *Client) Do(ctx context.Context, args ...interface{}) interface{} {
	val, err := r.rdb.Do(ctx, args...).Result()
	if err != nil {
		fmt.Println("redis do fail:", args, err)
	}
	return val
}

// List
func (r *Client) LPop(ctx context.Context, key string) string {
	val, err := r.rdb.LPop(ctx, key).Result()
	if err != nil {
		fmt.Println("redis lpop fail:", key, err)
	}
	return val
}

func (r *Client) LPush(ctx context.Context, key string, vals ...interface{}) int64 {
	val, err := r.rdb.LPush(ctx, key, vals...).Result()
	if err != nil {
		fmt.Println("redis lpush fail:", key, vals, err)
	}
	return val
}

func (r *Client) RPop(ctx context.Context, key string) string {
	val, err := r.rdb.RPop(ctx, key).Result()
	if err != nil && err != redis.Nil {
		fmt.Println("redis rpop fail:", key, err)
	}
	return val
}

func (r *Client) RPush(ctx context.Context, key string, vals ...interface{}) int64 {
	val, err := r.rdb.RPush(ctx, key, vals...).Result()
	if err != nil {
		fmt.Println("redis rpush fail:", key, vals, err)
	}
	return val
}

func (r *Client) LLen(ctx context.Context, key string) int64 {
	val, err := r.rdb.LLen(ctx, key).Result()
	if err != nil {
		fmt.Println("redis llen fail:", key, err)
	}
	return val
}

func (r *Client) LRange(ctx context.Context, key string, start, stop int64) []string {
	vals, err := r.rdb.LRange(ctx, key, start, stop).Result()
	if err != nil {
		fmt.Println("redis llen fail:", key, err)
	}
	return vals
}

// Hash
func (r *Client) HGet(ctx context.Context, key, field string) (string, error) {
	val, err := r.rdb.HGet(ctx, key, field).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *Client) HGetInt64(ctx context.Context, key, field string) (int64, error) {
	val, err := r.rdb.HGet(ctx, key, field).Int64()
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (r *Client) HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error) {
	vals, err := r.rdb.HMGet(ctx, key, fields...).Result()
	if err != nil {
		return nil, err
	}
	return vals, nil
}

func (r *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	val, err := r.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (r *Client) HIncrBy(ctx context.Context, key, field string, num int64) (int64, error) {
	val, err := r.rdb.HIncrBy(ctx, key, field, num).Result()
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (r *Client) HSet(ctx context.Context, key string, vals ...interface{}) error {
	err := r.rdb.HSet(ctx, key, vals...).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Client) HMSet(ctx context.Context, key string, vals ...interface{}) (bool, error) {
	ok, err := r.rdb.HMSet(ctx, key, vals...).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (r *Client) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	num, err := r.rdb.HDel(ctx, key, fields...).Result()
	if err != nil {
		return 0, err
	}
	return num, nil
}

func (r *Client) HExists(ctx context.Context, key, field string) (bool, error) {
	val, err := r.rdb.HExists(ctx, key, field).Result()
	if err != nil {
		return false, err
	}
	return val, nil
}

func (r *Client) HKeys(ctx context.Context, key string) ([]string, error) {
	val, err := r.rdb.HKeys(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return val, nil
}

// Set
func (r *Client) SRandMemberN(ctx context.Context, key string, count int64) map[string]string {
	m := make(map[string]string, count)
	vals, err := r.rdb.SRandMemberN(ctx, key, count).Result()
	if err != nil {
		fmt.Println("redis SRandMemberN fail:", key, err)
	} else {
		for k, val := range vals {
			m[val] = fmt.Sprintf("%v", vals[k])
		}
	}
	return m
}

// zset
func (r *Client) ZAdd(ctx context.Context, key string, members ...*redis.Z) error {
	err := r.rdb.ZAdd(ctx, key, members...).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Client) ZRemRangeByScore(ctx context.Context, key, min, max string) error {
	err := r.rdb.ZRemRangeByScore(ctx, key, min, max).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Client) ZRevRangeByScore(ctx context.Context, key string, max, min string) ([]string, error) {
	members, err := r.rdb.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: min,
		Max: max,
	}).Result()
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (r *Client) SetBit(ctx context.Context, key string, offset int64) bool {
	return r.rdb.SetBit(ctx, key, offset, 1).Val() > 0
}

func (r *Client) GetBit(ctx context.Context, key string, offset int64) (int64, error) {
	bit, err := r.rdb.GetBit(ctx, key, offset).Result()
	if err != nil {
		return 0, err
	}

	return bit, nil
}

func (r *Client) BitCount(ctx context.Context, key string, start, end int64) int64 {
	return r.rdb.BitCount(ctx, key, &redis.BitCount{Start: start, End: end}).Val()
}

func (r *Client) BitField(ctx context.Context, key string, args ...interface{}) ([]int64, error) {
	list, err := r.rdb.BitField(ctx, key, args...).Result()
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (r *Client) BitPos(ctx context.Context, key string, bit int64, pos ...int64) (int64, error) {
	idx, err := r.rdb.BitPos(ctx, key, bit, pos...).Result()
	if err != nil {
		return 0, err
	}

	return idx, nil
}

func (r *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := r.rdb.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return ttl, nil
}

// TODO BitOp
