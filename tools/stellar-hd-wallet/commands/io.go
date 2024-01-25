package commands

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)
var out io.Writer = os.Stdout

func readString() string {
	line, _ := reader.ReadString('\n')
	return strings.TrimRight(line, "\r\n")
}

func readUint() uint32 {
	line := readString()
	number, err := strconv.Atoi(line)
	if err != nil {
		log.Fatal("Invalid value")
	}

	return uint32(number)
}

func printf(format string, a ...interface{}) {
	fmt.Fprintf(out, format, a...)
}

func println(a ...interface{}) {
	fmt.Fprintln(out, a...)
}
