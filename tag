#!/bin/bash

os=$(uname)
tmp_file=/tmp/tag_aliases
RED=$(tput setaf 1)
BLUE=$(tput setaf 4)
CLEAR=$(tput sgr0)

function strip_ansi {
    # http://www.commandlinefu.com/commands/view/3584/remove-color-codes-special-characters-with-sed
    echo
    # echo $@ | if [[ "${os}" == 'Darwin' ]]; then
    #     sed -E "s/"$'\E'"\[([0-9]{1,2}(;[0-9]{1,2})*)?[a-zA-Z]//g"
    # else
    #     sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})*)?[a-zA-Z]//g"
    # fi
}

function parse_output {
    index=1
    curfile=""
    command ag --group --color $@ | while read line; do
        if [ -z "$line" ]; then
            echo && continue
        elif [[ $line != *":"* ]]; then
            curfile=$(strip_ansi $line)
            echo $line && continue
        else
            fpath=($(strip_ansi "$line" | tr ':' '\n'))
            echo "alias \"f${index}\"=\"vim $PWD/$curfile +${fpath[0]}\"" >> $tmp_file
            echo "${BLUE}[${RED}$index${BLUE}]${CLEAR} $line"
            let "index += 1"
        fi
    done
}

[ -f $tmp_file ] && rm -rf $tmp_file; touch $tmp_file
parse_output "$@"
