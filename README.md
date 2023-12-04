# A very simple WebDAV server written in Go

## Installing

```
go install src.agwa.name/webdavd@latest
```

## Usage

```
Usage of ./webdavd:
  -listen value
    	Socket to listen on, in go-listener syntax (repeatable)
  -public
    	Don't require authentication
  -readwrite
    	Allow read/write access (read-only is the default)
  -root string
    	Path to root directory (required)
  -users string
    	Path to users file (required unless -public is used)
For go-listener syntax, see https://pkg.go.dev/src.agwa.name/go-listener#readme-listener-syntax
Each line of the users file should contain a username and password separated by whitespace
```
