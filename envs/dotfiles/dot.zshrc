export LESSCHARSET=utf-8
export EDITOR=vim

autoload -Uz colors
colors
autoload -Uz compinit
compinit -u
setopt auto_pushd
setopt pushd_ignore_dups
setopt correct
setopt extended_glob
setopt print_eight_bit
setopt auto_cd
setopt no_beep
setopt nolistbeep

WORDCHARS='*?_-.[]~=&;!#$%^(){}<>'
# --- history ---
HISTFILE=~/.share/.zsh_history
HISTSIZE=100000
SAVEHIST=100000
HISTORY_IGNORE="(ls|cd|pwd|exit)"

setopt share_history
setopt hist_ignore_space
setopt hist_no_store
setopt hist_ignore_dups
setopt hist_verify
setopt extended_history
setopt hist_reduce_blanks
setopt hist_save_no_dups
setopt hist_ignore_all_dups
setopt hist_expand

bindkey -e
zstyle ':completion:*' matcher-list 'm:{a-z}={A-Z}'
zstyle ':completion:*:default' list-colors ${(s.:.)LS_COLORS}

alias ls='ls --color=auto'
alias grep='grep --color=auto'
alias ll='ls -alF'
alias reload='source ~/.zshrc'

autoload -Uz vcs_info
setopt prompt_subst
zstyle ':vcs_info:git:*' check-for-changes true
zstyle ':vcs_info:git:*' stagedstr "%F{magenta}!"
zstyle ':vcs_info:git:*' unstagedstr "%F{yellow}+"
zstyle ':vcs_info:*' formats " - %F{cyan}%c%u[%b]%f"
zstyle ':vcs_info:*' actionformats '[%b|%a]'
precmd () { vcs_info }

# プロンプトカスタマイズ
PROMPT='%F{red}[%n]%f - %F{green}[%~]%f%F{cyan}$vcs_info_msg_0_%f
%F{white}$%f '
RPROMPT="%F{white}[%D{%Y-%m-%d %H:%M:%S}]%f"
