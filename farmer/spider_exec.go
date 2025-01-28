package farmer

import (
	"strings"

	"github.com/XvKuoMing/farmer/internals"
	"github.com/go-rod/rod"
	"github.com/redis/go-redis/v9"
)

type SpiderExecutable struct {
	spider *Spider
	exec   func()
}

type SpiderExecutableOptions struct {
	Domain          string
	StartUrl        string
	ConcurencyLimit int
	Depth           int
	Redis           *redis.Options
	LoggerOn        bool
	LoggerDir       string
	El              string
	Callback        func(e *rod.Element) any
}

func NewSpiderExec(options *SpiderExecutableOptions) *SpiderExecutable {

	options.Domain = strings.TrimSuffix(options.Domain, "/")

	if options.StartUrl == "" {
		options.StartUrl = options.Domain
	}

	spider := NewSpider(options.Domain, options.ConcurencyLimit)

	if options.LoggerOn {
		spider.SetLogger(options.LoggerDir + internals.NameUrl(options.Domain) + "/")

		spider.OnPanic(func(r any, page *rod.Page) {
			spider.Err("run into panic from "+page.MustInfo().URL+"\nerror message:", r)
		})

		spider.OnSuccess(func(res any, page *rod.Page) {
			spider.Info("success from "+page.MustInfo().URL+"\n", res)
		})
	}

	if options.Redis != nil {
		spider.SetRedisStorage(options.Redis)
	}

	if strings.HasSuffix(options.Domain, ".com") {
		spider.SetRandomUserAgent("US")
	} else {
		spider.SetRandomUserAgent("Russia")
	}

	spider.Retrieve(options.El, options.Callback)

	executable := SpiderExecutable{
		spider: spider,
		exec:   func() { spider.Visit([]string{options.StartUrl}, options.Depth) },
	}
	return &executable
}

func (spiderExecutable *SpiderExecutable) AddRetrieval(el string, callback func(e *rod.Element) any) {
	spiderExecutable.spider.Retrieve(el, callback)
}

func (spiderExecutable *SpiderExecutable) AddOnSuccess(do func(res any, page *rod.Page)) {
	spiderExecutable.spider.OnSuccess(do)
}

func (spiderExecutable *SpiderExecutable) AddOnPanic(do func(r any, page *rod.Page)) {
	spiderExecutable.spider.OnPanic(do)
}

func (spiderExecutable *SpiderExecutable) Exec() {
	spiderExecutable.exec()
}
