#!/bin/bash

red=$(tput setaf 1)
green=$(tput setaf 2)
yellow=$(tput setaf 3)
cyan=$(tput setaf 6)
purple=$(tput setaf 177)
reset=$(tput sgr0)

p-debug()    { printf "${purple}%s${reset}\n" "$@"; }
p-info()     { printf "${cyan}%s${reset}\n" "$@"; }
p-warn()     { printf "${yellow}%s${reset}\n" "$@"; }
p-error()    { printf "${red}%s${reset}\n" "$@"; }
p-success()  { printf "${green}%s${reset}\n" "$@"; }

p-debugf()   { printf "${purple}$*${reset}" ; }
p-infof()    { printf "${cyan}$*${reset}" ; }
p-warnf()    { printf "${yellow}$*${reset}" ; }
p-errorf()   { printf "${red}$*${reset}" ; }
p-successf() { printf "${green}$*${reset}" ; }

show-help(){
    items="$(compgen -A "function" | grep "run-" | sed "s/run-//")"

    tasks=()
    while IFS= read -r task; do tasks+=("$task"); done<<<"$items"

    cmd="$(basename "$0")"
    usage=""
    case $1 in
    "--helpz")
        usage+="\nMore tasks:"
    ;;
    *)
        usage+="Usage: $cmd TASK [ARGUMENTS]\n\nTasks:"
    ;;
    esac
    p-infof "$usage\n"
    for task in "${tasks[@]}"; do p-infof "\t$task\n"; done
}

task=$1
case "$task" in
    "" | "-h" | "--help" | "--helpz")
        show-help $1
        if [[ "$1" && "$PWD" != "$REPO_ROOT" ]]; then
            bash -c "cd $REPO_ROOT && $REPO_ROOT/run --helpz"
        fi
        ;;
    *)
        shift
        if compgen -A "function" | grep "run-$task" >/dev/null ; then
            run-"${task}" "$@"
        elif [[ "$PWD" != "$REPO_ROOT" ]]; then
            bash -c "cd $REPO_ROOT && $REPO_ROOT/run $task $*"
        else
            echo "Task $task not found."
            exit 2
        fi
        ;;
esac
