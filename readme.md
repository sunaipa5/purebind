# purebind

**purebind** is a CLI tool that parses C-style header files and generates pure Go function bindings using [purego](https://github.com/ebitengine/purego).

[see generated example]()

## Features

- Parses `.h` header files
- Generates Go bindings using [purego](https://github.com/ebitengine/purego)

---

## Installation

```bash
go install github.com/sunaipa5/purebind
```

---

## Usage

```bash
purebind <project-name> <header-file-path>
```

Example:

```bash
purebind mylib ./mylib.h
```

This will generate:

```
mylib/
├── wrapper.go
├── lib.go
└── go.mod
```

---

## Example

Given a C header file:

```c
int add(int a, int b);
void print_message(const char *msg);
```

The generated Go code will look like:

```go
var (
	add func(int32, int32) int32
	print_message func(unsafe.Pointer)
)
```
