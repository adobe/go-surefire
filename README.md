# Surefire Result Parser

Reads Surefire Results from multiple XML files and aggregates then for further processing.

## Features
Surefire XML Test results are parsed into structs. [Data Model](doc/datamodel.md) shows 
the data-structure surefire results is converted to.

## Installation

- Install Go, at least version 1.21.0
- Run `make local-build`. This will resolve dependencies and run tests

## Usage

Read test-result XML files from `~/test-results` directory.

```
files, _ := listTestReportFiles()

testResults, err := NewJUnitReportsReaderBuilder().Build().FromReportFiles(files)
	
func listTestReportFiles() ([]string, error) {
	sureFireTestResult := []string{}
	items, err := os.ReadDir(~/test-results)

	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.Type().IsRegular() && strings.HasSuffix(item.Name(), ".xml") {
			sureFireTestResult = append(sureFireTestResult, fmt.Sprintf("%s/%s", path, item.Name()))
		}
	}

	return sureFireTestResult, nil
}
```

There is also support for adding labels to results on Suite level. 

```
files, _ := listTestReportFiles()

testResults, err := NewJUnitReportsReaderBuilder().WithLabeler(func(suite TestSuite) []string {
		return []string{"label"}
	}).Build().FromReportFiles(files)
```
### Contributing

Contributions are welcomed! Read the [Contributing Guide](./.github/CONTRIBUTING.md) for more information.

### Licensing

This project is licensed under the Apache V2 License. See [LICENSE](LICENSE) for more information.
