# Personal Website

## How to develop

Use commands from `make help`.

### Run

1. Install necessary tools:
   - [Go](https://go.dev/doc/install)
   - [Hugo](https://gohugo.io/installation/)

2. Run the command:

    ```sh
    make serve
    ```

3. Open <http://localhost:1313> in a browser.

### Lint

1. Install necessary tools:
   - [yamlfmt](github.com/google/yamlfmt)

2. Run the command to format files:

    ```sh
    make fmt
    ```

3. Run the command to perform all lints:

    ```sh
    make lint
    ```

### Spell

1. Install necessary tools:
   - [codespell](https://github.com/codespell-project/codespell)

2. Run the command to check spellings:

    ```sh
    make spell
    ```

### Check for dead links

1. Install necessary tools:
   - [linkcheck](https://github.com/filiph/linkcheck)

2. Run the website locally.
3. Execute the command:

    ```sh
    make linkcheck
    ```
