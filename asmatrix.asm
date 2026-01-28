global main

; -------------------------------------------------
; imports
; -------------------------------------------------

; ncurses
extern initscr
extern cbreak
extern noecho
extern printw
extern refresh
extern getch
extern endwin
extern keypad
extern stdscr

; libc
extern fopen
extern fread
extern fclose
extern strerror
%ifdef MACOS
extern __error
%else
extern __errno_location
%endif

section .data
    rmode   db "r", 0
    err_msg db "Error: %s", 10, 0
    usage   db "Usage: %s <filename>", 10, 0

section .bss
    ; 4KB buffer
    buffer  resb 4096

section .text
default rel

; -------------------------------------------------
; Utility function: print string to screen
; rdi = string address
; -------------------------------------------------
print_screen:
    ; since printw is a variadic function,
    ; we must declare number of registers used.
    ; xor eax, eax := 0 since we are using
    ; printw with 0 arguments
    xor eax, eax
    jmp printw

; -------------------------------------------------
; Main entry point: program initialization and execution
; rdi = argc (argument count)
; rsi = argv (argument vector)
; -------------------------------------------------
main:
    ; save argc/argv
    push rbp
    mov rbp, rsp

    ; Save argc and argv in callee-saved registers
    mov r12, rdi            ; r12 = argc
    mov r13, rsi            ; r13 = argv

    ; stack alignment (16-byte align before calls)
    sub rsp, 8
    push r12
    push r13

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

    ; -------------------------------------------------
    ; Check arguments
    ; -------------------------------------------------
    ; argc == 2?
    cmp r12, 2
    jne .show_usage

    ; -------------------------------------------------
    ; Open File: fopen(argv[1], "r")
    ; -------------------------------------------------
    mov rdi, [r13 + 8]          ; rdi = argv[1]
    lea rsi, [rmode]            ; rsi = "r"
    call fopen
    test rax, rax               ; Check for NULL
    jz .show_strerror           ; If null, display error
    mov rbx, rax                ; Save FILE* in rbx

    ; -------------------------------------------------
    ; Read File: fread(buffer, 1, 4096, file)
    ; -------------------------------------------------
    lea rdi, [buffer]
    mov rsi, 1
    mov rdx, 4096
    mov rcx, rbx
    call fread

    ; Null-terminate the string
    mov rcx, rax                ; Save bytes read
    lea rdi, [buffer]           ; Load buffer address
    mov byte [rdi + rcx], 0     ; Null terminate using saved count

    ; Close file
    mov rdi, rbx
    call fclose

    ; -------------------------------------------------
    ; Print content and exit
    ; -------------------------------------------------
    lea rdi, [buffer]
    call print_screen
    jmp .cleanup

; -------------------------------------------------
; Utility function: display usage message
; Shows program usage instructions when arguments are invalid
; -------------------------------------------------
.show_usage:
    lea rdi, [usage]
    mov rsi, [r13]              ; argv[0]
    xor eax, eax
    call printw
    jmp .cleanup

; -------------------------------------------------
; Utility function: display system error message
; Retrieves errno value and displays corresponding error string
; -------------------------------------------------
.show_strerror:
%ifdef MACOS
    call __error              ; rax = &errno
%else
    call __errno_location     ; rax = &errno
%endif
    mov edi, [rax]            ; edi = errno value
    call strerror             ; rax = error string
    lea rdi, [err_msg]
    mov rsi, rax
    xor eax, eax
    call printw

; -------------------------------------------------
; Utility function: cleanup and exit
; Handles ncurses cleanup, waits for user input, and exits
; Restores stack and returns with exit code 0
; -------------------------------------------------
.cleanup:
    call refresh
    call getch
    call endwin
    xor eax, eax
    add rsp, 24               ; Clean up stack (8 + 16 for pushed registers)
    pop rbp
    ret
