package main

import (
	"bufio"
	"fmt"
	"os"
)

type storyNode struct {
	text    string
	yesPath *storyNode
	noPath  *storyNode
}

func (node *storyNode) play() {
	fmt.Println(node.text)
	if node.yesPath != nil && node.noPath != nil {
		scanner := bufio.NewScanner(os.Stdin)

		for {
			scanner.Scan()
			answer := scanner.Text()
			if answer == "yes" {
				node.yesPath.play()
				break
			} else if answer == "no" {
				node.noPath.play()
				break
			} else {
				fmt.Println("That is not an option, please input yes or no.")
			}
		}

	}

}

func main() {
	root := &storyNode{"This is the entrance of the dark cave, do you want to step in? ", nil, nil}

	winning := &storyNode{"You have won.", nil, nil}
	lose := &storyNode{"You have lose.", nil, nil}

	root.yesPath = lose
	root.noPath = winning
	root.play()
}
