# Cube

**Init Project**
```sh
cube init <name>
```

**Start Project**
```sh
cube start <file>
```

## Configuration file
*example*

```
ports = "3000:3000"
dots = "my-user/dotfiles"
source {
	repo = "https://github.com/scirats/scirats.git"
	name = "my-user"
	email = "my-user@example.com"
	token = "my-token"
}
container = `
	FROM alpine:latest

	RUN apk update && \
		apk add --no-cache neovim git
	
	CMD ["sleep", "infinity"]
`
```

