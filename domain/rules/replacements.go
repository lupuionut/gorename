package rules

func GenerateReplacements(line []*Token) (key string, value string) {
    var k string
    var v string
    for _, t := range(line) {
        if t.Type == TokenTagValue {
            if len(k) == 0 {
                k = t.Value
            } else{
                v = t.Value
            }
        }
    }
    return k, v
}
