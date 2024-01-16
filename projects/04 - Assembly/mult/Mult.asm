// RAM[100] = RAM[200]
// @200
// D=M
// @100
// M=D
//
// @N creates a constant used as parameter when accessing M or D, @20 turns A=20, @var turn M=var
// Set var x=17: @17, D=A, @x, M=D
//@Rn -> direct access RAM[n], n: 0-15

@R2
M=0

@R0
D=M
@STEP
D;JGT

@END
0;JMP

(STEP)
    @R2
    D=M

    @R1
    D=D+M

    @R2
    M=D

    @R0
    D=M-1
    M=D

    @STEP
    D;JGT

(END)
    @END
    0;JMP