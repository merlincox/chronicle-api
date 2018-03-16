package utils

import (
    "regexp"
    "strings"
    "encoding/json"
    "fmt"
)

var re = regexp.MustCompile("[^a-z0-9]+")

func Slug(s string) string {
    return strings.Trim(re.ReplaceAllString(strings.ToLower(s), "-"), "-")
}

func JsonStringify(data interface{}) string {

    raw, err := json.Marshal(data)

    if err != nil {
        return fmt.Sprintf("JSON parse error: %v", err)
    }

    return string(raw[:])
}

