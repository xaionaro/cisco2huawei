package c2h

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Vlan struct {
	Id      int
	Enabled bool   `json:",omitempty"`
	Name    string `json:",omitempty"`
}

type Configuration struct {
	Vlan []Vlan `json:",omitempty"`
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func readLine(reader *bufio.Reader) (string, bool) {
	lineBytes, isPrefix, err := reader.ReadLine()
	if err != nil {
		return "", true
	}
	if isPrefix {
		panic("line is too long")
	}
	return string(lineBytes), false
}

func ParseCiscoFile(in io.Reader) (config Configuration, err error) {
	reader := bufio.NewReader(in)

	for {
		line, eof := readLine(reader)
		if eof {
			break
		}
		words := strings.Split(line, " ")
		section := words[0]

		switch section {
		case "vlan":
			if words[1] == "internal" {
				continue
			}
			vlanId, err := strconv.Atoi(words[1])
			checkErr(err)

			var vlan Vlan
			vlan.Id = vlanId
			vlan.Enabled = true
			for {
				line, eof := readLine(reader)
				if eof {
					break
				}
				if line == "!" {
					break
				}
				words := strings.Split(line, " ")
				switch words[1] {
				case "name":
					vlan.Name = strings.Join(words[2:], " ")
				}
			}
			config.Vlan = append(config.Vlan, vlan)
		}
	}

	return
}

func WriteToHuaweiFile(out io.Writer, config Configuration) (err error) {
	for _, vlan := range config.Vlan {
		if !vlan.Enabled {
			continue
		}
		if vlan.Name == "" {
			vlan.Name = fmt.Sprintf("vlan%v", vlan.Id)
		}
		_, err = fmt.Fprintf(out, "vlan %v\nname %v\nq\n", vlan.Id, vlan.Name)
		checkErr(err)
	}
	return

}
