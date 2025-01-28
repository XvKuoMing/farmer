package farmer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const ONE_PASS int = 1
const ONE_PAGE int = 0

var client *openai.Client

func Text2Spider(items string, domain string, target string) *SpiderExecutable {
	// openai uses prompt to parse all necessary elements from target page
	page := CollectElements(domain, target)
	if client == nil {
		InitOpenai()
	}
	args := Generate(items, page)

	args = strings.TrimPrefix(args, "```json")
	args = strings.TrimSuffix(args, "```")
	fmt.Println(args)

	var selectors map[string]interface{}
	err := json.Unmarshal([]byte(args), &selectors)
	if err != nil {
		panic(fmt.Sprintf("failed to create spider executable from %s", args))
	}

	var spider *SpiderExecutable

	for item, selector := range selectors {

		el := selector.(string)
		callback := func(e *rod.Element) any {
			fmt.Println("on page: ", e.Page().MustInfo().Title)
			fmt.Println(item, ",", e.MustElement(el).Timeout(30*time.Second).MustText())

			return 0
		}

		if spider == nil {
			exeOptions := &SpiderExecutableOptions{
				Domain:          domain,
				StartUrl:        target,
				ConcurencyLimit: 1,
				Depth:           1,
				LoggerOn:        true,
				LoggerDir:       "./logs/",
				El:              el,
				Callback:        callback,
			}

			spider = NewSpiderExec(exeOptions) // then simply go .Exec()

		} else {
			spider.AddRetrieval(el, callback)
		}

	}

	return spider

}

func CollectElements(domain string, target string) string {
	// init simple spider to parse html and feed openai, then rechecks
	var result string

	feeder := NewSpider(domain, ONE_PASS)
	feeder.Retrieve(
		"body",
		func(e *rod.Element) any {

			result = e.Page().MustHTML()
			return 0
		},
	)
	feeder.Visit([]string{target}, ONE_PAGE)
	return result
}

// Helper function to get an environment variable or a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func InitOpenai() {
	client = openai.NewClient(
		option.WithAPIKey(os.Getenv("OPENAI_API_KEY")), // defaults to os.LookupEnv("OPENAI_API_KEY")
		option.WithBaseURL(getEnv("OPENAI_BASE_URL", "https://api.openai.com/v1/")),
	)
}

func Generate(items string, html_as_str string) string {
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(
				fmt.Sprintf(`
			You are a backend server that auto translates human language request into a json of html elements needed to be parsed from given html.
			Here is what user wants to parse from html page: %s
			Yout task is to find appropriate html elements from given target html-page.
			<output format>
			{
			"<the item's name user's requested>": "html-element as css-selector or XPath"
			}
			</output format>
			Ensure you output the correct json.
			If there is no selector or path for enlisted items, then return null.
			Always output json with keys of all requested items.
			IMPORTANT: you must make it clear what are you targeting: content of the element or attribute value.
			When targeting attribute value of the element -> ALWAYS select it using selector
			When targeting content of the element only -> omit attributes at the end of the selector.
			This helps to avoid ambiguity for user.
			<example of targeting content of element>
				<user's request> 
				title of the text
				</user's request>

				<given html>
					<div class="text-holder">
						<h1 id=1>TITLE</h1>
						<p>some text</p>
					</div>
				</given html>
				<output>
				{
				"title": div.text-holder > h1
				}
				</output>
			</example of targeting content of element>
			<example of targeting element attribute>
				<user's request> 
				id of the title of the text
				</user's request>

				<given html>
					<div class="text-holder">
						<h1 id=1>TITLE</h1>
						<p>some text</p>
					</div>
				</given html>
				<output>
				{
				"title": div.text-holder > h1[id]
				}
				</output>
			</example of targeting element attribute>
			IMPORTANT: you MUST capture the most generic selectors for items, such as other similar products could be matched with your generated selectors
			
			Emerge in your role.

			JSON:
			`, items)),
			openai.UserMessage(html_as_str),
		}),
		Model:       openai.F(openai.ChatModelGPT4oMini),
		Temperature: openai.F(0.7),
	})
	if err != nil {
		panic(err.Error())
	}

	return chatCompletion.Choices[0].Message.Content
}
