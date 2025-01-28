# Farmer: Dynamic Web Crawling Library With No Headache

Farmer is a powerful and dynamic web crawling library built on the robust [go-rod](https://github.com/go-rod/rod) framework. It introduces a high-level API for creating spiders dynamically and offers features like stealth mode, proxy rotation, and Redis-based storage for parsed data. Farmer is designed for developers who want an easy, scalable, and flexible solution for web crawling with minimal setup.

---

## Key Features

- **Dynamic Spider Generation**: Automatically generate spiders from high-level requirements.
- **Stealth Mode**: Avoid detection by anti-crawling mechanisms.
- **Redis Storage**: Built-in support for Redis to store parsed data efficiently.
- **Proxy Rotation**: Seamlessly rotate proxies to enhance anonymity (feature in development).
- **Randomization**: Use randomized user agents and other browser configurations for diverse scraping sessions.
- **Docker Support**: Coming soon for seamless containerized deployments.
- **Microservice Ready**: Planned support for deploying spiders as microservices.
- **Easy-to-Use API**: High-level API for setting up spiders and defining callbacks.
- **Concurrency Control**: Manage multiple pages with configurable concurrency limits.
- **Customizable Logging**: Enable detailed logs for debugging and monitoring.

---

## Getting Started

### Installation

```bash
go get github.com/XvKuoMing/farmer
```

### Prerequisites

- [Go](https://golang.org/) 1.18 or higher
- Redis (optional but recommended for large-scale crawls)

---

## Example Usage

### Auto-Generated Spider

Generate a spider dynamically based on requirements and run it:

```go
package main

import (
    "fmt"
    "os"

    "github.com/XvKuoMing/farmer/farmer"
)

func main() {
    // .env file is suggested, this is just for the example
    _ = os.Setenv("OPENAI_API_KEY", "your-openai-api-key")
    _ = os.Setenv("OPENAI_BASE_URL", "https://api.openai.com/v1/") // Optional for default URL

    spider := farmer.Text2Spider(
        "product price, product name",
        "https://example.com",
        "https://example.com/sample-product")
    
    fmt.Println("Executing auto-generated spider...") // As for now this will simply yield items in console
    // in the future those items will be accessible with reference to the page and custom callbacks
    spider.Exec()
}
```

### Custom Spider with Redis Storage

Define a custom spider with specific elements and callbacks:

```go
package main

import (
    "strconv"

    "github.com/XvKuoMing/farmer/farmer"
    "github.com/go-rod/rod"
    "github.com/redis/go-redis/v9"
)

func main() {
    redisOptions := &redis.Options{
        Addr: "localhost:6379",
        DB:   0,
    }

    spiderOptions := &farmer.SpiderExecutableOptions{
        Domain:          "https://example.com",
        ConcurencyLimit: 3,
        Depth:           3,
        Redis:           redisOptions,
        LoggerOn:        true,
        LoggerDir:       "./logs/",
        El:              "div.product-item",
        Callback: func(e *rod.Element) any {
            price, _ := strconv.Atoi(e.MustElement(".price").MustText())
            fmt.Printf("Found product: %s, Price: %d\n", e.MustElement(".title").MustText(), price)
            return nil
        },
    }

    spider := farmer.NewSpiderExec(spiderOptions)
    spider.Exec()
}
```

---

## Configuration

### Environment Variables

- `OPENAI_API_KEY`: API key for OpenAI (used for spider generation).
- `OPENAI_BASE_URL`: Base URL for OpenAI API (default: `https://api.openai.com/v1/`).

### Spider Options

| Option              | Description                                           |
|---------------------|-------------------------------------------------------|
| `Domain`            | Domain to crawl                                       |
| `ConcurencyLimit`   | Maximum number of concurrent pages                    |
| `Depth`             | Depth of crawling                                     |
| `Redis`             | Redis configuration for storing crawled data          |
| `LoggerOn`          | Enable or disable logging                             |
| `LoggerDir`         | Directory for storing logs                            |
| `El`                | CSS selector for the target element                  |
| `Callback`          | Function to handle the extracted elements            |

---

## Roadmap

- [ ] Full integration of all `go-rod` features into the high-level API
- [ ] Support for proxy rotation and stealth mode
- [ ] Improved handling of auto-generated spiders
- [ ] Docker images for easy deployment
- [ ] Microservice-ready architecture
- [ ] Real-time monitoring dashboard for spiders

---

## Contributions

Contributions are welcome! Please open an issue or submit a pull request.

---

## License

Farmer is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.

--- 