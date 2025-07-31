package branchingstory

import (
	"bufio"
	"fmt"
	"os"
)

type StoryNode struct {
	Text    string
	YesPath *StoryNode
	NoPath  *StoryNode
}

func (node *StoryNode) play() {
	fmt.Println(node.Text)
	if node.YesPath != nil && node.NoPath != nil {
		scanner := bufio.NewScanner(os.Stdin)

		for {
			scanner.Scan()
			answer := scanner.Text()
			if answer == "yes" {
				node.YesPath.play()
				break
			} else if answer == "no" {
				node.NoPath.play()
				break
			} else {
				fmt.Println("That is not an option, please input yes or no.")
			}
		}

	}

}

func Run() {
	root := &StoryNode{"This is the entrance of the dark cave, do you want to step in? ", nil, nil}

	winning := &StoryNode{"You have won.", nil, nil}
	lose := &StoryNode{"You have lose.", nil, nil}

	root.YesPath = lose
	root.NoPath = winning
	root.play()
}
