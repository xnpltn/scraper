package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

type Book struct {
	Name   string `json:"name"`
	Price  string `json:"price"`
	Image  string `json:"image"`
	Rating uint   `json:"rating"`
}

type Books struct {
	Books []Book `json:"books"`
}

func main() {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var booksNode []*cdp.Node

	err := chromedp.Run(ctx,
		chromedp.Navigate("https://books.toscrape.com/"),
		chromedp.Nodes(".product_pod", &booksNode, chromedp.ByQueryAll),
	)
	if err != nil {
		log.Fatal(err)
	}
	books := new(Books)
	for _, node := range booksNode {
		book := new(Book)
		var rating string
		err := chromedp.Run(ctx,
			chromedp.AttributeValue(".image_container > a > img", "src", &book.Image, new(bool), chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.AttributeValue(".product_pod > p", "class", &rating, new(bool), chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(".product_pod > div.product_price >p", &book.Price, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(".product_pod > h3 >a", &book.Name, chromedp.ByQuery, chromedp.FromNode(node)),
		)
		if err != nil {
			log.Println(err)
			continue
		}
		switch true {
		case strings.Contains(rating, "One"):
			book.Rating = 1
		case strings.Contains(rating, "Two"):
			book.Rating = 2
		case strings.Contains(rating, "Three"):
			book.Rating = 3
		case strings.Contains(rating, "Four"):
			book.Rating = 4
		case strings.Contains(rating, "Five"):
			book.Rating = 5
		default:
			book.Rating = 0

		}
		books.Books = append(books.Books, *book)
	}
	data, err := json.Marshal(books)
	if err != nil {
		log.Printf("error marshalling json: %q\n", err.Error())
	}
	jsonFile, err := os.Create("out.json")
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			log.Printf("error creating file: %q\n", err.Error())
		}
	}

	_, err = jsonFile.Write(data)
	if err != nil {
		log.Printf("error saving to file: %q\n", err.Error())
	}

	fmt.Println("data saved to: ", jsonFile.Name())
}
