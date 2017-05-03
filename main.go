package main

import (
	"encoding/json"
	"fmt"
	"github.com/droundy/goopt"
	"github.com/xaionaro/cisco2huawei/c2h"
	"os"
	"strings"
)

type format int

const (
	FORMAT_UNDEFINED format = 0
	FORMAT_JSON      format = 1
	FORMAT_CISCO     format = 2
	FORMAT_HUAWEI    format = 3
)

var formatMap = map[string]format{"json": FORMAT_JSON, "cisco": FORMAT_CISCO, "huawei": FORMAT_HUAWEI}

func parseFromTo(out *format, arg string) (err error) {
	*out = formatMap[arg]
	if *out == 0 {
		return fmt.Errorf("Invalid format name: \"%v\"", arg)
	}

	return nil
}

func main() {
	var from format
	var to format

	var possibleToFromValuesSlice []string
	for formatName, _ := range formatMap {
		possibleToFromValuesSlice = append(possibleToFromValuesSlice, formatName)
	}
	possibleToFromValuesString := strings.Join(possibleToFromValuesSlice, ", ")

	goopt.ReqArg([]string{"-f", "--from"}, "formatName", "format of data in stdin   (possible values: "+possibleToFromValuesString+")", func(arg string) error { return parseFromTo(&from, arg) })
	goopt.ReqArg([]string{"-t", "--to"}, "formatName", "format of data for stdout (possible values: "+possibleToFromValuesString+")", func(arg string) error { return parseFromTo(&to, arg) })

	goopt.Description = func() string {
		return "Converter of cisco configuration to huawei configuration"
	}
	goopt.Version = "0.0"
	goopt.Summary = "cisco2huawei"
	goopt.Parse(nil)

	if from == FORMAT_UNDEFINED || to == FORMAT_UNDEFINED {
		panic(fmt.Errorf("--from or --to is not set"))
		os.Exit(-1)
	}

	var configuration c2h.Configuration
	switch from {
	case FORMAT_JSON:
		decoder := json.NewDecoder(os.Stdin)
		decoder.Decode(&configuration)

	case FORMAT_CISCO:
		var err error
		configuration, err = c2h.ParseCiscoFile(os.Stdin)
		if err != nil {
			panic(err)
		}

	case FORMAT_HUAWEI:
		panic("Not implemented, yet")
	}

	switch to {
	case FORMAT_JSON:
		b, err := json.Marshal(configuration)
		if err != nil {
			panic(err)
		}
		os.Stdout.Write(b)

	case FORMAT_CISCO:
		panic("Not implemented, yet")

	case FORMAT_HUAWEI:
		err := c2h.WriteToHuaweiFile(os.Stdout, configuration)
		if err != nil {
			panic(err)
		}
	}

	return
}
