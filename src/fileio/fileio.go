package fileio

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// WriteStructToFile writes the given struct to the specified file in JSON format.
// Parameters:
//
//	toWrite (interface{}): The struct to be written to the file. It must be serializable to JSON.
//	file (*os.File): A pointer to an os.File object representing the file where the JSON data should be written.
//
// Returns:
//
//	error: An error if any occurs during the marshaling of the struct to JSON or the writing of the JSON data to the file.
//	       If no error occurs, it returns nil.
func WriteStructToFile(toWrite interface{}, file *os.File) error {
	jsonData, err := json.Marshal(toWrite)
	if err != nil {
		return fmt.Errorf("error marshaling struct to JSON: %v", err)
	}

	if _, err = file.Write(jsonData); err != nil {
		return fmt.Errorf("error writing JSON to file: %v", err)
	}

	return nil
}

// WriteStructToFilePath writes a JSON representation of the provided struct
// to the file specified by the given path. It creates a new file if it does
// not exist, or truncates the file to zero length if it exists. The file
// permissions are set to 0644 (readable and writable by owner, readable by others).
//
// Parameters:
//   - toWrite: The struct to be serialized to JSON and written to the file.
//   - path: The file path where the JSON data will be written.
//
// Returns:
//   - error: An error if any operation fails, including file opening, writing,
//     and closing errors.
func WriteStructToFilePath(toWrite interface{}, path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	if err = WriteStructToFile(toWrite, file); err != nil {
		return fmt.Errorf("error writing JSON to file: %v", err)
	}

	return nil
}

// WriteFileToStruct reads the contents of the provided file and parses the JSON data
// into the given data structure.
//
// Parameters:
//   - file: An open file (os.File) containing JSON data to be parsed.
//   - data: A pointer to the target structure where the parsed JSON will be stored.
//
// Returns:
//   - error: Returns an error if reading from the file or parsing JSON fails.
func WriteFileToStruct(file *os.File, data interface{}) error {
	fileContents, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	if err := json.Unmarshal(fileContents, data); err != nil {
		return fmt.Errorf("could not unmarshal JSON: %w", err)
	}

	return nil
}

// WriteFilePathToStruct opens the file specified by the given filename,
// reads its contents, and parses the JSON data into the provided data structure.
// It automatically closes the file when done.
//
// Parameters:
//   - filename: The name of the file containing JSON data to be parsed.
//   - data: A pointer to the target structure where the parsed JSON will be stored.
//
// Returns:
//   - error: Returns an error if opening the file, reading from the file,
//     or parsing JSON fails.
func WriteFilePathToStruct(filename string, data interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	if err := WriteFileToStruct(file, data); err != nil {
		return fmt.Errorf("could not unmarshal JSON: %w", err)
	}

	return nil
}
