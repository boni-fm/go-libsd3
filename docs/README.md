# How to Use `go-libsd3` in Your Go Project

This guide explains how to use the local `go-libsd3` Go module in your new project.

---

## 1. Clone the `go-libsd3` Library Locally

```sh
git clone http://<your-gitea-host>/ITSD3/go-libsd3.git /path/to/go-libsd3
```
> Replace `/path/to/go-libsd3` with the directory where you want to store the library.

---

## 2. Create Your New Project

```sh
mkdir /path/to/myapp
cd /path/to/myapp
go mod init myapp
```

---

## 3. Add the `go-libsd3` Dependency

In your new project's code, import `go-libsd3` as needed. For example:

```go
import "go-libsd3/pkg/dbutil"
```

---

## 4. Add a `replace` Directive in Your `go.mod`

Open your project's `go.mod` file and add:

```go
require go-libsd3 v0.0.0

replace go-libsd3 => /path/to/go-libsd3
```
> Make sure `/path/to/go-libsd3` matches where you cloned the library.

---

## 5. Run and Build Your Project

Download your dependencies and run your project:

```sh
go mod tidy
go run main.go
```

---

## 6. Example Project Structure

```
/path/to/
├── go-libsd3/
│   ├── go.mod
│   └── pkg/
│       └── dbutil/
│           └── dbutil.go
└── myapp/
    ├── go.mod
    └── main.go
```

---

## 7. Example `main.go`

```go
package main

import (
    "fmt"
    "go-libsd3/pkg/dbutil"
)

func main() {
    fmt.Println("Using go-libsd3!")
    dbutil.DoSomething()
}
```

---

## 8. Troubleshooting

- **Import errors:** Ensure the `replace` path in `go.mod` is correct and absolute.
- **Version errors:** If you update `go-libsd3`, re-run `go mod tidy` in your new project.

---

## 9. Additional Notes

- You can edit `go-libsd3` directly and changes will reflect in your project immediately.
- If you want to share `go-libsd3` with others or use it from a remote repository, update the module path and remove the `replace` directive.

---