package parser

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"lacon-go-tiny-scrapy/crawler"
	"lacon-go-tiny-scrapy/logger"
	"strings"
	"sync"
	"time"
)

var (
	poolSize int

	workersGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "parser_worker_count",
		Help: "Number of active parsers",
	})

	parserSummary = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "parser_worker_summary",
		Help: "Performances per worker",
	},
	[]string{
		"website",
	})
)

func init() {
	config := newParserConfiguration()
	poolSize = config.getWorkerPoolSize()
}

func NewParser(waitingGroup *sync.WaitGroup, inputHTML <- chan crawler.HTMLPage) <- chan ParsedContent {
	waitingGroup.Add(1)
	outputContent := newParsedOutputChannel()

	innerWaitingGroup := sync.WaitGroup{}
	innerWaitingGroup.Add(poolSize)
	logger.INFO.Printf("Starting %d parsing workers\n", poolSize)
	workersGauge.Set(float64(poolSize))
	for i := 0; i < poolSize; i++ {
		go startParsingWorker(&innerWaitingGroup, inputHTML, outputContent)
	}
	go closeOutput(&innerWaitingGroup, waitingGroup, outputContent)

	return outputContent
}

func newParsedOutputChannel() chan ParsedContent {
	return make(chan ParsedContent)
}

func startParsingWorker(wg *sync.WaitGroup, inputHTML <-chan crawler.HTMLPage, output chan ParsedContent) {

	for htmlPage := range inputHTML {
		startTime := time.Now().UnixNano()

		website := htmlPage.GetWebsite()
		websiteFields := getFieldsByWebsite(website)
		content := htmlPage.GetContent()

		extractedContent := extractFieldsFromContent(websiteFields, content)
		extractedContent.setWebsite(website)
		extractedContent.setId(htmlPage.GetId())
		output <- extractedContent

		endTime := time.Now().UnixNano()
		parserSummary.WithLabelValues(website).Observe(float64(endTime - startTime))
	}

	workersGauge.Dec()
	wg.Done()
}

func extractFieldsFromContent(websiteFields map[string]rule, content []byte) ParsedContent {
	results := newParsedContent()
	for fieldName, fieldRules := range websiteFields {
		extractedValue := extractContentByFieldRules(content, fieldRules)

		results.setField(fieldName, extractedValue)
	}

	return results
}

func extractContentByFieldRules(content []byte, fieldRules rule) string {
	htmlPage, _ := goquery.NewDocumentFromReader(bytes.NewReader(content))

	var result string
	tag := fieldRules.getTag()
	htmlPage.Find(tag).Each(func(i int, selection *goquery.Selection) {
		res, found := findInSelection(selection, fieldRules)
		if found {
			result = res
			return
		}
	})
	return strings.TrimSpace(result)
}

func findInSelection(selection *goquery.Selection, fieldRules rule) (string, bool) {
	conditionAttribute := fieldRules.getConditionAttribute()
	val, found := selection.Attr(conditionAttribute)

	conditionValue := fieldRules.getConditionValue()
	extractAttribute := fieldRules.getExtractAttribute()
	if !found || val != conditionValue {
		return "", false
	} else if found && val == conditionValue && extractAttribute != "" {
		val, _ = selection.Attr(extractAttribute)
		return val, true
	} else {
		filterTags := fieldRules.getFilterTags()
		filterSelection(selection, filterTags)
		return selection.Text(), true
	}
}

func filterSelection(selection *goquery.Selection, tags []string) {
	for _, tag := range tags {
		selection.Find(tag).Each(func(i int, element *goquery.Selection) {
			element.Remove()
		})
	}
}

func closeOutput(innerWaitingGroup *sync.WaitGroup, waitingGroup *sync.WaitGroup, outputChannel chan ParsedContent) {
	logger.INFO.Println("Waiting to parse all html pages")
	innerWaitingGroup.Wait()
	logger.INFO.Println("Parsing done, closing channel")
	close(outputChannel)
	logger.INFO.Println("Parsing done, channel closed")
	waitingGroup.Done()
}
