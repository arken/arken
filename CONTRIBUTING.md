## Go Instructions

1. Install `golang` from your distribution repository.
2. Clone the repository `git clone https://github.com/arken/arken`
3. Type `go run arken.go` (This will pull down all dependencies and start the program.)
4. You're ready to go :P

## Configuration

##### arken.config (Main configuration file)

Arken stores it's `arken.config` file in a `.arken` folder in your home directory. This configuration file is written in TOML for an easy to understand format.

##### keysets.yaml

Arken stores the list of keysets it watches in an easily editable keysets.yaml file within the same `.arken` directory. To add a keyset simply add another line beginning with a dash.

```yaml
keysets:
- https://github.com/arken/core-keyset-testing
- https://github.com/alecbcs/core-fork
```
