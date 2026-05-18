package community

import (
	"context"
	"fmt"
	"strconv"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

const (
	likedKeyTTL  = 24 * time.Hour
	countKeyTTL  = 24 * time.Hour
)

// likedKey: user:liked:{userId}:{postId}
func likedKey(userID, postID int64) string {
	return fmt.Sprintf("user:liked:%d:%d", userID, postID)
}

// countKey: post:like:count:{postId}
func countKey(postID int64) string {
	return fmt.Sprintf("post:like:count:%d", postID)
}

// Lua script: 아토믹 좋아요 토글
// KEYS[1] = likedKey, KEYS[2] = countKey
// ARGV[1] = likedTtl(sec), ARGV[2] = countTtl(sec)
var toggleLuaScript = goredis.NewScript(`
local likedKey = KEYS[1]
local countKey = KEYS[2]
local likedTtl = tonumber(ARGV[1])
local countTtl  = tonumber(ARGV[2])

local current = redis.call('GET', likedKey)
local liked
local likeCount

if current == '1' then
    redis.call('SET', likedKey, '0', 'EX', likedTtl)
    likeCount = redis.call('DECR', countKey)
    liked = 0
else
    redis.call('SET', likedKey, '1', 'EX', likedTtl)
    likeCount = redis.call('INCR', countKey)
    liked = 1
end

redis.call('EXPIRE', countKey, countTtl)
return {liked, likeCount}
`)

type likeRedis struct {
	rdb *goredis.Client
}

func newLikeRedis(rdb *goredis.Client) *likeRedis {
	return &likeRedis{rdb: rdb}
}

func (r *likeRedis) initCache(ctx context.Context, postID, userID, dbCount int64, dbLiked bool) {
	likedVal := "0"
	if dbLiked {
		likedVal = "1"
	}
	r.rdb.Set(ctx, likedKey(userID, postID), likedVal, likedKeyTTL)
	r.rdb.Set(ctx, countKey(postID), dbCount, countKeyTTL)
}

func (r *likeRedis) toggle(ctx context.Context, postID, userID int64) (liked bool, count int64, err error) {
	ttlSec := int64(likedKeyTTL.Seconds())
	countTtlSec := int64(countKeyTTL.Seconds())

	keys := []string{likedKey(userID, postID), countKey(postID)}
	argv := []interface{}{ttlSec, countTtlSec}

	res, err := toggleLuaScript.Run(ctx, r.rdb, keys, argv...).Slice()
	if err != nil {
		return false, 0, err
	}

	likedInt, _ := res[0].(int64)
	countInt, _ := res[1].(int64)
	return likedInt == 1, countInt, nil
}

func (r *likeRedis) getLikedStatus(ctx context.Context, userID, postID int64) (bool, bool) {
	val, err := r.rdb.Get(ctx, likedKey(userID, postID)).Result()
	if err != nil {
		return false, false
	}
	return val == "1", true
}

func (r *likeRedis) getLikeCount(ctx context.Context, postID int64) (int64, bool) {
	val, err := r.rdb.Get(ctx, countKey(postID)).Result()
	if err != nil {
		return 0, false
	}
	count, err := strconv.ParseInt(val, 10, 64)
	return count, err == nil
}
