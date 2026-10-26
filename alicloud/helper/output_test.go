package helper

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests copied from alicloud/common_test.go; test names are unchanged.
func TestUnitCommonWriteToFile(t *testing.T) {
	tmpFile := fmt.Sprintf("%s/testfile.txt", t.TempDir())

	t.Run("write content", func(t *testing.T) {
		err := WriteToFile(tmpFile, "test content")
		assert.NoError(t, err)
		data, _ := os.ReadFile(tmpFile)
		assert.Equal(t, "test content", string(data))
	})

	t.Run("write json", func(t *testing.T) {
		content := map[string]interface{}{"key": "value"}
		err := WriteToFile(tmpFile, content)
		assert.NoError(t, err)
		data, _ := os.ReadFile(tmpFile)
		var result map[string]interface{}
		json.Unmarshal(data, &result)
		assert.Equal(t, content["key"], result["key"])
	})
}

func TestUnitCommonWriteToFileWithHomeDir(t *testing.T) {
	// Test with ~ in path
	tmpDir := t.TempDir()
	// We can't directly test ~ expansion without mocking, but we can test the logic
	err := WriteToFile(fmt.Sprintf("%s/test.txt", tmpDir), "content")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test with nil data
	err = WriteToFile(fmt.Sprintf("%s/test2.txt", tmpDir), nil)
	if err != nil {
		t.Errorf("Unexpected error for nil: %v", err)
	}

	// Test overwriting existing file
	err = WriteToFile(fmt.Sprintf("%s/test.txt", tmpDir), "new content")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestUnitCommonWriteToFileEdgeCases(t *testing.T) {
	tmpDir := t.TempDir()

	// Test with map data
	mapData := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	}
	err := WriteToFile(fmt.Sprintf("%s/map.json", tmpDir), mapData)
	if err != nil {
		t.Errorf("Unexpected error writing map: %v", err)
	}

	// Test with struct data
	structData := struct {
		Name string
		Age  int
	}{
		Name: "test",
		Age:  30,
	}
	err = WriteToFile(fmt.Sprintf("%s/struct.json", tmpDir), structData)
	if err != nil {
		t.Errorf("Unexpected error writing struct: %v", err)
	}

	// Verify file was created and can be read
	content, err := os.ReadFile(fmt.Sprintf("%s/struct.json", tmpDir))
	if err != nil {
		t.Errorf("Failed to read written file: %v", err)
	}
	if len(content) == 0 {
		t.Error("Written file is empty")
	}
}

func TestUnitCommonGetUserHomeDir(t *testing.T) {
	// Test getting user home directory
	homeDir, err := GetUserHomeDir()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if homeDir == "" {
		t.Error("Expected non-empty home directory")
	}
	// Verify it's an absolute path
	if !strings.HasPrefix(homeDir, "/") && !strings.Contains(homeDir, ":\\") {
		t.Errorf("Expected absolute path, got %s", homeDir)
	}
}

func TestUnitCommonWriteToFileComprehensive(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name     string
		filePath string
		data     interface{}
		wantErr  bool
	}{
		{
			name:     "write string",
			filePath: fmt.Sprintf("%s/string.txt", tmpDir),
			data:     "test string content",
			wantErr:  false,
		},
		{
			name:     "write nil",
			filePath: fmt.Sprintf("%s/nil.txt", tmpDir),
			data:     nil,
			wantErr:  false,
		},
		{
			name:     "write map",
			filePath: fmt.Sprintf("%s/map.json", tmpDir),
			data:     map[string]interface{}{"key": "value", "num": 123},
			wantErr:  false,
		},
		{
			name:     "write struct",
			filePath: fmt.Sprintf("%s/struct.json", tmpDir),
			data: struct {
				Name string
				Age  int
			}{"test", 30},
			wantErr: false,
		},
		{
			name:     "write slice",
			filePath: fmt.Sprintf("%s/slice.json", tmpDir),
			data:     []string{"a", "b", "c"},
			wantErr:  false,
		},
		{
			name:     "unmarshalable data",
			filePath: fmt.Sprintf("%s/channel.json", tmpDir),
			data:     make(chan int),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WriteToFile(tt.filePath, tt.data)

			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Verify file was created for non-error cases
			if !tt.wantErr && tt.data != nil {
				if _, err := os.Stat(tt.filePath); os.IsNotExist(err) {
					t.Error("File was not created")
				}
			}
		})
	}

	// Test overwriting existing file
	t.Run("overwrite existing file", func(t *testing.T) {
		filePath := fmt.Sprintf("%s/overwrite.txt", tmpDir)

		// Write first time
		err := WriteToFile(filePath, "first content")
		if err != nil {
			t.Fatalf("First write failed: %v", err)
		}

		// Write second time (overwrite)
		err = WriteToFile(filePath, "second content")
		if err != nil {
			t.Fatalf("Overwrite failed: %v", err)
		}

		// Verify content was updated
		content, _ := os.ReadFile(filePath)
		if string(content) != "second content" {
			t.Errorf("Expected 'second content', got %s", string(content))
		}
	})
}

func TestUnitCommonDataResourceIdHashFromMap(t *testing.T) {
	first := DataResourceIdHashFromMap(map[string]interface{}{"a": "x", "b": 2})
	second := DataResourceIdHashFromMap(map[string]interface{}{"b": 2, "a": "x"})
	assert.Equal(t, first, second)
	assert.NotEmpty(t, first)

	other := DataResourceIdHashFromMap(map[string]interface{}{"a": "y", "b": 2})
	assert.NotEqual(t, first, other)
}
