#!/bin/bash

red=$(tput setaf 1)
green=$(tput setaf 2)
yellow=$(tput setaf 3)
cyan=$(tput setaf 6)
purple=$(tput setaf 177)
reset=$(tput sgr0)

p_debug()    { printf "${purple}%s${reset}\n" "$@"; }
p_info()     { printf "${cyan}%s${reset}\n" "$@"; }
p_warn()     { printf "${yellow}%s${reset}\n" "$@"; }
p_error()    { printf "${red}%s${reset}\n" "$@"; }
p_success()  { printf "${green}%s${reset}\n" "$@"; }

p_debugf()   { printf "${purple}$*${reset}" ; }
p_infof()    { printf "${cyan}$*${reset}" ; }
p_warnf()    { printf "${yellow}$*${reset}" ; }
p_errorf()   { printf "${red}$*${reset}" ; }
p_successf() { printf "${green}$*${reset}" ; }

task_help(){
    task_list="$(compgen -A "function" | grep "task_" | grep -v "task_help" | sed "s/task_//")"

    tasks=()
    while IFS= read -r task; do tasks+=("$task"); done<<<"$task_list"

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
    p_infof "$usage\n"
    for task in "${tasks[@]}"; do p_infof "\t$task\n"; done
}

task=$1
case "$task" in
    "" | "-h" | "--help" | "--helpz")
        task_help $1
        if [[ "$1" && "$PWD" != "$REPO_ROOT" ]]; then
            bash -c "cd $REPO_ROOT && $REPO_ROOT/run --helpz"
        fi
        ;;
    *)
        shift
        if compgen -A "function" | grep "task_$task" >/dev/null ; then
            task_"${task}" "$@"
        elif [[ "$PWD" != "$REPO_ROOT" ]]; then
            bash -c "cd $REPO_ROOT && $REPO_ROOT/run $task $*"
        else
            echo "Task $task not found."
            exit 2
        fi
        ;;
esac
