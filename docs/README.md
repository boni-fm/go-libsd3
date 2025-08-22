# How to Use `libsd3` in Your Go Project

This guide explains how to use the local `libsd3` Go module in your new project.

---

## 1. Clone the `libsd3` Library Locally

```sh
git clone http://<your-gitea-host>/ITSD3/libsd3.git /path/to/libsd3
```
> Replace `/path/to/libsd3` with the directory where you want to store the library.

---

## 2. Create Your New Project

```sh
mkdir /path/to/myapp
cd /path/to/myapp
go mod init myapp
```

---

## 3. Add the `libsd3` Dependency

In your new project's code, import `libsd3` as needed. For example:

```go
import "libsd3/pkg/dbutil"
```

---

## 4. Add a `replace` Directive in Your `go.mod`

Open your project's `go.mod` file and add:

```go
require libsd3 v0.0.0

replace libsd3 => /path/to/libsd3
```
> Make sure `/path/to/libsd3` matches where you cloned the library.

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
├── libsd3/
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
    "libsd3/pkg/dbutil"
)

func main() {
    fmt.Println("Using libsd3!")
    dbutil.DoSomething()
}
```

---

## 8. Troubleshooting

- **Import errors:** Ensure the `replace` path in `go.mod` is correct and absolute.
- **Version errors:** If you update `libsd3`, re-run `go mod tidy` in your new project.

---

## 9. Additional Notes

- You can edit `libsd3` directly and changes will reflect in your project immediately.
- If you want to share `libsd3` with others or use it from a remote repository, update the module path and remove the `replace` directive.

---