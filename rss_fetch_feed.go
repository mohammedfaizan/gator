package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title	string  `xml:"title"`
		Link 	string	`xml:"link"`
		Description string `xml:"description"`
		Item 		[]RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title string `xml:"title"`
	Link string `xml:"link"`
	Description string `xml:"description"`
	PubDate	string `xml:"pubDate"`
}


func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error)  {
	req, err := http.NewRequest("GET", feedURL, nil)
	if err != nil {
		fmt.Println("Error sending request: ", err)
		return &RSSFeed{}, err
	}

	req.Header.Add("User-Agent","gator")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request: ", err)
		return &RSSFeed{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the response", err)
		return &RSSFeed{}, err
	}

	var rssFeed RSSFeed
	err = xml.Unmarshal(body, &rssFeed)
	if err != nil {
		fmt.Println("error unmarshalling the body", err)
		return &RSSFeed{}, err
	}

	defer resp.Body.Close()

	//unescape title update
	channelTitle := html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Title = channelTitle

	//unescape description update
	channelDes := html.UnescapeString(rssFeed.Channel.Description)
	rssFeed.Channel.Description = channelDes
	
	
	// updating all item's title and description
	for i, item := range rssFeed.Channel.Item {

		//item title update
		itemTitle := html.UnescapeString(item.Title)
		rssFeed.Channel.Item[i].Title = itemTitle

		//item description update
		itemDes := html.UnescapeString(item.Description)
		rssFeed.Channel.Item[i].Description = itemDes
	}


	return &rssFeed, nil
}