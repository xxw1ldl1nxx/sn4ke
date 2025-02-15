package sn4ke

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/xxw1ldl1nxx/lynn"
)

func Stdout(matrix [][]int, score int) {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()

	emptyColor := color.New(color.FgYellow).SprintFunc()
	snakeColor := color.New(color.FgGreen).SprintFunc()
	appleColor := color.New(color.FgRed).SprintFunc()
	borderColor := color.New(color.FgCyan).SprintFunc()

	border := "+"
	for range matrix[0] {
		border += "--"
	}
	border += "-+"
	border = borderColor(border)
	sep := borderColor("|")

	fmt.Println(border)
	for _, intLine := range matrix {
		strLine := lynn.Map(intLine, func(v int) string {
			switch v {
			case EMPTY:
				return emptyColor("+")
			case HEAD:
				return snakeColor("H")
			case TAIL:
				return snakeColor("S")
			case APPLE:
				return appleColor("A")
			default:
				return "$!"
			}
		},
		)
		str := strings.Join(strLine, " ")
		fmt.Println(sep, str, sep)
	}
	fmt.Println(border)
	fmt.Printf("\nScore: %d\n\n", score)
}

func Stdin() Direction {

	var button string
	for {
		if _, err := fmt.Scan(&button); err != nil {
			fmt.Println(err)
			continue
		}
		dir, ok := dirButton[button]
		if ok {
			return dir
		}
		fmt.Println("wrong button")
	}
}
