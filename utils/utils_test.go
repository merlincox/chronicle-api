package utils

import (
    "testing"
)

func TestSlug(t *testing.T) {

    slug := Slug(" This    is a- 'test?")

    if slug != "this-is-a-test" {
        t.Fail()
    }
}

type testStruct struct {
    StringVar1       string `json:"stringvar-1"`
    StringVar2 string `json:"stringvar-2,omitempty"`
    IntVar1 int `json:"intvar-1"`
    IntVar2 int `json:"intvar-2,omitempty"`
}

func TestJsonStringify(t *testing.T) {

    data1 := testStruct{
        StringVar1: "testing",
        IntVar1: 999,
    }

    string1 := JsonStringify(data1)
    desired1 := `{"stringvar-1":"testing","intvar-1":999}`

    if string1 != desired1 {
        t.Fail()
    }

    data2 := testStruct{
        StringVar1: "testing",
        IntVar1: 999,
        StringVar2: "testing2",
        IntVar2: 666,
    }

    string2 := JsonStringify(data2)
    desired2 := `{"stringvar-1":"testing","stringvar-2":"testing2","intvar-1":999,"intvar-2":666}`

    if string2 != desired2 {
        t.Fail()
    }
}
