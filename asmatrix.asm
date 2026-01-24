global main

extern initscr
extern cbreak
extern noecho
extern printw
extern refresh
extern getch
extern endwin
extern keypad
extern stdscr

section .data
    msg db "ass", 10, 0

section .text
main:
    sub rsp, 8                  ; stack alignment

    call initscr
    call cbreak
    call noecho

    ; keypad(stdscr, TRUE)
    mov rdi, [rel stdscr]       ; WINDOW *win
    mov esi, 1                  ; TRUE
    call keypad

    lea rdi, [rel msg]
    xor eax, eax                ; variadic ABI rule
    call printw

    call refresh
    call getch                  ; blocks
    call endwin

    add rsp, 8
    xor eax, eax
    ret
