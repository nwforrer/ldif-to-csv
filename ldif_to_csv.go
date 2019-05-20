package ldiftocsv

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"strings"
	"fmt"
)

type NameValue struct {
	name  string
	value string
}

type LdifEntry struct {
	properties []NameValue
}

func ReadLdifFile(r io.Reader, headers []string) []LdifEntry {
	ldifEntries := []LdifEntry{}

	scanner := bufio.NewScanner(r)
	currentProperties := []NameValue{}
	for scanner.Scan() {
		txt := scanner.Text()
		if len(txt) == 0 {
			if len(currentProperties) > 0 {
				ldifEntries = append(ldifEntries, parseLdifProperties(currentProperties))
				currentProperties = []NameValue{}
			}
		} else {
			for _, header := range headers {
				if strings.HasPrefix(txt, header) {
					value := txt[len(header)+2:]
					currentProperties = addNameValue(currentProperties, header, value)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// ensure the final entry gets recorded
	if len(currentProperties) != 0 {
		ldifEntries = append(ldifEntries, parseLdifProperties(currentProperties))
		currentProperties = []NameValue{}
	}

	return ldifEntries
}

func WriteToCsv(out io.Writer, entries []LdifEntry) {
	w := csv.NewWriter(out)

	lines := ldifEntriesToCsvFormat(entries)
	for _, line := range lines {
		if err := w.Write(line); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

func addNameValue(properties []NameValue, name, value string) []NameValue {
	duplicate := false
	for i, prop := range properties {
		if prop.name == name {
			duplicate = true
			properties[i].value += fmt.Sprintf("%s\n", value)
		}
	}

	if !duplicate {
		properties = append(properties, NameValue{
			name: name,
			value: value,
		})
	}

	return properties
}

func parseLdifProperties(properties []NameValue) LdifEntry {
	entry := LdifEntry{
		properties,
	}
	return entry
}

func ldifEntriesToCsvFormat(entries []LdifEntry) [][]string {
	lines := [][]string{}
	for _, entry := range entries {
		propArr := []string{}
		for _, prop := range entry.properties {
			propArr = append(propArr, prop.value)
		}
		lines = append(lines, propArr)
	}

	return lines
}
