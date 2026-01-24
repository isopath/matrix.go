global main

; imports
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
    ; msg = ["a", "s", "s", "\n" (10), "\0" (0)]
    msg db "ass", 10, 0

section .text
default rel

main:
    ; stack alignment
    sub rsp, 8

    ; initialise the screen
    call initscr
    call cbreak
    call noecho

    ; keypad (stdscr, TRUE)
%ifdef MACOS
    mov rdi, [stdscr wrt ..gotpcrel]
    mov rdi, [rdi]
%else
    mov rdi, [stdscr]
%endif
    mov esi, 1
    call keypad

    ; load address of msg
    lea rdi, [msg]

    ; since printw is a variadic function,
    ; we must declare number of registers used.
    ; xor eax, eax := 0 since we are using
    ; printw with 0 arguments
    xor eax, eax
    call printw
    call refresh
    call getch
    call endwin

    ; restore stack pointer
    add rsp, 8
    xor eax, eax
    ret
