package main

import (
	"fmt"
)

type storyPage struct {
	text     string
	nextPage *storyPage
}

// func (page *storyPage) playPage() {
// 	if page == nil {
// 		return
// 	}
// 	fmt.Println(page.text)
// 	page.nextPage.playPage()
// }

func (page *storyPage) playPage() {
	for page != nil {
		fmt.Println(page.text)
		page = page.nextPage
	}

}

func (page *storyPage) addToEnd(text string) {
	for page.nextPage != nil {
		page = page.nextPage
	}
	page.nextPage = &storyPage{text, nil}

}

func (page *storyPage) addAfterHead(text string) {
	newPage := &storyPage{text, page.nextPage}
	page.nextPage = newPage
}

func main() {

	page1 := &storyPage{"It was a dark and stormy night.", nil}
	page1.addToEnd("You are alone, and you need to find the sacred helmet before the bad guys do")
	page1.addToEnd("You see a troll ahead")

	page1.playPage()
	page1.addAfterHead("add after head")
	page1.playPage()

}
