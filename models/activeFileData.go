package models

type FileMetadata struct {
    FileName  string `json:"fileName"`
    FilePath  string `json:"filePath"`
    Lines     struct {
        Start int `json:"start"`
        End   int `json:"end"`
    } `json:"lines"`
    Timestamp string `json:"timestamp"`
	Language string `json:"language"`
}