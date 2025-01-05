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
deps = "npm", "nodejs"
workspace = "https://github.com/scirats/scirats.git"
dots {
	repo = "https://github.com/heaveless/dotfiles.git"
	ext {
		git {
			config {
				user {
					email = "user@scirats.com"
					name = "example name"
				}
			}
			credentials = "https://<user>:<token>@github.com", "https://<user>:<token>@gitlab.com"
		}
	}
}

```

