package rules

import (
    "strings"
    "unicode/utf8"
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
    max := utf8.RuneCount([]byte(parser.Content[parser.Line]))

    if parser.Cursor > max {
        token := &Token {
            Type: TokenEOL,
            Value: "",
        }
        parser.Tokens[parser.Line] = append(parser.Tokens[parser.Line], token)
        parser.Cursor = 0
        return nil
    }

    current, size := utf8.DecodeRuneInString(parser.Content[parser.Line][parser.Cursor:])

    if current == ' ' {
        parser.Cursor += size
        return parser.ParseLine()
    }

    if current == '<' {
        token := &Token {
            Type: TokenTagStart,
            Value: "<",
        }
        parser.Tokens[parser.Line] = append(parser.Tokens[parser.Line], token)
        parser.Cursor += size
        parser.ConsumeTag(max)
        return parser.ParseLine()
    }

    if current == '>' {
        token := &Token {
            Type: TokenTagEnd,
            Value: ">",
        }
        parser.Tokens[parser.Line] = append(parser.Tokens[parser.Line], token)
        parser.Cursor += size
        return parser.ParseLine()
    }

    for i := parser.Cursor; i < max; i += size  {
        current, size = utf8.DecodeRuneInString(parser.Content[parser.Line][parser.Cursor:])
        parser.Cursor += size
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
        token := &Token {
            Type: TokenTagValue,
            Value: parser.Buffer,
        }
        parser.Buffer = ""
        parser.Tokens[parser.Line] = append(parser.Tokens[parser.Line], token)
        return
    }
    current, size := utf8.DecodeRuneInString(parser.Content[parser.Line][parser.Cursor:])
    if current == '>' && parser.LastChar != '\\' {
        token := &Token {
            Type: TokenTagValue,
            Value: parser.Buffer,
        }
        parser.Buffer = ""
        parser.Tokens[parser.Line] = append(parser.Tokens[parser.Line], token)
        return
    }
    if current == '\\' {
        if parser.LastChar == '\\' {
            var c rune
            parser.LastChar = c
            parser.Buffer += "\\"
            parser.Cursor += size
            parser.ConsumeTag(max)
        } else {
            parser.LastChar = rune(current)
            parser.Cursor += size
            parser.ConsumeTag(max)
        }
    } else {
        parser.LastChar = current
        parser.Buffer += string(current)
        parser.Cursor += size
        parser.ConsumeTag(max)
    }
}
