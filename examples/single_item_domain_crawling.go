package examples

import (
	"strconv"

	"github.com/XvKuoMing/farmer/farmer"
	"github.com/go-rod/rod"
	"github.com/redis/go-redis/v9"
)

const MAX_THREADS = 3
const MAX_DEPTH = 3
const LOGGER_DIR string = "./logs/"

func getRedisOptions() *redis.Options {
	return &redis.Options{
		Addr:     "localhost:6379",
		Password: "password", // no password set
		DB:       0,          // use default DB
		Protocol: 3,
	}
}

type Product struct {
	Url         string `json:"url"`
	Shop        string `json:"shop"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	PriceUnit   string `json:"price_unit"`
	ImgPath     string `json:"img_path"`
	Description string `json:"description,omitempty"`
}

func ExeSpider() *farmer.SpiderExecutable {

	exeOptions := &farmer.SpiderExecutableOptions{
		Domain:          "https://domain.com",
		ConcurencyLimit: MAX_THREADS,
		Depth:           MAX_DEPTH,
		Redis:           getRedisOptions(),
		LoggerOn:        true,
		LoggerDir:       LOGGER_DIR,
		El:              "div.catalog-element",
		Callback: func(e *rod.Element) any {
			price, _ := strconv.Atoi(e.MustElementX(`//span[@class="prices__actual"]/span[1]/text()[1]`).MustText())
			product := Product{
				Url:         e.Page().MustInfo().URL,
				Shop:        "spar",
				Name:        e.MustElement("h1.catalog-element__title").MustText(),
				Price:       price,
				PriceUnit:   *e.MustElement("span.prices__actual > span.prices__unit[data-unit]").MustAttribute("data-unit"),
				ImgPath:     *e.MustElement("a.catalog-element__picture-area").MustAttribute("href"),
				Description: e.MustElement("div.catalog-element__desc > p").MustText(),
			}
			// you can here insert product into your storage
			return product
		},
	}

	return farmer.NewSpiderExec(exeOptions) // then simply go .Exec()
}
