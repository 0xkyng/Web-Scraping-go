package main

import (
	"encoding/csv"
	"fmt"
	// "go/scanner"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

var urls = []string{
	"https://www.amazon.com/iPhone-Pro-256GB-Sierra-Blue/dp/B0BGYFDQJX",
	"https://www.amazon.com/AMD-Ryzen-5600X-12-Thread-Processor/dp/B08166SLDF",
	"https://www.amazon.com/SAMSUNG-Border-Less-TUV-Certified-Intelligent-LS32A700NWNXZA/dp/B08V6MNW9P/",
}

type Product struct {
	name  string
	price string
}

func main() {
	// Measure time taken
	// To extarct the data
	t := time.Now()
	var products []Product

	// Wait group
	var wg sync.WaitGroup
	// Declare a channel
	ch := make(chan Product)

	// Add the lengthof the 
	// Urls to the wait group
	wg.Add(len(urls))

	for _, url := range urls {
		// Make the channel
		// Concurrent
		go scrape(url, ch)
	}

	// Receive data
	for range urls {
		// Create concurrent
		// Annonymous function
		go func() {
			// Defer waitgroup
			defer wg.Done()
			product := <-ch
			products = append(products, product)
		}()
	}

	// Wait for it 
	// To complete
	wg.Wait()

	// Create a file
	file, err := os.Create("data.csv")
	// Handle error
	if err != nil {
		log.Fatalln(err)
	}

	// Use defer to make sure
	// the file gets closed
	defer file.Close()

	// Create a csv writer
	writer := csv.NewWriter(file)
	// Flush writer
	defer writer.Flush()

	// add a header
	writer.Write([]string{"name", "price"})

	// Loop over the products 
	// And write them to csv
	for _, product := range products {
		writer.Write([]string{
			product.name,
			product.price,
		})
	}

	// Close channel
	close(ch)

	// Check time taken
	// And print it
	elapsed := time.Since(t).Seconds()
	fmt.Printf("%.2f", elapsed)
}

// Restrict the channel
func scrape(url string, ch chan<- Product) {
	var product Product
	// Use a scraping library
	// By initializing a ne collector
	c := colly.NewCollector()

	c.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL)
	})

	c.OnHTML("#ppd", func(e *colly.HTMLElement) {
		name := e.ChildText("#centerCol #productTitle")
		price := e.ChildText("#rightCol #corePrice_feature_div")
		product = Product{name, price}
	})

	c.Visit(url)
	// Send product
	// To the channel
	ch <- product
}
