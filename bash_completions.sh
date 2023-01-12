_gobrew()
{
    COMP_WORDBREAKS=${COMP_WORDBREAKS//:}
    case "$COMP_CWORD" in
        1)
            COMMANDS="ls ls-remote use install uninstall prune version help"
            ;;
        2)
            COMMANDS=`/Users/pulkit.kathuria/git/gobrew/main ls |sed '/*/d'| sed '/current/d' |awk NF`
            ;;
    esac
    COMPREPLY=(`compgen -W "$COMMANDS" -- "${COMP_WORDS[COMP_CWORD]}"`)
    return 0
}

complete -F _gobrew gobrew