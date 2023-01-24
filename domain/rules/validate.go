package rules

func IsValid(line []*Token) bool {
    valid := [9]TokenType{TokenReplace, TokenTagStart, TokenTagValue, TokenTagEnd,
                        TokenWith, TokenTagStart, TokenTagValue, TokenTagEnd, TokenEOL}
    for i, rule := range(line) {
        if rule.Type != valid[i] {
            return false
        }
    }
    return true
}
