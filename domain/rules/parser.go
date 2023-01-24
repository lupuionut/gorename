package rules

import (
    "strings"
)

type TokenType int
const (
    TokenReplace TokenType = iota
    TokenWith
    TokenTagStart
    TokenTagEnd
    TokenTagValue
    TokenEOL
    TokenUnknown
)

type Token struct {
    Type TokenType
    Value string
}

type Parser struct {
    Tokens [][]*Token
    Content []string
    Line int
    Cursor int
    LastChar rune
    Buffer string
}

func (parser *Parser) FilterValidLines() {
   var newcontent []string
   for _, line := range(parser.Content) {
        maxc := len(line)
        text := strings.Trim(line, " ")
        if maxc == 0 || text[0] == '#' {
            continue
        }
        newcontent = append(newcontent, line)
   }
   parser.Content = newcontent
}

func (parser *Parser) Parse() error {
    parser.FilterValidLines()
    parser.Tokens = make([][]*Token, 2)
    for i := range(parser.Content) {
       	parser.Line = i
        err := parser.ParseLine()
        if err != nil {
            return err
        }
    }
    return nil
}

func (parser *Parser) ParseLine() error {
    max := len(parser.Content[parser.Line])

    if parser.Cursor >= max {
        token := &Token {
            Type: TokenEOL,
            Value: "",
        }
        parser.Tokens[parser.Line] = append(parser.Tokens[parser.Line], token)
        parser.Cursor = 0
        return nil
    }

    if parser.Content[parser.Line][parser.Cursor] == ' ' {
        parser.Cursor++
        return parser.ParseLine()
    }

    if parser.Content[parser.Line][parser.Cursor] == '<' {
        token := &Token {
            Type: TokenTagStart,
            Value: "<",
        }
        parser.Tokens[parser.Line] = append(parser.Tokens[parser.Line], token)
        parser.Cursor++
        parser.ConsumeTag(max)
        return parser.ParseLine()
    }

    if parser.Content[parser.Line][parser.Cursor] == '>' {
        token := &Token {
            Type: TokenTagEnd,
            Value: ">",
        }
        parser.Tokens[parser.Line] = append(parser.Tokens[parser.Line], token)
        parser.Cursor++
        return parser.ParseLine()
    }

    for i := parser.Cursor; i < max; i++  {
        current := parser.Content[parser.Line][i]
        parser.Cursor++
        if current == ' ' || parser.Cursor == max-1 {
            if parser.Buffer == "replace" {
                token := &Token {
                    Type: TokenReplace,
                    Value: "replace",
                }
                parser.Tokens[parser.Line] = append(parser.Tokens[parser.Line], token)
            } else if  parser.Buffer == "with" {
                token := &Token {
                    Type: TokenWith,
                    Value: "with",
                }
                parser.Tokens[parser.Line] = append(parser.Tokens[parser.Line], token)
            } else {
                token := &Token {
                    Type: TokenUnknown,
                    Value: parser.Buffer,
                }
                parser.Tokens[parser.Line] = append(parser.Tokens[parser.Line], token)
            }
            parser.Buffer = ""
            return parser.ParseLine()
        }
        parser.Buffer += string(current)
    }

    return parser.ParseLine()
}

func (parser *Parser) ConsumeTag(max int) {
    if parser.Cursor == max {
        return
    }
    var content string
    for i := parser.Cursor; i < max; i++ {
        current := parser.Content[parser.Line][i]
        if current == '>' && parser.LastChar != '\\' {
            token := &Token {
                Type: TokenTagValue,
                Value: content,
            }
            parser.Tokens[parser.Line] = append(parser.Tokens[parser.Line], token)
            return
        }
        if current == '\\' {
            if parser.LastChar == '\\' {
                var c rune
                parser.LastChar = c
                content += "\\"
                parser.Cursor++
            } else {
                parser.LastChar = rune(current)
                parser.Cursor++
                continue
            }
        } else {
            parser.LastChar = rune(current)
            content += string(current)
            parser.Cursor++
        }
    }
}
