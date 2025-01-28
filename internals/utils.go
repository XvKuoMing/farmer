package internals

import (
	"math/rand/v2"
	"strings"

	"github.com/go-rod/rod/lib/proto"
)

// -------------------------------------url processing----------------------------
func IsUrlDomain(url string) bool {
	if strings.HasPrefix(url, "/") {
		return true
	} else {
		return false
	}
}

func AbsUrl(domain string, url string) string {
	return domain + url
}

func NameUrl(domain string) string {
	rawDomain := strings.Split(domain, "://")[1]
	rawDomain = strings.TrimPrefix(rawDomain, "www.")
	return strings.Split(rawDomain, ".")[0] // domain name between www. or https?:// AND first .
}

// -------------------------------------user agent-----------------------------------

var userAgents []string = []string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:130.0) Gecko/20100101 Firefox/130.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36",
}

var acceptedLanguages map[string]string = map[string]string{
	"Russia": "ru-RU,ru;q=0.9",
	"US":     "en-US,en;q=0.5",
}

func SetRandomUserAgent(region string) *proto.NetworkSetUserAgentOverride {
	// random choice
	choice := rand.IntN(len(userAgents))
	return &proto.NetworkSetUserAgentOverride{
		UserAgent:      userAgents[choice],
		AcceptLanguage: acceptedLanguages[region],
	}

}
