import "github.com/PuerkitoBio/goquery/v2"

doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
