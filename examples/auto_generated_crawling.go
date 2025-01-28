package examples

import (
	"fmt"
	"os"

	"github.com/XvKuoMing/farmer/farmer"
)

func Run_auto_generated_spider() {
	_ = os.Setenv("OPENAI_API_KEY", "sk-")
	_ = os.Setenv("OPENAI_BASE_URL", "https://api.../v1/") // omit if openai

	spider := farmer.Text2Spider(
		"product price,",
		"https://someshop.com",
		"https://someshop.com/goods/good-example-of-target-item")
	fmt.Println("WARNING: as for now this feature simply prints data; storing parsed data will be available soon!")
	spider.Exec()
}
