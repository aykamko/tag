tag - Tag your ag matches
====
![revolv++](tag.gif)

**tag** is a lightweight wrapper around **[ag](https://github.com/ggreer/the_silver_searcher)** that generates shell aliases for **ag** matches. tag is a very fast Golang rewrite of [sack](https://github.com/sampson-chen/sack).

tag only supports ag. There are no plans to support ack or grep. Support for pt may be added if users show interest.

## Why should I use tag?

tag makes it easy to _immediately_ jump to an ag match in your favorite editor. It eliminates the tedious task of typing `vim foo/bar/baz.qux +42` to jump to a match by automatically generating these commands for you as shell aliases.

Inside vim, [vim-grepper](https://github.com/mhinz/vim-grepper) or [ag.vim](https://github.com/rking/ag.vim) is probably the way to go. Outside vim (or inside a Neovim `:terminal`), tag is your best friend.

Finally, tag is unobtrusive. It should behave exactly like ag under most circumstances.

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

1. Install the `tag` binary using one of the following methods.
    - Homebrew (OSX)
      ```
      $ brew tap aykamko/tag-ag
      $ brew install tag-ag
      ```

    - [Download a compressed binary for your platform](https://github.com/aykamko/tag/releases)

    - Developers and other platforms
      ```
      $ go get -u github.com/aykamko/tag/...
      $ go install github.com/aykamko/tag
      ```

1. Since tag generates a file with command aliases for your shell, you'll have to drop the following in your `bashrc`/`zshrc` to actually pick up those aliases.
    - `bash`
      ```bash
      if hash ag 2>/dev/null; then
        tag() { command tag "$@"; source ${TAG_ALIAS_FILE:-/tmp/tag_aliases} 2>/dev/null; }
        alias ag=tag
      fi
      ```

    - `zsh`
      ```zsh
      if (( $+commands[tag] )); then
        tag() { command tag "$@"; source ${TAG_ALIAS_FILE:-/tmp/tag_aliases} 2>/dev/null }
        alias ag=tag
      fi
      ```

    - `fish - ~/.config/fish/functions/tag.fish`
      ```fish
      function tag
          set -q TAG_ALIAS_FILE; or set -l TAG_ALIAS_FILE /tmp/tag_aliases
          command tag $argv; and source $TAG_ALIAS_FILE ^/dev/null
          alias ag tag
      end
      ```

# Configuration

tag exposes the following configuration options via environment variables:

- `TAG_ALIAS_FILE`
  - Path where shortcut alias file will be generated.
  - Default: `/tmp/tag_aliases`
- `TAG_ALIAS_PREFIX`
  - Prefix for alias commands, e.g. the `e` in generated alias `e42`.
  - Default: `e`
- `TAG_CMD_FMT_STRING`
  - Format string for alias commands. Must contain `{{.Filename}}` and `{{.LineNumber}}` for proper substitution.
  - Default: `vim {{.Filename}} +{{.LineNumber}}`

# License

[MIT](LICENSE)

# Author

[aykamko](https://github.com/aykamko)
