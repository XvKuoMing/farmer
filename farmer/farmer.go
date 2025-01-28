package farmer

// import (
// 	"context"
// 	"strings"
// 	"sync"
// 	"time"

// 	"github.com/XvKuoMing/farmer/internals"
// 	"github.com/go-rod/rod"
// 	"github.com/go-rod/rod/lib/proto"
// 	"github.com/redis/go-redis/v9"
// )

// type SpiderConfig struct {
// 	Domain          string // domain to parse
// 	ConcurencyLimit int    // limit of parallel pages
// 	Stealth         bool
// }

// type job struct {
// 	el       string                   // if el exists, call callback
// 	callback func(e *rod.Element) any //function that user defines
// }

// type Spider struct {
// 	engine        *rod.Browser
// 	pool          *rod.Pool[rod.Page]
// 	userAgent     *proto.NetworkSetUserAgentOverride
// 	options       *SpiderConfig
// 	storage       internals.SpiderStorage // optimized way to fastly check if url has been visited
// 	logger        *internals.Logger
// 	jobs          []job
// 	panicCounter  *func(r any, page *rod.Page)
// 	sucessCounter *func(res any, page *rod.Page)
// }

// // } // old way to wait gotouines https://stackoverflow.com/questions/18207772/how-to-wait-for-all-goroutines-to-finish-without-using-time-sleep

// // ------------------------------------------------init-------------------------------------

// func NewSpider(domain string, concurencyLimit int) *Spider {
// 	// init with default, any custom changes are done using method recievers
// 	spider := Spider{
// 		engine: rod.New(),
// 		options: &SpiderConfig{
// 			Domain:          domain,
// 			ConcurencyLimit: concurencyLimit,
// 		},
// 		storage: internals.GetinMemmorySpiderStorage(),
// 		logger:  internals.GetLogger(nil),
// 		jobs:    []job{},
// 		//panicCounter: func(page *rod.Page) {},
// 	}
// 	return &spider
// }

// // -------------------------------------------------config------------------------------------

// func (spider *Spider) SetRandomUserAgent(region string) {
// 	spider.userAgent = internals.SetRandomUserAgent(region)
// }

// func (spider *Spider) SetSpecificUserAgent(userAgent string, acceptLanguage string) {
// 	spider.userAgent = &proto.NetworkSetUserAgentOverride{
// 		UserAgent:      userAgent,
// 		AcceptLanguage: acceptLanguage,
// 	}
// }

// func (spider *Spider) SetStorage(ss internals.SpiderStorage) {
// 	// custom storage
// 	spider.storage = ss
// }

// func (spider *Spider) SetRedisStorage(opts *redis.Options) {
// 	// ready to use redis storage
// 	spider.storage = internals.GetRedisSpiderStorage(
// 		context.Background(),
// 		opts,
// 	)
// }

// func (spider *Spider) SetLogger(from string) {
// 	spider.logger = internals.GetLogger(&from)
// }

// func (spider *Spider) launch() {
// 	spider.engine.MustConnect()
// 	_pool := rod.NewPagePool(spider.options.ConcurencyLimit)
// 	spider.pool = &_pool
// }

// func (spider *Spider) close() {
// 	spider.engine.Close()
// }

// // ------------------------------------------------logging------------------------------------------

// func (spider *Spider) Info(v ...any) {
// 	spider.logger.Info.Println(v...)
// }

// func (spider *Spider) Warn(v ...any) {
// 	spider.logger.Warn.Println(v...)
// }

// func (spider *Spider) Err(v ...any) {
// 	spider.logger.Err.Println(v...)
// }

// func (spider *Spider) FatalErr(v ...any) {
// 	spider.logger.Err.Fatalln(v...)
// }

// // ----------------------------------------------------handlers-------------------------------------

// func (spider *Spider) Retrieve(el string, callback func(e *rod.Element) any) {
// 	// runs callback if and only if el realy exists on a page
// 	job := job{
// 		el,
// 		callback,
// 	}
// 	spider.jobs = append(spider.jobs, job)
// }

// func (spider *Spider) OnSuccess(do func(res any, page *rod.Page)) {
// 	// runs do every successfull retrieval, if there are several -> runs for each
// 	spider.sucessCounter = &do
// }

// func (spider *Spider) OnPanic(do func(r any, page *rod.Page)) {
// 	// runs do when spider encounters a panic
// 	spider.panicCounter = &do
// }

// func (spider *Spider) retrieveAll(page *rod.Page) {
// 	// runs all jobs in spider.jobs on given page
// 	for _, job := range spider.jobs {
// 		if page.MustHas(job.el) {
// 			// spider.Info("Found " + job.el + " in page " + page.MustInfo().URL)
// 			el := page.Timeout(time.Minute).MustElement(job.el).Timeout(3 * time.Minute)
// 			if el != nil {
// 				res := job.callback(el)
// 				if spider.sucessCounter != nil {
// 					(*spider.sucessCounter)(res, page)
// 				}
// 			}
// 		}
// 	}
// }

// // func (spider *Spider) OnPanic(do func(page *rod.Page)) {
// // 	// runs do when spider encounters a panic / TODO /
// // 	spider.panicCounter = do
// // }

// // ----------------------------------------------------------internal for main launcher-------------------------------------------

// func (spider *Spider) getPage() *rod.Page {
// 	instance := func() (*rod.Page, error) {
// 		return spider.engine.Page(proto.TargetCreateTarget{})
// 	}
// 	page, _ := spider.pool.Get(instance)
// 	return page
// }

// func (spider *Spider) putPage(page *rod.Page) {
// 	spider.pool.Put(page)
// }

// func (spider *Spider) updateQueue(curPage *rod.Page) {
// 	links := "a[href]"
// 	if curPage.MustHas(links) {
// 		els := curPage.Timeout(3 * time.Minute).MustElements(links)
// 		for _, el := range els {
// 			_nextUrl := *el.MustAttribute("href")
// 			if !internals.IsUrlDomain(_nextUrl) {
// 				continue
// 			}

// 			nextUrl := internals.AbsUrl(spider.options.Domain, _nextUrl)
// 			nextUrl = strings.TrimSuffix(nextUrl, "/")
// 			if !(spider.storage.IsVisited(nextUrl) || spider.storage.InQueue(nextUrl)) {
// 				spider.storage.SetQueue(nextUrl)
// 			}

// 		}
// 	}
// }

// func (spider *Spider) visit(urls []string, depth int) {
// 	// visits all pages in spider.pages with depth = depth

// 	run := func(url string) {

// 		var page *rod.Page

// 		defer spider.putPage(page)
// 		if spider.panicCounter != nil {
// 			defer func() {
// 				if r := recover(); r != nil {
// 					(*spider.panicCounter)(r, page)
// 				}
// 			}()
// 		}

// 		page.Timeout(1 * time.Minute).MustWindowFullscreen().MustWaitStable()
// 		isVisited := spider.storage.IsVisited(url)
// 		if !(depth == 0 && isVisited) {
// 			page.SetUserAgent(spider.userAgent)
// 			spider.Info("Navigating to " + url)
// 			page.Navigate(url)
// 			page.Timeout(3 * time.Minute).MustWindowFullscreen().MustWaitStable()
// 		}

// 		if isVisited {
// 			spider.Info(url + " has been already visited, skiping")
// 		} else {
// 			spider.Info("Retrieving info from  " + url)
// 			spider.retrieveAll(page)
// 			spider.storage.SetVisited(url) // в дальнейшем нужно сделать так， чтобы setVisited не вызывалась
// 		}

// 		if depth != 0 {
// 			spider.updateQueue(page) // collection new urls for parsing
// 		}
// 	}

// 	wg := sync.WaitGroup{}
// 	for i, url := range urls {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			run(url)
// 		}()

// 		if (i+1)%spider.options.ConcurencyLimit == 0 {
// 			wg.Wait()
// 		}
// 	}
// 	wg.Wait()

// 	if depth != 0 {
// 		queue := spider.storage.GetQueue(spider.options.Domain)
// 		if len(queue) > 0 {
// 			spider.visit(queue, depth-1)
// 		}
// 	}

// }

// // --------------------------------------------------------------main launcher-----------------------------

// func (spider *Spider) Visit(urls []string, depth int) {
// 	spider.launch()
// 	defer spider.close()

// 	spider.visit(urls, depth)

// 	spider.Info("Done")
// }
