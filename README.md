tag - Tag your ag matches
====
![revolv++](tag.gif)

**tag** is a lightweight wrapper around **[ag](https://github.com/ggreer/the_silver_searcher)** and **[ripgrep](https://github.com/BurntSushi/ripgrep)** that generates shell aliases for your search matches. tag is a very fast Golang reimagining of [sack](https://github.com/sampson-chen/sack).

tag supports `ag` and `ripgrep (rg)`. There are no plans to support ack or grep. If you'd like to add support for more search backends, I encourage you to contribute!

## Why should I use tag?

tag makes it easy to _immediately_ jump to a search match in your favorite editor. It eliminates the tedious task of typing `vim foo/bar/baz.qux +42` to jump to a match by automatically generating these commands for you as shell aliases.

Inside vim, [vim-grepper](https://github.com/mhinz/vim-grepper) or [ag.vim](https://github.com/rking/ag.vim) is probably the way to go. Outside vim (or inside a Neovim `:terminal`), tag is your best friend.

Finally, tag is unobtrusive. It should behave exactly like `ag` or `ripgrep` under most circumstances.

## Performance Benchmarks

tag processes ag's output on-the-fly with Golang using pipes so the performance loss is neglible. In other words, **tag is just as fast as ag**!

```
$ cd ~/github/torvalds/linux
$ time ( for _ in {1..10}; do  ag EXPORT_SYMBOL_GPL >/dev/null 2>&1; done )
16.66s user 16.54s system 347% cpu 9.562 total
$ time ( for _ in {1..10}; do tag EXPORT_SYMBOL_GPL >/dev/null 2>&1; done )
16.84s user 16.90s system 356% cpu 9.454 total
```

# Installation

1. Update to the latest versions of [`ag`](https://github.com/ggreer/the_silver_searcher) or [`ripgrep`](https://github.com/BurntSushi/ripgrep). `ag` in particular must be version `>= 0.25.0`.

1. Install the `tag` binary using one of the following methods.
    - Homebrew (OSX)
      ```
      $ brew tap aykamko/tag-ag
      $ brew install tag-ag
      ```

    - AUR (Arch Linux)

      Using your favorite [AUR helper](https://wiki.archlinux.org/index.php/AUR_helpers), install the [tag-ag AUR package](https://aur.archlinux.org/packages/tag-ag/) like this:
      ```
      $ aura -A tag-ag
      ```

    - [Download a compressed binary for your platform](https://github.com/aykamko/tag/releases)

    - Developers and other platforms
      ```
      $ go get -u github.com/aykamko/tag/...
      $ go install github.com/aykamko/tag
      ```

1. By default, `tag` uses `ag` as its search backend. To use `ripgrep` instead, set the environment variable `TAG_SEARCH_PROG=rg`.
> To persist this setting, put it in your `bashrc`/`zshrc`. 

1. Since tag generates a file with command aliases for your shell, you'll have to drop the following in your `bashrc`/`zshrc` to actually pick up those aliases.
    - `bash`
      ```bash
      if hash ag 2>/dev/null; then
        tag() { command tag "$@"; source ${TAG_ALIAS_FILE:-/tmp/tag_aliases} 2>/dev/null; }
        alias ag=tag  # replace with rg for ripgrep
      fi
      ```

    - `zsh`
      ```zsh
      if (( $+commands[tag] )); then
        tag() { command tag "$@"; source ${TAG_ALIAS_FILE:-/tmp/tag_aliases} 2>/dev/null }
        alias ag=tag  # replace with rg for ripgrep
      fi
      ```

    - `fish - ~/.config/fish/functions/tag.fish`
      ```fish
      function tag
          set -q TAG_ALIAS_FILE; or set -l TAG_ALIAS_FILE /tmp/tag_aliases
          command tag $argv; and source $TAG_ALIAS_FILE ^/dev/null
          alias ag tag  # replace with rg for ripgrep
      end
      ```

# Configuration

`tag` exposes the following configuration options via environment variables:

- `TAG_SEARCH_PROG`
  - Determines whether to use `ag` or `ripgrep` as the search backend. Must be one of `ag` or `rg`.
  - Default: `ag`
- `TAG_ALIAS_FILE`
  - Path where shortcut alias file will be generated.
  - Default: `/tmp/tag_aliases`
- `TAG_ALIAS_PREFIX`
  - Prefix for alias commands, e.g. the `e` in generated alias `e42`.
  - Default: `e`
- `TAG_CMD_FMT_STRING`
  - Format string for alias commands. Must contain `{{.Filename}}`, `{{.LineNumber}}`, and `{{.ColumnNumber}}` for proper substitution.
  - Default: `vim -c 'call cursor({{.LineNumber}}, {{.ColumnNumber}})' '{{.Filename}}'`

# License

[MIT](LICENSE)

# Author

[aykamko](https://github.com/aykamko)
