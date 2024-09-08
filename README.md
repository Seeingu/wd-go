# wd-go

A local html based dev server demo.

## Usage

1. Install CLI

```bash
go install github.com/Seeingu/wd-go@v0.0.1
```

2. Add this script to your html head
```html
<script src="https://raw.githubusercontent.com/Seeingu/wd-go/main/static/reload.js?host=localhost"></script>
```

3. Run 

```bash
# use `wd-go -h` to see all arguments 
wd-go ./index.html
```

open `http://localhost:3012`, edit and preview html in your browser

## Support Features

- [x] reload after change html file 

- [x] reload after change referenced JavaScript/CSS

- [ ] live reload when changed global javascript variable

- [ ] built-in popular libraries support