package internals

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type SpiderStorage interface {
	IsVisited(url string) bool       // gets url string and checks whether it was visited
	SetVisited(url string) error     // gets url string and set it as it was visited
	SetQueue(url string) error       // gets url and sets to queue for visiting
	InQueue(url string) bool         // gets url string and checks whether it is in queue for visiting
	GetQueue(domain string) []string // get the whole list of urls in queue for visiting
}

const QUEUE_KEY string = "QUEUE>>" // constant prefix used for distinguising visited url from queued url

// --------------------------in memmory -------------

type inMemmorySpiderStorage struct {
	// as name suggests -> purely RAM memory
	urls   map[string]struct{} // a set of urls
	locker sync.Mutex          // an instance of mutex
}

var singletonInMemmory *inMemmorySpiderStorage // global variable that is used to ensure singleton

func GetinMemmorySpiderStorage() *inMemmorySpiderStorage {
	if singletonInMemmory == nil {
		singletonInMemmory = &inMemmorySpiderStorage{
			urls:   make(map[string]struct{}),
			locker: sync.Mutex{},
		}
	}
	return singletonInMemmory
}

func (ss *inMemmorySpiderStorage) IsVisited(url string) bool {
	// returns true if url not in storage, else returns false meaning : do not visit it since it is already visited
	ss.locker.Lock()
	_, ok := ss.urls[url]
	ss.locker.Unlock()
	return ok
}

func (ss *inMemmorySpiderStorage) SetVisited(url string) error {
	ss.locker.Lock()
	ss.urls[url] = struct{}{}
	ss.locker.Unlock()
	return nil
}

func (ss *inMemmorySpiderStorage) InQueue(url string) bool {
	return ss.IsVisited(QUEUE_KEY + url)
}

func (ss *inMemmorySpiderStorage) SetQueue(url string) error {
	return ss.SetVisited(QUEUE_KEY + url)
}

func (ss *inMemmorySpiderStorage) GetQueue(domain string) []string {
	var queue []string
	for url := range ss.urls {
		if strings.HasPrefix(url, QUEUE_KEY+domain) {
			queue = append(queue, strings.TrimPrefix(url, QUEUE_KEY))
		}
	}
	return queue
}

// ---------------redis--------------------

type redisSpiderStorage struct {
	ctx   context.Context
	redis *redis.Client
}

var singletonRedis *redisSpiderStorage

func GetRedisSpiderStorage(ctx context.Context, opts *redis.Options) *redisSpiderStorage {
	if singletonRedis == nil {
		singletonRedis = &redisSpiderStorage{
			ctx:   ctx,
			redis: redis.NewClient(opts),
		}
	}
	return singletonRedis
}

func errAsResult(err error) bool {
	if err == redis.Nil {
		return false
	} else if err == nil {
		return true
	} else {
		panic(err)
	}
}

func (rss *redisSpiderStorage) SetQueue(url string) error {
	return rss.redis.Set(rss.ctx, QUEUE_KEY+url, 0, 1*time.Hour).Err()
}

func (rss *redisSpiderStorage) InQueue(url string) bool {
	_, err := rss.redis.Get(rss.ctx, QUEUE_KEY+url).Result()
	return errAsResult(err)
}

func (rss *redisSpiderStorage) IsVisited(url string) bool {
	_, err := rss.redis.Get(rss.ctx, url).Result()
	return errAsResult(err)
}

func (rss *redisSpiderStorage) SetVisited(url string) error {
	return rss.redis.Set(rss.ctx, url, 0, (24*30)*time.Hour).Err()
}

func (rss *redisSpiderStorage) GetQueue(domain string) []string {
	var queue []string
	iter := rss.redis.Scan(rss.ctx, 0, QUEUE_KEY+domain+"*", 0).Iterator()
	for iter.Next(rss.ctx) {
		queue = append(queue, strings.TrimPrefix(iter.Val(), QUEUE_KEY))
	}
	return queue
}
