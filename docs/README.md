

# How to Use `go-libsd3` in Your Go Project

This guide explains how to use the `go-libsd3` Go module in your project, either from GitHub or from a local repository.

---

## Option 1: Use from GitHub (Recommended)

### 1. Create Your New Project

```sh
mkdir /path/to/myapp
cd /path/to/myapp
go mod init myapp
```

### 2. Add the `go-libsd3` Dependency from GitHub

In your code, import as:

```go
import "github.com/boni-fm/go-libsd3/pkg/dbutil"
```

Then, add the dependency:

```sh
go get github.com/boni-fm/go-libsd3
```

### 3. Run and Build Your Project

```sh
go mod tidy
go run main.go
```

---

## Option 2: Use from Local Repository (Advanced/Development)

### 1. Clone the `go-libsd3` Library Locally

```sh
git clone http://url-ke-gitea/ITSD3/go-libsd3.git /path/to/go-libsd3
```

### 2. Create Your New Project

```sh
mkdir /path/to/myapp
cd /path/to/myapp
go mod init myapp
```

### 3. Add the Dependency and Replace Directive

In your code, import as:

```go
import "go-libsd3/pkg/dbutil"
```

In your `go.mod`:

```go
require go-libsd3 v0.0.0
replace go-libsd3 => /path/to/go-libsd3
```

### 4. Run and Build Your Project

```sh
go mod tidy
go run main.go
```

---

## Example Project Structure

```
/path/to/myapp/
├── go.mod
└── main.go
```

---

## Example `main.go`

```go
package main

import (
    "fmt"
    // For GitHub import:
    // "github.com/boni-fm/go-libsd3/pkg/dbutil"
    // For local import:
    // "go-libsd3/pkg/dbutil"
)

func main() {
    fmt.Println("Using go-libsd3!")
    // dbutil.DoSomething()
}
```

---

## Troubleshooting

- **Import errors:** Ensure the import path and/or `replace` directive are correct.
- **Version errors:** If you update `go-libsd3`, re-run `go mod tidy` in your new project.

---

## Additional Notes

- You can always update to the latest version with `go get -u github.com/boni-fm/go-libsd3`.
- If you want to use a specific version, specify it in your `go.mod`.
- If you use the local repo, edits to `go-libsd3` will reflect immediately in your project.

---