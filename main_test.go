package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestListAllFiles(t *testing.T) {
	data := []struct {
		FilePath  string
		Want      int
		WantError error
	}{
		{
			FilePath:  "D:/Dummy",
			Want:      0,
			WantError: nil,
		},
		{
			FilePath:  "D:/files",
			Want:      2,
			WantError: nil,
		},
		{
			FilePath:  "D:/Duuumy",
			WantError: errors.New("directory does not exist"),
		},
	}

	for _, item := range data {
		got, err := listAllFiles(item.FilePath)
		if err != nil {
			gotError := err.Error()
			wantError := ""
			if item.WantError != nil {
				wantError = item.WantError.Error()
			}
			if !reflect.DeepEqual(gotError, wantError) {
				t.Errorf("Got %v, want %v", gotError, wantError)
			}
		}

		want := item.Want
		if len(got) != want {
			t.Errorf("Got %v, want %v", got, want)
		}
	}
}

func TestSearchFiles(t *testing.T) {
	data := []struct {
		FilePath  string
		FileName  string
		Want      string
		WantError error
	}{
		{
			FilePath:  "D:/Dummy",
			FileName:  "doc.pdf",
			Want:      "file is not present in the given directory",
			WantError: nil,
		},
		{
			FilePath:  "D:/files",
			FileName:  "Data.jpg",
			Want:      "file is present in the given directory",
			WantError: nil,
		},
		{
			FilePath:  "D:/Duuumy",
			FileName:  "Data.jpg",
			WantError: errors.New("directory does not exist"),
		},
	}

	for _, item := range data {
		got, err := searchFiles(item.FilePath, item.FileName)
		if err != nil && item.WantError != nil {
			gotError := err.Error()
			wantError := item.WantError.Error()
			if !reflect.DeepEqual(gotError, wantError) {
				t.Errorf("Got %v, want %v", gotError, wantError)
			}
		}

		want := item.Want
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Got %v, want %v", got, want)
		}
	}
}

func TestCopyAndMoveFile(t *testing.T) {
	data := []struct {
		SourceFilePath      string
		DestinationFilePath string
		FileName            string
		Operation           string
		Want                string
		WantError           error
	}{
		{
			SourceFilePath:      "D:/files",
			DestinationFilePath: "D:/Dummy3",
			FileName:            "test_gc.pdf",
			Operation:           "copy",
			Want:                "File copied successfully!",
			WantError:           nil,
		},
		{
			SourceFilePath:      "D:/Dummy2",
			DestinationFilePath: "D:/files",
			FileName:            "test_gc.pdf",
			Operation:           "move",
			Want:                "File moved successfully!",
			WantError:           nil,
		},
		{
			SourceFilePath:      "D:/D",
			DestinationFilePath: "",
			Operation:           "copy",
			FileName:            "doc.pdf",
			WantError:           errors.New("failed to open doc.pdf file in given directory D:/D"),
		},
		{
			SourceFilePath:      "D:/files",
			DestinationFilePath: "D:/Dummy2",
			FileName:            "test_gc.pdf",
			Operation:           "copy",
			Want:                "File copied successfully!",
			WantError:           nil,
		},
		{
			SourceFilePath:      "C:/",
			DestinationFilePath: "D:/Dummy2",
			FileName:            "test_gc.pdf",
			Operation:           "move",
			Want:                "",
			WantError:           errors.New("remove C://test_gc.pdf: Access is denied."),
		},
		{
			SourceFilePath:      "D:/Dummy2",
			DestinationFilePath: "E:/Dummy2",
			FileName:            "test_gc.pdf",
			Operation:           "copy",
			Want:                "",
			WantError:           errors.New("mkdir E:: The system cannot find the path specified."),
		},
		{
			SourceFilePath:      "D:/files",
			DestinationFilePath: "D:/Dummy2",
			FileName:            "test_gc.pdf",
			Operation:           "d",
			Want:                "",
			WantError:           errors.New("invalid operation d"),
		},
		{
			SourceFilePath:      "D:/files",
			DestinationFilePath: "D:/Dummy3",
			FileName:            "",
			Operation:           "copy",
			Want:                "Directory copied successfully",
			WantError:           nil,
		},
		{
			SourceFilePath:      "D:/Dummy2",
			DestinationFilePath: "D:/Dummy3",
			FileName:            "",
			Operation:           "move",
			Want:                "Directory moved successfully",
			WantError:           nil,
		},
	}

	for _, item := range data {
		got, err := copyAndMoveFile(item.SourceFilePath, item.DestinationFilePath, item.FileName, item.Operation)
		if err != nil {
			gotError := err.Error()
			wantError := ""
			if item.WantError != nil {
				wantError = item.WantError.Error()
			}
			if !reflect.DeepEqual(gotError, wantError) {
				t.Errorf("Got %v, want %v", gotError, wantError)
			}
		}

		want := item.Want
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Got %v, want %v", got, want)
		}
	}
}

func TestDeleteFile(t *testing.T) {
	data := []struct {
		FilePath  string
		FileName  string
		Want      string
		WantError error
	}{
		{
			FilePath:  "D:/Dummy",
			FileName:  "doc.pdf",
			WantError: errors.New("file is not present in the given directory"),
		},
		{
			FilePath:  "D:/Dummy3",
			FileName:  "test_gc.pdf",
			Want:      "File deleted successfully",
			WantError: nil,
		},
		{
			FilePath:  "D:/Duuumy",
			FileName:  "Data.jpg",
			Want:      "",
			WantError: errors.New("file is not present in the given directory"),
		},
		{
			FilePath:  "C:/",
			FileName:  "test_gc.pdf",
			Want:      "",
			WantError: errors.New("remove C://test_gc.pdf: Access is denied."),
		},
		{
			FilePath:  "D:/Dummy3",
			FileName:  "",
			Want:      "Directory deleted successfully",
			WantError: nil,
		},
	}
	for _, item := range data {
		got, err := deleteFile(item.FilePath, item.FileName)
		if err != nil && item.WantError != nil {
			gotError := err.Error()
			wantError := item.WantError.Error()
			if !reflect.DeepEqual(gotError, wantError) {
				t.Errorf("Got %v, want %v", gotError, wantError)
			}
		}

		want := item.Want
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Got %v, want %v", got, want)
		}
	}
}

func TestCreateFolderIfNotExist(t *testing.T) {
	data := []struct {
		FilePath  string
		WantError error
	}{
		{
			FilePath:  "D:/Dummy2",
			WantError: nil,
		},
		{
			FilePath:  "E:/Dummy2",
			WantError: errors.New("mkdir E:: The system cannot find the path specified."),
		},
	}

	for _, item := range data {
		err := createFolderIfNotExist(item.FilePath)
		if err != nil {
			gotError := err.Error()
			wantError := ""
			if item.WantError != nil {
				wantError = item.WantError.Error()
			}
			if !reflect.DeepEqual(gotError, wantError) {
				t.Errorf("Got %v, want %v", gotError, wantError)
			}
		}
	}
}

func TestCopyDirectory(t *testing.T) {
	data := []struct {
		SourceFilePath      string
		DestinationFilePath string
		Operation           string
		WantError           error
	}{
		{
			SourceFilePath:      "D:/files",
			DestinationFilePath: "D:/NewFiles",
			Operation:           "copy",
			WantError:           nil,
		},
		{
			SourceFilePath:      "D:/D",
			DestinationFilePath: "",
			Operation:           "copy",
			WantError:           errors.New("CreateFile D:/D: The system cannot find the file specified."),
		},
		{
			SourceFilePath:      "D:/Dummy2",
			DestinationFilePath: "E:/Dummy2",
			Operation:           "copy",
			WantError:           errors.New("mkdir E:: The system cannot find the path specified."),
		},
		{
			SourceFilePath:      "D:/NewFiles",
			DestinationFilePath: "D:/NewFiles2",
			Operation:           "move",
			WantError:           nil,
		},
	}

	for _, item := range data {
		err := copyDirectory(item.SourceFilePath, item.DestinationFilePath, item.Operation)
		gotError := ""
		wantError := ""
		if err != nil {
			gotError = err.Error()
		}
		if item.WantError != nil {
			wantError = item.WantError.Error()
		}
		if !reflect.DeepEqual(gotError, wantError) {
			t.Errorf("Got %v, want %v", gotError, wantError)
		}
	}
}

func TestCopyFileDirectory(t *testing.T) {
	data := []struct {
		SourceFilePath      string
		DestinationFilePath string
		WantError           error
	}{
		{
			SourceFilePath:      "D:/files/test_gc.pdf",
			DestinationFilePath: "D:/NewFiles2/test_gc.pdf",
			WantError:           nil,
		},
		{
			SourceFilePath:      "D:/D",
			DestinationFilePath: "",
			WantError:           errors.New("open D:/D: The system cannot find the file specified."),
		},
		{
			SourceFilePath:      "D:/Dummy2",
			DestinationFilePath: "E:/Dummy2",
			WantError:           errors.New("open E:/Dummy2: The system cannot find the path specified."),
		},
	}

	for _, item := range data {
		err := copyFileDirectory(item.SourceFilePath, item.DestinationFilePath)
		gotError := ""
		wantError := ""
		if err != nil {
			gotError = err.Error()
		}
		if item.WantError != nil {
			wantError = item.WantError.Error()
		}
		if !reflect.DeepEqual(gotError, wantError) {
			t.Errorf("Got %v, want %v", gotError, wantError)
		}
	}
}

func TestDeleteDirectory(t *testing.T) {
	data := []struct {
		FilePath  string
		WantError error
	}{
		{
			FilePath:  "D:/Duumy",
			WantError: errors.New("directory not found"),
		},
		{
			FilePath:  "D:/NewFiles2",
			WantError: nil,
		},
		{
			FilePath:  "C:/test_gc.pdf",
			WantError: errors.New("remove C:/test_gc.pdf: Access is denied."),
		},
	}
	for _, item := range data {
		err := deleteDirectory(item.FilePath)
		gotError := ""
		wantError := ""
		if err != nil {
			gotError = err.Error()
		}
		if item.WantError != nil {
			wantError = item.WantError.Error()
		}
		if !reflect.DeepEqual(gotError, wantError) {
			t.Errorf("Got %v, want %v", gotError, wantError)
		}
	}
}
