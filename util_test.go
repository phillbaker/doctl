package main

import ()

type testOutputableResource string

func (a testOutputableResource) Headers() []string           { return []string{"Header1"} }
func (a testOutputableResource) FormatString() string        { return "%s\t\n" }
func (a testOutputableResource) RowObject(i int) interface{} { return a }
func (a testOutputableResource) RowValues(datum interface{}) []interface{} {
	return []interface{}{datum}
}
func (a testOutputableResource) Len() int { return 1 }

func ExampleTableResourceWriteOutputable() {
	OutputFormat = "table"
	WriteOutputable(testOutputableResource("StringA"))

	// Output:
	// Header1
	// StringA
}

func ExampleJsonResourceWriteOutputable() {
	OutputFormat = "json"
	WriteOutputable(testOutputableResource("StringA"))

	// Output:
	// "StringA"
}

func ExampleYamlResourceWriteOutputable() {
	OutputFormat = "yaml"
	WriteOutputable(testOutputableResource("StringA"))

	// Output:
	// StringA
}

type testOutputableCollection []string

func (a testOutputableCollection) Headers() []string           { return []string{"Header1", "Header2"} }
func (a testOutputableCollection) FormatString() string        { return "%s\t%s\n" }
func (a testOutputableCollection) RowObject(i int) interface{} { return a[i] }
func (a testOutputableCollection) RowValues(datum interface{}) []interface{} {
	return []interface{}{datum, datum}
}
func (a testOutputableCollection) Len() int { return len(a) }

func ExampleTableCollectionWriteOutputable() {
	OutputFormat = "table"
	WriteOutputable(testOutputableCollection([]string{"StringA", "StringB"}))

	// Output:
	// Header1		Header2
	// StringA		StringA
	// StringB		StringB
}

func ExampleJsonCollectionWriteOutputable() {
	OutputFormat = "json"
	WriteOutputable(testOutputableCollection([]string{"StringA", "StringB"}))

	// Output:
	// ["StringA","StringB"]
}

func ExampleYamlCollectionWriteOutputable() {
	OutputFormat = "yaml"
	WriteOutputable(testOutputableCollection([]string{"StringA", "StringB"}))

	// Output:
	// - StringA
	// - StringB
}
