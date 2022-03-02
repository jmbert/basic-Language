package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func isWhiteSpace(char byte) bool {
	return char == ' ' || char == '\t' || char == '\n'
}

func isLetter(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}

func isDigit(char byte) bool {
	return (char >= '0' && char <= '9')
}

type Token int

const (

	// Invalid token
	INVTOK Token = iota

	// Binary operators
	ADD
	SUB
	MUL
	DIV

	// Misc
	WS
	LBRACKET
	RBRACKET

	// Variable stuff
	ASSIGN
	REFERENCE
	DECLARE

	// Literals
	INT
	FLOAT

	// Identifiers
	IDENT
	PRINT

	// End of file
	EOF
)

var dispatch map[string]Token = make(map[string]Token)

func setUpDispatch() {
	dispatch["EOF"] = EOF
	dispatch["INVTOK"] = INVTOK
	dispatch["ADD"] = ADD
	dispatch["SUB"] = SUB
	dispatch["MUL"] = MUL
	dispatch["DIV"] = DIV
	dispatch["PRINT"] = PRINT
	dispatch["ASSIGN"] = ASSIGN
	dispatch["REFERENCE"] = REFERENCE
	dispatch["DECLARE"] = DECLARE
	dispatch["WS"] = WS
	dispatch["LBRACKET"] = LBRACKET
	dispatch["RBRACKET"] = RBRACKET
	dispatch["INT"] = INT
	dispatch["FLOAT"] = FLOAT
}

type Lexer struct {
	reader      *bufio.Reader
	char        byte
	positionIdx int
	tokens      []string
	buffer      []byte
}

func (l *Lexer) advance() {
	char, err := l.reader.ReadByte()
	l.positionIdx++
	if err != nil {
		l.char = 0
	} else {
		l.char = char
	}
}

func (l *Lexer) unAdvance() {
	_ = l.reader.UnreadByte()
	if len(l.buffer) != 0 {
		l.buffer = l.buffer[:len(l.buffer)-1]
	}
	l.positionIdx--
}

func (l *Lexer) scanThroughWhiteSpace() {
	for {
		l.advance()
		if l.char == 0 {
			l.tokens = append(l.tokens, "WS")
			l.tokens = append(l.tokens, "EOF")
			break
		}
		if !isWhiteSpace(l.char) {
			l.unAdvance()
			l.tokens = append(l.tokens, "WS")
			break
		}
	}

}

func (l *Lexer) scanDigit() {
	dots := 0

	for {
		l.advance()
		l.buffer = append(l.buffer, l.char)

		if l.char == '.' && dots == 0 {
			dots++
		}
		if l.char == 0 {
			l.unAdvance()
			if dots == 0 {
				l.tokens = append(l.tokens, "INT:"+string(l.buffer))
			} else {
				l.tokens = append(l.tokens, "FLOAT:"+string(l.buffer))
			}
			l.tokens = append(l.tokens, "EOF")
			break
		}
		if !isDigit(l.char) && l.char != '.' {
			if len(l.buffer) != 0 {
				l.buffer = l.buffer[:len(l.buffer)-1]
			}
			if dots == 0 {
				l.tokens = append(l.tokens, "INT:"+string(l.buffer))
			} else {
				l.tokens = append(l.tokens, "FLOAT:"+string(l.buffer))
			}
			if isWhiteSpace(l.char) {
				l.buffer = make([]byte, 0)
				l.scanThroughWhiteSpace()
				if l.char == 0 {
					break
				}
			}
			break
		}
	}

}

func (l *Lexer) parseTokName() string {

	switch string(l.buffer) {
	case "EOF":
		return "EOF"
	case "INVTOK":
		return "INVTOK"
	case "+":
		return "ADD"
	case "-":
		return "SUB"
	case "*":
		return "MUL"
	case "PRINT":
		return "PRINT"
	case "=":
		return "ASSIGN"
	case " ":
		return "WS"
	case "(":
		return "LBRACKET"
	case ")":
		return "RBRACKET"
	default:
		return string(l.buffer)
	}
}

func (l *Lexer) appendIdentifier() {

	if l.buffer[0] == '$' {
		varName := l.buffer[1:len(l.buffer)]
		l.tokens = append(l.tokens, "REFERENCE:"+string(varName))
	} else if l.buffer[0] == '@' {
		varName := l.buffer[1:len(l.buffer)]
		l.tokens = append(l.tokens, "DECLARE:"+string(varName))
	} else {
		buf := l.parseTokName()

		if _, ok := dispatch[buf]; ok {
			l.tokens = append(l.tokens, buf)
		} else {
			fmt.Printf("%s\n", buf)
			l.tokens = append(l.tokens, "INVTOK")
		}
	}

}

func (l *Lexer) scan() {
	for {
		l.advance()
		if isWhiteSpace(l.char) {
			if len(l.buffer) != 0 {
				l.appendIdentifier()
			}
			l.buffer = make([]byte, 0)
			l.scanThroughWhiteSpace()
			if l.char == 0 {
				break
			}
		} else if isDigit(l.char) {
			l.unAdvance()
			if len(l.buffer) != 0 {
				l.appendIdentifier()
			}
			l.buffer = make([]byte, 0)
			l.scanDigit()
			if l.char == 0 {
				break
			}
			l.buffer = make([]byte, 0)
		} else {
			l.buffer = append(l.buffer, l.char)
		}

		if l.char == 0 {
			if len(l.buffer) != 0 {
				l.unAdvance()
				l.appendIdentifier()
			}
			l.tokens = append(l.tokens, "EOF")
			break
		}
	}
}

func NewLexer(r *strings.Reader) Lexer {
	l := Lexer{reader: bufio.NewReader(r), char: 0, positionIdx: -1, tokens: make([]string, 0)}
	return l
}

func (l *Lexer) String() string {
	var out = l.tokens[0]
	for i := 1; i < len(l.tokens); i++ {
		out += ", "
		out += l.tokens[i]
	}
	return out
}

func main() {
	fText, err := ioutil.ReadFile("code.txt")

	if err != nil {
		log.Fatal("Invalid File")
	}

	setUpDispatch()

	l := NewLexer(strings.NewReader(string(fText)))
	l.scan()

	fmt.Printf("%s\n", l.String())
}
