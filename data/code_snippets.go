package data

type CodeSnippet struct {
	Code         string
	RequiredKeys string
	Language     string
}

var Snippets = []CodeSnippet{
	{Code: "if true { return }", RequiredKeys: "iftrue{rn}", Language: "go"},
	{Code: "for i := 0; i < 10; i++ {", RequiredKeys: "fori:=0;<1+{", Language: "go"},
	{Code: "def split(words):", RequiredKeys: "defsplit(wor):", Language: "python"},
	{Code: "result := 0", RequiredKeys: "result:=0", Language: "go"},
	{Code: "err != nil", RequiredKeys: "er!=nil", Language: "go"},
	{Code: "type status struct {}", RequiredKeys: "typesauc{}", Language: "go"},
	{Code: "for k, v := range items {", RequiredKeys: "fork,v:=angeitms{", Language: "go"},
	{Code: "import os", RequiredKeys: "importas", Language: "python"},
	{Code: "while true:", RequiredKeys: "whlietru:", Language: "python"},
	{Code: "print(sorted(list))", RequiredKeys: "pint(sored(ls))", Language: "python"},
}
