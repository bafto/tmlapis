package main

import (
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/net/html"
)

const url string = "http://javid.ddns.net/tModLoader/tools/ranksbysteamid.php?steamid64="

type ModStats struct {
	RankTotal          int
	DisplayName        string
	DownloadsTotal     int
	DownloadsYesterday int
}

func GetHtml(steamId string) (*html.Node, error) {
	resp, err := http.Get(url + steamId)
	if err != nil {
		return nil, err
	}
	return html.Parse(resp.Body)
}

func GetNodesByTag(doc *html.Node, tag string) ([]*html.Node, error) {
	var Node []*html.Node
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == tag {
			Node = append(Node, node)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	if Node != nil {
		return Node, nil
	}
	return nil, fmt.Errorf("missing %s in the doc", tag)
}

func getNodeContent(node *html.Node) string {
	var ret string
	if node.Type == html.TextNode {
		ret += node.Data
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		ret += getNodeContent(child)
	}
	return ret
}

func GetAuthorStats(steamId string) ([]ModStats, error) {
	doc, err := GetHtml(steamId)
	if err != nil {
		return nil, err
	}
	tBody, err := GetNodesByTag(doc, "tbody")
	if err != nil {
		return nil, err
	}
	table, err := GetNodesByTag(tBody[0], "tr")
	if err != nil {
		return nil, err
	}
	var modStats []ModStats = make([]ModStats, 0)
	for _, v := range table[1:] {
		tds, err := GetNodesByTag(v, "td")
		if err != nil {
			return nil, err
		}
		rankTotal, err := strconv.Atoi(getNodeContent(tds[0]))
		if err != nil {
			return nil, err
		}
		downloadsTotal, err := strconv.Atoi(getNodeContent(tds[2]))
		if err != nil {
			return nil, err
		}
		downloadsYesterday, err := strconv.Atoi(getNodeContent(tds[3]))
		if err != nil {
			return nil, err
		}
		modStats = append(modStats, ModStats{
			RankTotal:          rankTotal,
			DisplayName:        getNodeContent(tds[1]),
			DownloadsTotal:     downloadsTotal,
			DownloadsYesterday: downloadsYesterday,
		})
	}
	return modStats, nil
}
