package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"github.com/MustafaLo/noted/models"
)


func GetCurrentFileMetadata() (models.FileMetadata, error) {
    var metadata models.FileMetadata
    data, err := os.ReadFile("fileMetadata.json")
    if err != nil {
        return metadata, fmt.Errorf("failed to open fileMetaData -- make sure to enable File Tracker Extension")
    }

    err = json.Unmarshal(data, &metadata)
    if err != nil {
        return metadata, fmt.Errorf("failed to parse fileMetaData: %w", err)
    }

    return metadata, nil
}

func PrintFileMetaData(metadata models.FileMetadata) {
    fmt.Printf("File Name: %s\n", metadata.FileName)
    fmt.Printf("File Path: %s\n", metadata.FilePath)
    fmt.Printf("Lines: Start=%d, End=%d\n", metadata.Lines.Start, metadata.Lines.End)
    fmt.Printf("Timestamp: %s\n", metadata.Timestamp)
	fmt.Printf("Language: %s\n", metadata.Language)
}