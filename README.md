# My MessagePack Library

This repository contains a custom implementation of MessagePack encoding and decoding in Go.

## How to Run
To run the command-line interface (CLI) for this library, execute the following command from the root of the repository:

```sh
# go run ./cmd/msgpack-cli/main.go
```

## How to Test
To run the tests for this library, execute the following command from the root of the repository:

```sh
# go test .
```
* This will execute all the tests in the current module and display the results.

## Testing Different JSON Data
To test different JSON data, you need to modify the main.go file located in the cmd/msgpack-cli/ directory. Update the jsonData variable with your desired JSON data. For example:

1. Open cmd/msgpack-cli/main.go in your text editor.

2. Locate the jsonData variable assignment, which might look like this:
    ```go
    jsonData := []byte(`{"key": "value"}`)
    ```

3. Replace the JSON string with your test data:
    ```go
    jsonData := []byte(`{"newKey": "newValue"}`)
    ```

4. Save the file and run the CLI again using:
    ```sh
    # go run ./cmd/msgpack-cli/main.go
    ```
* This will process the new JSON data and output the result.