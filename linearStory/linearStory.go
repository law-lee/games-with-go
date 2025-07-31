package linestore

import (
	"fmt"
)

type StoryPage struct {
	Text     string
	NextPage *StoryPage
}

// func (page *StoryPage) playPage() {
// 	if page == nil {
// 		return
// 	}
// 	fmt.Println(page.Text)
// 	page.NextPage.playPage()
// }

func (page *StoryPage) playPage() {
	for page != nil {
		fmt.Println(page.Text)
		page = page.NextPage
	}

}

func (page *StoryPage) addToEnd(Text string) {
	for page.NextPage != nil {
		page = page.NextPage
	}
	page.NextPage = &StoryPage{Text, nil}

}

func (page *StoryPage) addAfterHead(Text string) {
	newPage := &StoryPage{Text, page.NextPage}
	page.NextPage = newPage
}

func Run() {

	page1 := &StoryPage{"It was a dark and stormy night.", nil}
	page1.addToEnd("You are alone, and you need to find the sacred helmet before the bad guys do")
	page1.addToEnd("You see a troll ahead")

	page1.playPage()
	page1.addAfterHead("add after head")
	page1.playPage()

}
