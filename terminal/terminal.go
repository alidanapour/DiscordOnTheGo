// Helper package which changes the colour of the server side CLI via helper functions

// Credit: https://twin.sh/articles/35/how-to-add-colors-to-your-console-terminal-output-in-go
package terminal

import (
	"fmt"
	"runtime"
)

// Enum
type Prefix int

const (
	MESSAGE Prefix = iota
	COMMAND
	ERROR
	SUCCESS
	HTTPGET
)

// Colors
var reset = "\033[0m"
var red = "\033[31m"
var green = "\033[32m"
var yellow = "\033[33m"
var blue = "\033[2;80;119;178m"

//var blue = "\033[34m"
//var purple = "\033[35m"
//var cyan = "\033[36m"
//var gray = "\033[37m"
var white = "\033[97m"

// Windows CLI sometimes doesnt work, this prevents clogging up stdout
func init() {
	if runtime.GOOS == "windows" {
		reset = ""
		red = ""
		green = ""
		yellow = ""
		blue = ""
		//purple = ""
		//cyan = ""
		//gray = ""
		white = ""
	}
}

// Calls fmt.Println with prefix and color
func Print(pre Prefix, str string) {
	var color string
	var prefixStr string

	switch pre {
	case MESSAGE:
		color = white
		prefixStr = "MESSAGE"
	case COMMAND:
		color = yellow
		prefixStr = "COMMAND"
	case ERROR:
		color = red
		prefixStr = "-ERROR-"
	case SUCCESS:
		color = green
		prefixStr = "SUCCESS"
	case HTTPGET:
		color = blue
		prefixStr = "HTTPGET"
	default:
		color = reset
		prefixStr = "--BOT--"
	}

	//result := "[" + color + prefixStr + reset + "]: " + str
	//fmt.Println(result)
	fmt.Printf("[%s%-7s%s] %s\n", color, prefixStr, reset, str)
}
