# pkg/fileutil

Reusable file-system helpers for Go services. Covers general file operations, CSV marshalling, ZIP archiving, and FTP transfers. All functions that touch the file system accept an `afero.Fs` interface — pass `afero.NewMemMapFs()` in tests and `afero.NewOsFs()` in production for full testability.

---

## General file operations (`file.go`)

| Function | Description |
|---|---|
| `EnsureDir(fs, dirPath)` | Creates `dirPath` and all parents. Safe to call when it already exists. |
| `FileExists(fs, path)` | Returns `true` if `path` is a regular file. |
| `DirExists(fs, path)` | Returns `true` if `path` is a directory. |
| `CopyFile(fs, src, dst)` | Copies `src` → `dst`, creating parent directories as needed. |
| `ReadLines(fs, path)` | Reads a UTF-8 text file and returns lines as `[]string`. |
| `WriteLines(fs, path, lines)` | Writes `lines` to `path`, separated by newlines. |

```go
fs := afero.NewOsFs()

if err := fileutil.EnsureDir(fs, "/data/output"); err != nil { ... }

if fileutil.FileExists(fs, "/data/input.txt") {
    lines, _ := fileutil.ReadLines(fs, "/data/input.txt")
    fileutil.WriteLines(fs, "/data/output.txt", lines)
}
```

---

## CSV helpers (`csv.go`)

| Function | Description |
|---|---|
| `WriteStructsToCSV[T](data, filePath)` | Marshals `[]T` (gocsv-tagged structs) to a CSV file on the OS filesystem. |
| `WriteStructsToCSVFs[T](fs, data, filePath)` | afero-aware version. |
| `ReadCSVToStructs[T](filePath)` | Reads a CSV file and unmarshals rows into `[]T`. |
| `ReadCSVToStructsFs[T](fs, filePath)` | afero-aware version. |
| `WriteSliceToCSV(data, filePath)` | Writes `[][]string` rows to CSV (first row is the header). |
| `WriteSliceToCSVFs(fs, data, filePath)` | afero-aware version. |

```go
type Record struct {
    Name  string `csv:"name"`
    Score int    `csv:"score"`
}

records := []Record{{Name: "Alice", Score: 90}, {Name: "Bob", Score: 85}}
if err := fileutil.WriteStructsToCSV(records, "/out/scores.csv"); err != nil { ... }

loaded, err := fileutil.ReadCSVToStructs[Record]("/out/scores.csv")
```

---

## ZIP helpers (`zip.go`)

| Function | Description |
|---|---|
| `ZipFiles(destPath, sourcePaths)` | Creates a ZIP archive from a list of files/directories on the OS filesystem. |
| `ZipFilesFs(fs, destPath, sourcePaths)` | afero-aware version. |
| `UnzipTo(srcPath, destDir)` | Extracts a ZIP archive into `destDir` on the OS filesystem. Protected against zip-slip attacks. |
| `UnzipToFs(fs, srcPath, destDir)` | afero-aware version. |

```go
// Create archive
if err := fileutil.ZipFiles("/out/bundle.zip", []string{"/data/report.csv", "/data/images"}); err != nil { ... }

// Extract archive
if err := fileutil.UnzipTo("/out/bundle.zip", "/data/extracted"); err != nil { ... }
```

---

## FTP client (`ftp.go`)

`FTPClient` wraps an authenticated FTP connection. Not safe for concurrent use — create one per goroutine or synchronise externally.

```go
client, err := fileutil.NewFTPClient(fileutil.FTPConfig{
    Host:     "ftp.example.com",
    Port:     21,           // optional, defaults to 21
    Username: "user",
    Password: "pass",
    Timeout:  30 * time.Second, // optional, defaults to 30s
})
if err != nil { ... }
defer client.Close()

// Upload
if err := client.Upload("/local/report.csv", "/remote/reports/report.csv"); err != nil { ... }

// Download
if err := client.Download("/remote/data.zip", "/local/data.zip"); err != nil { ... }

// List directory
names, err := client.List("/remote/reports")

// Delete
if err := client.Delete("/remote/old.csv"); err != nil { ... }

// Create remote directory
if err := client.MakeDir("/remote/new-folder"); err != nil { ... }
```
