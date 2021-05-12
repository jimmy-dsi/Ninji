package lexer

import (
	"unicode/utf8"
	"fmt"
	"strconv"
	"strings"
	. "io/ioutil"
)

type Lexer struct {
	FilePath string
	FileData string

	Line    int
	Column  int

	StateID int
}

type Token struct {
	ID       string
	RawValue [] byte
	Value    interface {}

	FilePath string
	Line     int
	Column   int

	Error    bool
}

func (this Lexer) Init(filePath string) Lexer {
	fileData, err := ReadFile(filePath)

	if err == nil {
		return Lexer {
			FilePath: filePath,
			FileData: string(fileData) + "\n\n",

			Line:     1,
			Column:   1,
			StateID:  0,
		}
	} else {
		return Lexer {}
	}
}

func (this *Lexer) Lex() [] Token {
	tokens      := [] Token {}
	finalTokens := [] Token {}

	tempToken := Token {}

	tokenize := func() {
		utf8_counter := 0

		for i := 0; i < len(this.FileData); i++ {
			var c    byte = '\n'
			var cc   byte = '\n'
			var ccc  byte = '\n'
			var cccc byte = '\n'
	
			if i < len(this.FileData) {
				c = this.FileData[i]
			}
			if i < len(this.FileData) - 1 {
				cc = this.FileData[i + 1]
			}
			if i < len(this.FileData) - 2 {
				ccc = this.FileData[i + 2]
			}
			if i < len(this.FileData) - 3 {
				cccc = this.FileData[i + 3]
			}

			if utf8_counter == 0 {
				if c & 0b11100000 ^ 0b11000000 == 0 && cc & 0b11000000 ^ 0b10000000 == 0 {
					utf8_counter = 1
				} else if c & 0b11110000 ^ 0b11100000 == 0 && cc & 0b11000000 ^ 0b10000000 == 0 && ccc & 0b11000000 ^ 0b10000000 == 0 {
					utf8_counter = 2
				} else if c & 0b11111000 ^ 0b11110000 == 0 && cc & 0b11000000 ^ 0b10000000 == 0 && ccc & 0b11000000 ^ 0b10000000 == 0 && cccc & 0b11000000 ^ 0b10000000 == 0 {
					utf8_counter = 3
				}
			} else {
				utf8_counter -= 1
			}
	
			new_token, is_new := this.ReadChar(c, cc, ccc, &tempToken)
			if is_new && new_token.ID != "whitespace" && new_token.ID != "comment" && new_token.ID != "line end" {
				tokens = append(tokens, new_token)
			}
	
			if c == '\n' {
				this.Line++
				this.Column = 1
			} else if c == '\r' {
				// Do not advance line/column
			} else if utf8_counter == 0 {
				this.Column++
			}
		}
	}

	consolidateUnicodeCharacters := func() {
		for i := 0; i < len(tokens); i++ {
			token := tokens[i]

			if len(finalTokens) > 0 {
				prevToken := finalTokens[len(finalTokens) - 1]

				if token.ID == prevToken.ID && token.Line == prevToken.Line && token.Column == prevToken.Column {
					finalTokens[len(finalTokens) - 1].Value    = finalTokens[len(finalTokens) - 1].Value.(string) + token.Value.(string)
					finalTokens[len(finalTokens) - 1].RawValue = []byte(finalTokens[len(finalTokens) - 1].Value.(string))
				} else {
					finalTokens = append(finalTokens, token)
				}
			} else {
				finalTokens = append(finalTokens, token)
			}
		}
	}

	//condense_line_ends := func() {
	//	for i := 0; i < len(tokens); i++ {
	//		token := tokens[i]
	//		if token.ID != "line end" || len(final_tokens) == 0 || final_tokens[len(final_tokens) - 1].ID != "line end" {
	//			final_tokens = append(final_tokens, token)
	//		}
	//	}
	//}

	tokenize()
	consolidateUnicodeCharacters()
	//condenseLineEnds()

	return finalTokens
	//return tokens
}

const (
	Start                      = iota
	AcceptWhiteSpace           = iota
	AcceptComment              = iota
	AcceptLineEnd              = iota
     
	AcceptIdent                = iota
	AcceptKeyword              = iota
	AcceptInteger              = iota
	AcceptDecimal              = iota
	AcceptFloat                = iota
	AcceptDouble               = iota
	AcceptFixed                = iota
	AcceptByte                 = iota
	AcceptWord                 = iota
	AcceptLong                 = iota
	AcceptDWord                = iota
	AcceptQWord                = iota
	AcceptOperator             = iota
     
	AcceptCStringStart         = iota
	AcceptWStringPrefix        = iota
	AcceptUStringPrefix        = iota
	AcceptStringChar           = iota
	AcceptStringEnd            = iota
     
	AcceptChar                 = iota
	AcceptUCharPrefix          = iota
	AcceptWCharPrefix          = iota
     
	RejectChar                 = iota
      
	RejectLiteral              = iota
	RejectNumber               = iota
	RejectCharacter            = iota
	RejectStringChar           = iota
	RejectStringEnd            = iota
      
	InvalidReadNumber          = iota
	    
	InvalidReadString          = iota
	InvalidReadBackslashU1     = iota
	InvalidReadBackslashU2     = iota
	InvalidReadBackslashU3     = iota
	InvalidReadBackslashX1     = iota
	
	InvalidReadChar            = iota
  
	ReadWhitespace             = iota
	ReadComment                = iota
	ReadMultiLineComment       = iota
	ReadMultiLineComment2      = iota
     
	ReadIdent                  = iota
	     
	ReadString                 = iota
	ReadStringBackslash        = iota
	ReadStringBackslashX       = iota
	ReadStringBackslashX1      = iota
	ReadStringBackslashU       = iota
	ReadStringBackslashU1      = iota
	ReadStringBackslashU2      = iota
	ReadStringBackslashU3      = iota
	     
	ReadChar                   = iota
	ReadCharFinal              = iota
	ReadCharBackslash          = iota
	ReadCharBackslashX         = iota
	ReadCharBackslashX1        = iota
	ReadCharBackslashU         = iota
	ReadCharBackslashU1        = iota
	ReadCharBackslashU2        = iota
	ReadCharBackslashU3        = iota
      
	ReadDot1                   = iota
	ReadDot2                   = iota
	ReadEquals1                = iota
	ReadEquals2                = iota
	ReadColon                  = iota
	ReadPlus                   = iota
	ReadMinus                  = iota
	ReadStar                   = iota
	ReadSlash                  = iota
	ReadPercent                = iota
	ReadAmp                    = iota
	ReadPipe                   = iota
	ReadCaret                  = iota
	ReadGreaterThan            = iota
	ReadShiftRight             = iota
	ReadLessThan               = iota
	ReadShiftLeft              = iota
	ReadBang                   = iota
	ReadBangEquals             = iota
     
	ReadNumberZero             = iota
	ReadNumberDigit            = iota
	ReadBinaryDigit            = iota
	ReadHexDigit               = iota
	ReadDecimalDigit           = iota
      
	ReadDollarSign             = iota
      
	ReadHexDigit1              = iota
	ReadHexDigit2              = iota
	ReadHexDigit3              = iota
	ReadHexDigit4              = iota
	ReadHexDigit5              = iota
	ReadHexDigit6              = iota
	ReadHexDigit7              = iota
	ReadHexDigit8              = iota
	ReadHexDigitQ              = iota
      
	ReadBinaryDigit0           = iota
      
	ReadBinaryDigit1           = iota
	ReadBinaryDigit2           = iota
	ReadBinaryDigit3           = iota
	ReadBinaryDigit4           = iota
	ReadBinaryDigit5           = iota
	ReadBinaryDigit6           = iota
	ReadBinaryDigit7           = iota
	ReadBinaryDigit8           = iota
      
	ReadBinaryDigit9           = iota
	ReadBinaryDigit10          = iota
	ReadBinaryDigit11          = iota
	ReadBinaryDigit12          = iota
	ReadBinaryDigit13          = iota
	ReadBinaryDigit14          = iota
	ReadBinaryDigit15          = iota
	ReadBinaryDigit16          = iota
      
	ReadBinaryDigit17          = iota
	ReadBinaryDigit18          = iota
	ReadBinaryDigit19          = iota
	ReadBinaryDigit20          = iota
	ReadBinaryDigit21          = iota
	ReadBinaryDigit22          = iota
	ReadBinaryDigit23          = iota
	ReadBinaryDigit24          = iota
      
	ReadBinaryDigit25          = iota
	ReadBinaryDigit26          = iota
	ReadBinaryDigit27          = iota
	ReadBinaryDigit28          = iota
	ReadBinaryDigit29          = iota
	ReadBinaryDigit30          = iota
	ReadBinaryDigit31          = iota
	ReadBinaryDigit32          = iota
      
	ReadBinaryDigitQ           = iota
      
	ReadKeywordsA              = iota
	ReadKeywordsAL             = iota
	ReadKeywordsALI            = iota
	ReadKeywordsALIA           = iota
	ReadKeywordsAN             = iota
	ReadKeywordsB              = iota
	ReadKeywordsBR             = iota
	ReadKeywordsBRE            = iota
	ReadKeywordsBREA           = iota
	ReadKeywordsC              = iota
	ReadKeywordsCA             = iota
	ReadKeywordsCAC            = iota
	ReadKeywordsCACH           = iota
	ReadKeywordsCACHE          = iota
	ReadKeywordsCAS            = iota
	ReadKeywordsCAT            = iota
	ReadKeywordsCATC           = iota
	ReadKeywordsCO             = iota
	ReadKeywordsCON            = iota
	ReadKeywordsCONS           = iota
	ReadKeywordsCONST          = iota
	ReadKeywordsCONSTR         = iota
	ReadKeywordsCONSTRA        = iota
	ReadKeywordsCONSTRAI       = iota
	ReadKeywordsCONT           = iota
	ReadKeywordsCONTI          = iota
	ReadKeywordsCONTIN         = iota
	ReadKeywordsCONTINU        = iota
	ReadKeywordsD              = iota
	ReadKeywordsDE             = iota
	ReadKeywordsDEF            = iota
	ReadKeywordsDEFA           = iota
	ReadKeywordsDEFAU          = iota
	ReadKeywordsDEFAUL         = iota
	ReadKeywordsE              = iota
	ReadKeywordsEA             = iota
	ReadKeywordsEAC            = iota
	ReadKeywordsEL             = iota
	ReadKeywordsELS            = iota
	ReadKeywordsEN             = iota
	ReadKeywordsENU            = iota
	ReadKeywordsEV             = iota
	ReadKeywordsEVA            = iota
	ReadKeywordsF              = iota
	ReadKeywordsFA             = iota
	ReadKeywordsFAL            = iota
	ReadKeywordsFALS           = iota
	ReadKeywordsFI             = iota
	ReadKeywordsFIN            = iota
	ReadKeywordsFINA           = iota
	ReadKeywordsFINAL          = iota
	ReadKeywordsFINALL         = iota
	ReadKeywordsFO             = iota
	ReadKeywordsFU             = iota
	ReadKeywordsFUN            = iota
	ReadKeywordsG              = iota
	ReadKeywordsI              = iota
	ReadKeywordsIM             = iota
	ReadKeywordsIMP            = iota
	ReadKeywordsIMPL           = iota
	ReadKeywordsIMPLE          = iota
	ReadKeywordsIMPLEM         = iota
	ReadKeywordsIMPLEME        = iota
	ReadKeywordsIMPLEMEN       = iota
	ReadKeywordsIMPO           = iota
	ReadKeywordsIMPOR          = iota
	ReadKeywordsIN             = iota
	ReadKeywordsINH            = iota
	ReadKeywordsINHE           = iota
	ReadKeywordsINHER          = iota
	ReadKeywordsINHERI         = iota
	ReadKeywordsINT            = iota
	ReadKeywordsINTE           = iota
	ReadKeywordsINTER          = iota
	ReadKeywordsINTERF         = iota
	ReadKeywordsINTERFA        = iota
	ReadKeywordsINTERFAC       = iota
	ReadKeywordsL              = iota
	ReadKeywordsLE             = iota
	ReadKeywordsM              = iota
	ReadKeywordsMA             = iota
	ReadKeywordsMO             = iota
	ReadKeywordsMOD            = iota
	ReadKeywordsMODU           = iota
	ReadKeywordsMODUL          = iota
	ReadKeywordsN              = iota
	ReadKeywordsNA             = iota
	ReadKeywordsNO             = iota
	ReadKeywordsNOT            = iota
	ReadKeywordsNOTH           = iota
	ReadKeywordsNOTHI          = iota
	ReadKeywordsNOTHIN         = iota
	ReadKeywordsNU             = iota
	ReadKeywordsNUL            = iota
	ReadKeywordsO              = iota
	ReadKeywordsOP             = iota
	ReadKeywordsOPE            = iota
	ReadKeywordsP              = iota
	ReadKeywordsPR             = iota
	ReadKeywordsPRO            = iota
	ReadKeywordsR              = iota
	ReadKeywordsRE             = iota
	ReadKeywordsREP            = iota
	ReadKeywordsREPE           = iota
	ReadKeywordsREPEA          = iota
	ReadKeywordsRET            = iota
	ReadKeywordsRETU           = iota
	ReadKeywordsRETUR          = iota
	ReadKeywordsS              = iota
	ReadKeywordsSE             = iota
	ReadKeywordsSEL            = iota
	ReadKeywordsSELE           = iota
	ReadKeywordsSELEC          = iota
	ReadKeywordsSI             = iota
	ReadKeywordsSIZ            = iota
	ReadKeywordsSIZE           = iota
	ReadKeywordsSIZEO          = iota
	ReadKeywordsST             = iota
	ReadKeywordsSTR            = iota
	ReadKeywordsSTRU           = iota
	ReadKeywordsSTRUC          = iota
	ReadKeywordsSW             = iota
	ReadKeywordsSWI            = iota
	ReadKeywordsSWIT           = iota
	ReadKeywordsSWITC          = iota
	ReadKeywordsT              = iota
	ReadKeywordsTE             = iota
	ReadKeywordsTES            = iota
	ReadKeywordsTEST           = iota
	ReadKeywordsTH             = iota
	ReadKeywordsTHI            = iota
	ReadKeywordsTHR            = iota
	ReadKeywordsTHRO           = iota
	ReadKeywordsTR             = iota
	ReadKeywordsTRU            = iota
	ReadKeywordsTY             = iota
	ReadKeywordsTYP            = iota
	ReadKeywordsTYPE           = iota
	ReadKeywordsTYPEO          = iota
	ReadKeywordsU              = iota
	ReadKeywordsUN             = iota
	ReadKeywordsUNS            = iota
	ReadKeywordsUNSA           = iota
	ReadKeywordsUNSAF          = iota
	ReadKeywordsUNT            = iota
	ReadKeywordsUNTI           = iota
	ReadKeywordsV              = iota
	ReadKeywordsVA             = iota
	ReadKeywordsW              = iota
	ReadKeywordsWH             = iota
	ReadKeywordsWHE            = iota
	ReadKeywordsWHI            = iota
	ReadKeywordsWHIL           = iota
	ReadKeywordsWI             = iota
	ReadKeywordsWIT            = iota
	ReadKeywordsX              = iota
	ReadKeywordsXO             = iota
)

func (this *Lexer) ReadChar(thisChar byte, nextChar byte, nextNextChar byte, tempToken *Token) (Token, bool) {
	return_token := Token {}
	result       := false

	// Acceptance States
	switch this.StateID {
		case AcceptWhiteSpace:
			return_token = Token {
				ID:       "whitespace",
				RawValue: tempToken.RawValue,
				Value:    string(tempToken.RawValue),

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: false,
			}
			this.StateID = Start
			result = true

		case AcceptComment:
			return_token = Token {
				ID:       "comment",
				RawValue: tempToken.RawValue,
				Value:    string(tempToken.RawValue),

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: false,
			}
			this.StateID = Start
			result = true

		case AcceptLineEnd:
			return_token = Token {
				ID:       "line end",
				RawValue: tempToken.RawValue,
				Value:    "\\n",

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: false,
			}
			this.StateID = Start
			result = true

		case AcceptIdent:
			return_token = Token {
				ID:       "identifier",
				RawValue: tempToken.RawValue,
				Value:    string(tempToken.RawValue),

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: false,
			}
			this.StateID = Start
			result = true

		case AcceptKeyword:
			return_token = Token {
				ID:       "keyword: " + string(tempToken.RawValue),
				RawValue: tempToken.RawValue,
				Value:    string(tempToken.RawValue),

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: false,
			}
			this.StateID = Start
			result = true

		case AcceptInteger:
			value    := string(tempToken.RawValue)
			valueInt := int64(0)
			if strings.HasPrefix(value, "0x") {
				valueInt, _ = strconv.ParseInt(value[2:], 16, 64)
			} else if strings.HasPrefix(value, "0b") {
				valueInt, _ = strconv.ParseInt(value[2:], 2, 64)
			} else {
				valueInt, _ = strconv.ParseInt(value, 10, 64)
			}

			return_token = Token {
				ID:       "integer literal",
				RawValue: tempToken.RawValue,
				Value:    valueInt,

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: false,
			}
			this.StateID = Start
			result = true

		case AcceptByte:
			value    := string(tempToken.RawValue)
			valueInt := int64(0)
			if strings.HasPrefix(value, "$%") {
				valueInt, _ = strconv.ParseInt(value[2:], 2, 9)
			} else if strings.HasPrefix(value, "$") {
				valueInt, _ = strconv.ParseInt(value[1:], 16, 9)
			} else {
				valueInt, _ = strconv.ParseInt(value[:len(value)-1], 10, 9)
			}

			return_token = Token {
				ID:       "byte literal",
				RawValue: tempToken.RawValue,
				Value:    uint8(valueInt & 0xFF),

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: false,
			}
			this.StateID = Start
			result = true

		case AcceptWord:
			value    := string(tempToken.RawValue)
			valueInt := int64(0)
			if strings.HasPrefix(value, "$%") {
				valueInt, _ = strconv.ParseInt(value[2:], 2, 17)
			} else if strings.HasPrefix(value, "$") {
				valueInt, _ = strconv.ParseInt(value[1:], 16, 17)
			} else {
				valueInt, _ = strconv.ParseInt(value[:len(value)-1], 10, 17)
			}

			return_token = Token {
				ID:       "word literal",
				RawValue: tempToken.RawValue,
				Value:    uint16(valueInt & 0xFFFF),

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: false,
			}
			this.StateID = Start
			result = true

		case AcceptLong:
			value    := string(tempToken.RawValue)
			valueInt := int64(0)
			if strings.HasPrefix(value, "$%") {
				valueInt, _ = strconv.ParseInt(value[2:], 2, 25)
			} else if strings.HasPrefix(value, "$") {
				valueInt, _ = strconv.ParseInt(value[1:], 16, 25)
			} else {
				valueInt, _ = strconv.ParseInt(value[:len(value)-1], 10, 25)
			}

			return_token = Token {
				ID:       "long literal",
				RawValue: tempToken.RawValue,
				Value:    uint32(valueInt & 0xFFFFFF),

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: false,
			}
			this.StateID = Start
			result = true

		case AcceptDWord:
			value    := string(tempToken.RawValue)
			valueInt := int64(0)
			if strings.HasPrefix(value, "$%") {
				valueInt, _ = strconv.ParseInt(value[2:], 2, 33)
			} else if strings.HasPrefix(value, "$") {
				valueInt, _ = strconv.ParseInt(value[1:], 16, 33)
			} else {
				valueInt, _ = strconv.ParseInt(value[:len(value)-1], 10, 33)
			}
			
			return_token = Token {
				ID:       "dword literal",
				RawValue: tempToken.RawValue,
				Value:    uint32(valueInt & 0xFFFFFFFF),

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: false,
			}
			this.StateID = Start
			result = true

		case AcceptQWord:
			value    := string(tempToken.RawValue)
			valueInt := uint64(0)
			if strings.HasPrefix(value, "$%") {
				valueInt, _ = strconv.ParseUint(value[2:], 2, 64)
			} else if strings.HasPrefix(value, "$") {
				valueInt, _ = strconv.ParseUint(value[1:], 16, 64)
			} else {
				valueInt, _ = strconv.ParseUint(value[:len(value)-1], 10, 64)
			}

			return_token = Token {
				ID:       "qword literal",
				RawValue: tempToken.RawValue,
				Value:    valueInt,

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: false,
			}
			this.StateID = Start
			result = true

		case AcceptDecimal:
			value    := string(tempToken.RawValue)
			valueNum := float64(0.0)
			valueNum, _ = strconv.ParseFloat(value, 64)

			return_token = Token {
				ID:       "decimal literal",
				RawValue: tempToken.RawValue,
				Value:    valueNum,

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: false,
			}
			this.StateID = Start
			result = true

		case AcceptFloat:
			value    := string(tempToken.RawValue)
			valueNum := float64(0.0)
			valueNum, _ = strconv.ParseFloat(value[:len(value)-1], 64)

			return_token = Token {
				ID:       "floating-point literal",
				RawValue: tempToken.RawValue,
				Value:    valueNum,

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: false,
			}
			this.StateID = Start
			result = true

		case AcceptDouble:
			value    := string(tempToken.RawValue)
			valueNum := float64(0.0)
			valueNum, _ = strconv.ParseFloat(value[:len(value)-1], 64)

			return_token = Token {
				ID:       "double-precision-floating-point literal",
				RawValue: tempToken.RawValue,
				Value:    valueNum,

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: false,
			}
			this.StateID = Start
			result = true

		case AcceptFixed:
			value    := string(tempToken.RawValue)
			valueNum := float64(0.0)
			valueNum, _ = strconv.ParseFloat(value[:len(value)-1], 64)

			return_token = Token {
				ID:       "fixed-point literal",
				RawValue: tempToken.RawValue,
				Value:    valueNum,

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: false,
			}
			this.StateID = Start
			result = true

		case AcceptOperator:
			return_token = Token {
				ID:       string(tempToken.RawValue),
				RawValue: tempToken.RawValue,
				Value:    string(tempToken.RawValue),

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: false,
			}
			this.StateID = Start
			result = true

		case AcceptCStringStart:
			return_token = Token {
				ID:       "string start",
				RawValue: tempToken.RawValue,
				Value:    string(tempToken.RawValue),

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: false,
			}
			this.StateID = ReadString
			result = true

		case AcceptWStringPrefix:
			return_token = Token {
				ID:       "w-string prefix",
				RawValue: tempToken.RawValue,
				Value:    string(tempToken.RawValue),

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: false,
			}
			this.StateID = Start
			result = true

		case AcceptUStringPrefix:
			return_token = Token {
				ID:       "u-string prefix",
				RawValue: tempToken.RawValue,
				Value:    string(tempToken.RawValue),

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: false,
			}
			this.StateID = Start
			result = true

		case AcceptStringChar:
			value := string(tempToken.RawValue)
			if value[0] == '\\' {
				if value[1] == 'x' {
					x, _ := strconv.ParseUint(value[2:], 16, 8)
					value = string(byte(x))
				} else if value[1] == 'u' {
					u, _ := strconv.ParseUint(value[2:], 16, 32)
					value = string(rune(u))
				} else {
					switch value[1] {
						case '"':  value = "\""
						case '\\': value = "\\"
						case 'b':  value = "\b"
						case 'f':  value = "\f"
						case 'n':  value = "\n"
						case 'r':  value = "\r"
						case 't':  value = "\t"
					}
				}
			}

			return_token = Token {
				ID:       "string character",
				RawValue: tempToken.RawValue,
				Value:    value,

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: false,
			}
			this.StateID = ReadString
			result = true

		case AcceptStringEnd:
			return_token = Token {
				ID:       "string end",
				RawValue: tempToken.RawValue,
				Value:    string(tempToken.RawValue),

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: false,
			}
			this.StateID = Start
			result = true

		case AcceptUCharPrefix:
			return_token = Token {
				ID:       "u-character prefix",
				RawValue: tempToken.RawValue,
				Value:    string(tempToken.RawValue),

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: false,
			}
			this.StateID = Start
			result = true

		case AcceptWCharPrefix:
			return_token = Token {
				ID:       "w-character prefix",
				RawValue: tempToken.RawValue,
				Value:    string(tempToken.RawValue),

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: false,
			}
			this.StateID = Start
			result = true

		case AcceptChar:
			value := string(tempToken.RawValue)
			if value[1] == '\\' {
				if value[2] == 'x' {
					x, _ := strconv.ParseUint(value[3:len(value)-1], 16, 8)
					value = string(byte(x))
				} else if value[2] == 'u' {
					u, _ := strconv.ParseUint(value[3:len(value)-1], 16, 32)
					value = string(rune(u))
				} else {
					switch value[2] {
						case '\'':  value = "'"
						case '\\': value = "\\"
						case 'b':  value = "\b"
						case 'f':  value = "\f"
						case 'n':  value = "\n"
						case 'r':  value = "\r"
						case 't':  value = "\t"
					}
				}
			} else {
				value = value[1:len(value)-1]
			}

			return_token = Token {
				ID:       "character",
				RawValue: tempToken.RawValue,
				Value:    value,

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: false,
			}
			this.StateID = Start
			result = true

		case RejectStringChar:
			return_token = Token {
				ID:       "invalid string character",
				RawValue: tempToken.RawValue,
				Value:    string(tempToken.RawValue),

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: true,
			}
			this.StateID = ReadString
			result = true

		case RejectStringEnd:
			return_token = Token {
				ID:       "line end in string literal",
				RawValue: tempToken.RawValue,
				Value:    string(tempToken.RawValue),

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: true,
			}
			this.StateID = Start
			result = true

		case RejectChar:
			return_token = Token {
				ID:       "invalid character in string",
				RawValue: tempToken.RawValue,
				Value:    string(tempToken.RawValue),

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: true,
			}
			this.StateID = ReadString
			result = true

		case RejectNumber:
			return_token = Token {
				ID:       "invalid number literal",
				RawValue: tempToken.RawValue,
				Value:    string(tempToken.RawValue),

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: true,
			}
			this.StateID = Start
			result = true

		case RejectLiteral:
			return_token = Token {
				ID:       "invalid literal",
				RawValue: tempToken.RawValue,
				Value:    string(tempToken.RawValue),

				FilePath: tempToken.FilePath,
				Line:     tempToken.Line,
				Column:   tempToken.Column,

				Error: true,
			}
			this.StateID = Start
			result = true

		case RejectCharacter:
			if utf8.RuneCountInString(string(tempToken.RawValue)) == 3 && len(string(tempToken.RawValue)) != 3 && string(tempToken.RawValue)[0] == '\'' && string(tempToken.RawValue)[len(string(tempToken.RawValue))-1] == '\'' {
				return_token = Token {
					ID:       "character",
					RawValue: tempToken.RawValue,
					Value:    string(tempToken.RawValue)[1:len(string(tempToken.RawValue))-1],

					FilePath: tempToken.FilePath,
					Line:     tempToken.Line,
					Column:   tempToken.Column,

					Error: false,
				}
			} else {
				return_token = Token {
					ID:       "invalid character literal",
					RawValue: tempToken.RawValue,
					Value:    string(tempToken.RawValue),

					FilePath: tempToken.FilePath,
					Line:     tempToken.Line,
					Column:   tempToken.Column,

					Error: true,
				}
			}
			this.StateID = Start
			result = true
	}

	// Transitions
	switch this.StateID {
		case Start: // 0
			*tempToken = Token {
				ID:       "unknown",
				RawValue: [] byte {},
				Value:    "",

				FilePath: this.FilePath,
				Line:     this.Line,
				Column:   this.Column,

				Error: false,
			}

			if thisChar == 'w' && nextChar == '\'' {
				this.StateID = AcceptWCharPrefix
			} else if thisChar == 'w' && nextChar == '"' {
				this.StateID = AcceptWStringPrefix
			} else if thisChar == 'u' && nextChar == '\'' {
				this.StateID = AcceptUCharPrefix
			} else if thisChar == 'u' && nextChar == '"' {
				this.StateID = AcceptUStringPrefix
			} else if IsKeywordChar(thisChar) && !IsIdentChar(nextChar) {
				this.StateID = AcceptIdent
			} else if IsKeywordChar(thisChar) && IsKeywordChar(nextChar) {
				switch thisChar {
					case 'a': this.StateID = ReadKeywordsA
					case 'b': this.StateID = ReadKeywordsB
					case 'c': this.StateID = ReadKeywordsC
					case 'd': this.StateID = ReadKeywordsD
					case 'e': this.StateID = ReadKeywordsE
					case 'f': this.StateID = ReadKeywordsF
					//case 'g': this.StateID = ReadKeywordsG
					case 'i': this.StateID = ReadKeywordsI
					case 'l': this.StateID = ReadKeywordsL
					case 'm': this.StateID = ReadKeywordsM
					case 'n': this.StateID = ReadKeywordsN
					case 'o': this.StateID = ReadKeywordsO
					case 'p': this.StateID = ReadKeywordsP
					case 'r': this.StateID = ReadKeywordsR
					case 's': this.StateID = ReadKeywordsS
					case 't': this.StateID = ReadKeywordsT
					case 'u': this.StateID = ReadKeywordsU
					case 'v': this.StateID = ReadKeywordsV
					case 'w': this.StateID = ReadKeywordsW
					case 'x': this.StateID = ReadKeywordsX
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(thisChar) && !IsNumChar(thisChar) && !IsIdentChar(nextChar) {
				this.StateID = AcceptIdent
			} else if IsIdentChar(thisChar) && !IsNumChar(thisChar) && IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else if thisChar == '0' && IsIdentChar(nextChar) {
				this.StateID = ReadNumberZero
			} else if thisChar == '0' && nextChar == '.' {
				if nextNextChar == '.' {
					this.StateID = AcceptInteger
				} else {
					this.StateID = ReadNumberZero
				}
			} else if thisChar == '0' {
				this.StateID = AcceptInteger
			} else if IsNumChar(thisChar) && IsIdentChar(nextChar) {
				this.StateID = ReadNumberDigit
			} else if IsNumChar(thisChar) && nextChar == '.' {
				if nextNextChar == '.' {
					this.StateID = AcceptInteger
				} else {
					this.StateID = ReadNumberDigit
				}
			} else if thisChar == '$' && nextChar == 'w' {
				//this.StateID = ReadHexDigit1
			} else if thisChar == '$' && nextChar == 'u' {
				//this.StateID = ReadHexDigit1
			} else if thisChar == '$' && nextChar == '"' {
				//this.StateID = ReadHexDigit1
			} else if thisChar == '$' && IsIdentChar(nextChar) {
				this.StateID = ReadHexDigit1
			} else if thisChar == '$' && nextChar == '.' {
				if nextNextChar == '.' {
					this.StateID = RejectLiteral
				} else {
					this.StateID = ReadHexDigit1
				}
			} else if thisChar == '$' && nextChar == '%' {
				this.StateID = ReadBinaryDigit0
			} else if thisChar == '$' {
				this.StateID = RejectLiteral
			} else if IsNumChar(thisChar) {
				this.StateID = AcceptInteger
			} else if thisChar == '~' || thisChar == '@' || thisChar == '#' || thisChar == '?' || thisChar == ',' || thisChar == '(' || thisChar == ')' || thisChar == '[' || thisChar == ']' || thisChar == '{' || thisChar == '}' {
				this.StateID = AcceptOperator
			} else if thisChar == '.' && nextChar == '.' {
				this.StateID = ReadDot1
			} else if thisChar == '.' {
				this.StateID = AcceptOperator
			} else if thisChar == '=' && nextChar == '=' {
				this.StateID = ReadEquals1
			} else if thisChar == '=' {
				this.StateID = AcceptOperator
			} else if thisChar == ':' && (nextChar == '=' || nextChar == ':') {
				this.StateID = ReadColon
			} else if thisChar == ':' {
				this.StateID = AcceptOperator
			} else if thisChar == '+' && nextChar == '=' {
				this.StateID = ReadPlus
			} else if thisChar == '+' {
				this.StateID = AcceptOperator
			} else if thisChar == '-' && (nextChar == '=' || nextChar == '>' || nextChar == '-') {
				this.StateID = ReadMinus
			} else if thisChar == '-' {
				this.StateID = AcceptOperator
			} else if thisChar == '*' && nextChar == '=' {
				this.StateID = ReadStar
			} else if thisChar == '*' {
				this.StateID = AcceptOperator
			} else if thisChar == '/' && (nextChar == '=' || nextChar == '/' || nextChar == '*') {
				this.StateID = ReadSlash
			} else if thisChar == '/' {
				this.StateID = AcceptOperator
			} else if thisChar == '%' && nextChar == '=' {
				this.StateID = ReadPercent
			} else if thisChar == '%' {
				this.StateID = AcceptOperator
			} else if thisChar == '&' && nextChar == '=' {
				this.StateID = ReadAmp
			} else if thisChar == '&' {
				this.StateID = AcceptOperator
			} else if thisChar == '|' && nextChar == '=' {
				this.StateID = ReadPipe
			} else if thisChar == '|' {
				this.StateID = AcceptOperator
			} else if thisChar == '^' && nextChar == '=' {
				this.StateID = ReadCaret
			} else if thisChar == '^' {
				this.StateID = AcceptOperator
			} else if thisChar == '>' && (nextChar == '=' || nextChar == '>') {
				this.StateID = ReadGreaterThan
			} else if thisChar == '>' {
				this.StateID = AcceptOperator
			} else if thisChar == '<' && (nextChar == '=' || nextChar == '<' || (nextChar == '-' && nextNextChar != '-')) {
				this.StateID = ReadLessThan
			} else if thisChar == '<' {
				this.StateID = AcceptOperator
			} else if thisChar == '!' && nextChar == '=' {
				this.StateID = ReadBang
			} else if thisChar == '!' {
				this.StateID = AcceptOperator
			} else if thisChar == '\\' {
				this.StateID = AcceptOperator
			} else if thisChar == '"' {
				this.StateID = AcceptCStringStart
			} else if thisChar == '\'' {
				this.StateID = ReadChar
			} else if thisChar == ';' && (nextChar == '\r' || nextChar == '\n') {
				this.StateID = AcceptComment
			} else if thisChar == ';' {
				this.StateID = ReadComment
			} else if IsWhitespace(thisChar) && !IsWhitespace(nextChar) {
				this.StateID = AcceptWhiteSpace
			} else if IsWhitespace(thisChar) && IsWhitespace(nextChar) {
				this.StateID = ReadWhitespace
			} else if thisChar == '\n' {
				this.StateID = AcceptLineEnd
			} else {
				this.StateID = RejectCharacter
			}

		case ReadWhitespace:
			if !IsWhitespace(nextChar) {
				this.StateID = AcceptWhiteSpace
			}

		case ReadComment:
			if nextChar == '\r' || nextChar == '\n' {
				this.StateID = AcceptComment
			}

		case ReadMultiLineComment:
			if thisChar == '*' && nextChar == '/' {
				this.StateID = ReadMultiLineComment2
			}

		case ReadMultiLineComment2:
			if thisChar == '/' {
				this.StateID = AcceptComment
			}

		case ReadIdent:
			if !IsIdentChar(nextChar) {
				this.StateID = AcceptIdent
			}

		case ReadString:
			*tempToken = Token {
				ID:       "unknown",
				RawValue: [] byte {},
				Value:    "",

				FilePath: this.FilePath,
				Line:     this.Line,
				Column:   this.Column,

				Error: false,
			}

			if thisChar == '\n' {
				this.StateID = RejectStringEnd
			} else if thisChar == '\r' {
				this.StateID = RejectStringChar
			} else if thisChar == '"' {
				this.StateID = AcceptStringEnd
			} else if thisChar == '\\' && nextChar != '\r' && nextChar != '\n' {
				this.StateID = ReadStringBackslash
			} else if thisChar == '\\' {
				this.StateID = RejectStringChar
			} else {
				this.StateID = AcceptStringChar
			}

		case ReadStringBackslash:
			if thisChar == 'u' && (nextChar == '\r' || nextChar == '\n' || nextChar == '"') {
				this.StateID = RejectStringChar
			} else if thisChar == 'u' {
				this.StateID = ReadStringBackslashU
			} else if thisChar == 'x' && (nextChar == '\r' || nextChar == '\n' || nextChar == '"') {
				this.StateID = RejectStringChar
			} else if thisChar == 'x' {
				this.StateID = ReadStringBackslashX
			} else if thisChar == '\n' {
				this.StateID = RejectStringEnd
			} else if thisChar != 'b' && thisChar != 'f' && thisChar != 'n' && thisChar != 'r' && thisChar != 't' && thisChar != '\\' && thisChar != '"' {
				this.StateID = RejectStringChar
			} else {
				this.StateID = AcceptStringChar
			}

		case ReadStringBackslashU:
			if nextChar == '\r' || nextChar == '\n' || nextChar == '"' {
				this.StateID = RejectStringChar
			} else if IsHexChar(thisChar) {
				this.StateID = ReadStringBackslashU1
			} else {
				this.StateID = InvalidReadBackslashU1
			}

		case ReadStringBackslashU1:
			if nextChar == '\r' || nextChar == '\n' || nextChar == '"' {
				this.StateID = RejectStringChar
			} else if IsHexChar(thisChar) {
				this.StateID = ReadStringBackslashU2
			} else {
				this.StateID = InvalidReadBackslashU2
			}

		case ReadStringBackslashU2:
			if nextChar == '\r' || nextChar == '\n' || nextChar == '"' {
				this.StateID = RejectStringChar
			} else if IsHexChar(thisChar) {
				this.StateID = ReadStringBackslashU3
			} else {
				this.StateID = InvalidReadBackslashU3
			}

		case ReadStringBackslashU3:
			if IsHexChar(thisChar) {
				this.StateID = AcceptStringChar
			} else {
				this.StateID = RejectStringChar
			}

		case InvalidReadBackslashU1:
			if nextChar == '\r' || nextChar == '\n' || nextChar == '"' {
				this.StateID = RejectStringChar
			} else {
				this.StateID = InvalidReadBackslashU2
			}

		case InvalidReadBackslashU2:
			if nextChar == '\r' || nextChar == '\n' || nextChar == '"' {
				this.StateID = RejectStringChar
			} else {
				this.StateID = InvalidReadBackslashU3
			}

		case InvalidReadBackslashU3:
			this.StateID = RejectStringChar

		case ReadStringBackslashX:
			if nextChar == '\r' || nextChar == '\n' || nextChar == '"' {
				this.StateID = RejectStringChar
			} else if IsHexChar(thisChar) {
				this.StateID = ReadStringBackslashX1
			} else {
				this.StateID = InvalidReadBackslashX1
			}

		case ReadStringBackslashX1:
			if IsHexChar(thisChar) {
				this.StateID = AcceptStringChar
			} else {
				this.StateID = RejectStringChar
			}

		case InvalidReadBackslashX1:
			this.StateID = RejectStringChar

		case ReadChar:
			if thisChar == '\n' {
				this.StateID = RejectCharacter
			} else if thisChar == '\r' {
				this.StateID = InvalidReadChar
			} else if thisChar == '\'' {
				this.StateID = RejectCharacter
			} else if thisChar == '\\' && nextChar != '\r' && nextChar != '\n' {
				this.StateID = ReadCharBackslash
			} else if thisChar == '\\' {
				this.StateID = InvalidReadChar
			} else {
				this.StateID = ReadCharFinal
			}

		case ReadCharBackslash:
			if thisChar == 'u' && (nextChar == '\r' || nextChar == '\n' || nextChar == '\'') {
				this.StateID = InvalidReadChar
			} else if thisChar == 'u' {
				this.StateID = ReadCharBackslashU
			} else if thisChar == 'x' && (nextChar == '\r' || nextChar == '\n' || nextChar == '\'') {
				this.StateID = InvalidReadChar
			} else if thisChar == 'x' {
				this.StateID = ReadCharBackslashX
			} else if thisChar == '\n' {
				this.StateID = RejectCharacter
			} else if thisChar != 'b' && thisChar != 'f' && thisChar != 'n' && thisChar != 'r' && thisChar != 't' && thisChar != '\\' && thisChar != '\'' {
				this.StateID = InvalidReadChar
			} else {
				this.StateID = ReadCharFinal
			}

		case ReadCharBackslashU:
			if thisChar == '\r' || thisChar == '\n' || thisChar == '\'' {
				this.StateID = RejectCharacter
			} else if IsHexChar(thisChar) {
				this.StateID = ReadCharBackslashU1
			} else {
				this.StateID = InvalidReadChar
			}

		case ReadCharBackslashU1:
			if thisChar == '\r' || thisChar == '\n' || thisChar == '\'' {
				this.StateID = RejectCharacter
			} else if IsHexChar(thisChar) {
				this.StateID = ReadCharBackslashU2
			} else {
				this.StateID = InvalidReadChar
			}

		case ReadCharBackslashU2:
			if thisChar == '\r' || thisChar == '\n' || thisChar == '\'' {
				this.StateID = RejectCharacter
			} else if IsHexChar(thisChar) {
				this.StateID = ReadCharBackslashU3
			} else {
				this.StateID = InvalidReadChar
			}

		case ReadCharBackslashU3:
			if IsHexChar(thisChar) {
				this.StateID = ReadCharFinal
			} else {
				this.StateID = InvalidReadChar
			}

		case ReadCharBackslashX:
			if thisChar == '\r' || thisChar == '\n' || thisChar == '\'' {
				this.StateID = RejectCharacter
			} else if IsHexChar(thisChar) {
				this.StateID = ReadCharBackslashX1
			} else {
				this.StateID = InvalidReadChar
			}

		case ReadCharBackslashX1:
			if IsHexChar(thisChar) {
				this.StateID = ReadCharFinal
			} else {
				this.StateID = InvalidReadChar
			}

		case InvalidReadChar:
			if thisChar == '\'' || thisChar == '\n' {
				this.StateID = RejectCharacter
			}

		case ReadCharFinal:
			if thisChar == '\'' {
				this.StateID = AcceptChar
			} else if thisChar == '\n' {
				this.StateID = RejectCharacter
			} else {
				this.StateID = InvalidReadChar
			}

		case ReadSlash:
			if thisChar == '=' {
				this.StateID = AcceptOperator
			} else if thisChar == '/' && (nextChar == '\r' || nextChar == '\n') {
				this.StateID = AcceptComment
			} else if thisChar == '/' {
				this.StateID = ReadComment
			} else if thisChar == '*' {
				this.StateID = ReadMultiLineComment
			}

		case ReadMinus:
			if thisChar == '=' || thisChar == '>' {
				this.StateID = AcceptOperator
			} else if thisChar == '-' && (nextChar == '\r' || nextChar == '\n') {
				this.StateID = AcceptComment
			} else {
				this.StateID = ReadComment
			}

		case ReadDot1:
			if thisChar == '.' && nextChar == '.' {
				this.StateID = ReadDot2
			} else if thisChar == '.' {
				this.StateID = AcceptOperator
			}

		case ReadDot2:
			this.StateID = AcceptOperator

		case ReadEquals1:
			if thisChar == '=' && nextChar == '=' {
				this.StateID = ReadEquals2
			} else if thisChar == '=' {
				this.StateID = AcceptOperator
			}

		case ReadEquals2:
			this.StateID = AcceptOperator

		case ReadColon:
			this.StateID = AcceptOperator

		case ReadPlus:
			this.StateID = AcceptOperator

		case ReadStar:
			this.StateID = AcceptOperator

		case ReadPercent:
			this.StateID = AcceptOperator

		case ReadAmp:
			this.StateID = AcceptOperator

		case ReadPipe:
			this.StateID = AcceptOperator

		case ReadCaret:
			this.StateID = AcceptOperator

		case ReadGreaterThan:
			if (thisChar == '>' && nextChar != '=') || thisChar == '=' {
				this.StateID = AcceptOperator
			} else if thisChar == '>' && nextChar == '=' {
				this.StateID = ReadShiftRight
			}

		case ReadShiftRight:
			this.StateID = AcceptOperator

		case ReadLessThan:
			if (thisChar == '<' && nextChar != '=') || thisChar == '=' || thisChar == '-' {
				this.StateID = AcceptOperator
			} else if thisChar == '<' && nextChar == '=' {
				this.StateID = ReadShiftLeft
			}

		case ReadShiftLeft:
			this.StateID = AcceptOperator

		case ReadBang:
			if thisChar == '=' && nextChar == '=' {
				this.StateID = ReadBangEquals
			} else if thisChar == '=' {
				this.StateID = AcceptOperator
			}

		case ReadBangEquals:
			this.StateID = AcceptOperator

		case ReadNumberZero:
			if thisChar == 'B' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptByte
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptByte
				} else {
					this.StateID = InvalidReadNumber
				}
			} else if thisChar == 'w' || thisChar == 'W' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptWord
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptWord
				} else {
					this.StateID = InvalidReadNumber
				}
			} else if thisChar == 'l' || thisChar == 'L' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptLong
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptLong
				} else {
					this.StateID = InvalidReadNumber
				}
			} else if thisChar == 'd' || thisChar == 'D' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptDWord
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptDWord
				} else {
					this.StateID = InvalidReadNumber
				}
			} else if thisChar == 'q' || thisChar == 'Q' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptQWord
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptQWord
				} else {
					this.StateID = InvalidReadNumber
				}
			} else if thisChar == 'f' || thisChar == 'F' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptFloat
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptFloat
				} else {
					this.StateID = InvalidReadNumber
				}
			} else if thisChar == 'b' && !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = AcceptByte
			} else if thisChar == 'b' && nextChar == '.' && nextNextChar == '.' {
				this.StateID = AcceptByte
			} else if thisChar == 'b' && IsNumChar(nextChar) {
				this.StateID = ReadBinaryDigit
			} else if thisChar == 'b' {
				this.StateID = InvalidReadNumber
			} else if thisChar == 'x' && IsHexChar(nextChar) {
				this.StateID = ReadHexDigit
			} else if (thisChar == 'x' || thisChar == 'X') && nextChar == '.' && nextNextChar == '.' {
				this.StateID = AcceptFixed
			} else if thisChar == 'x' && (IsIdentChar(nextChar) || nextChar == '.') {
				this.StateID = InvalidReadNumber
			} else if (thisChar == 'x' || thisChar == 'X') && !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = AcceptFixed
			} else if thisChar == '.' && IsNumChar(nextChar) {
				this.StateID = ReadDecimalDigit
			} else if thisChar == '.' && !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsNumChar(thisChar) && nextChar == '.' && nextNextChar == '.' {
				this.StateID = AcceptInteger
			} else if IsNumChar(thisChar) && !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = AcceptInteger
			} else if IsNumChar(thisChar) {
				this.StateID = ReadNumberDigit
			} else {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = RejectNumber
				} else {
					this.StateID = InvalidReadNumber
				}
			}

		case ReadNumberDigit:
			if IsNumChar(thisChar) {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptInteger
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptInteger
				}
			} else if (thisChar == 'b' || thisChar == 'B') && !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = AcceptByte
			} else if (thisChar == 'b' || thisChar == 'B') && nextChar == '.' && nextNextChar == '.' {
				this.StateID = AcceptByte
			} else if (thisChar == 'w' || thisChar == 'W') && !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = AcceptWord
			} else if (thisChar == 'w' || thisChar == 'W') && nextChar == '.' && nextNextChar == '.' {
				this.StateID = AcceptWord
			} else if (thisChar == 'l' || thisChar == 'L') && !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = AcceptLong
			} else if (thisChar == 'l' || thisChar == 'L') && nextChar == '.' && nextNextChar == '.' {
				this.StateID = AcceptLong
			} else if (thisChar == 'd' || thisChar == 'D') && !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = AcceptDWord
			} else if (thisChar == 'd' || thisChar == 'D') && nextChar == '.' && nextNextChar == '.' {
				this.StateID = AcceptDWord
			} else if (thisChar == 'q' || thisChar == 'Q') && !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = AcceptQWord
			} else if (thisChar == 'q' || thisChar == 'Q') && nextChar == '.' && nextNextChar == '.' {
				this.StateID = AcceptQWord
			} else if (thisChar == 'f' || thisChar == 'F') && !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = AcceptFloat
			} else if (thisChar == 'f' || thisChar == 'F') && nextChar == '.' && nextNextChar == '.' {
				this.StateID = AcceptFloat
			} else if (thisChar == 'x' || thisChar == 'X') && !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = AcceptFixed
			} else if (thisChar == 'x' || thisChar == 'X') && nextChar == '.' && nextNextChar == '.' {
				this.StateID = AcceptFixed
			} else if thisChar == '.' && IsNumChar(nextChar) {
				this.StateID = ReadDecimalDigit
			} else if thisChar == '.' && !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			}

		case ReadDecimalDigit:
			if IsNumChar(thisChar) {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptDecimal
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptDecimal
				}
			} else if (thisChar == 'f' || thisChar == 'F') && !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = AcceptFloat
			} else if (thisChar == 'f' || thisChar == 'F') && nextChar == '.' && nextNextChar == '.' {
				this.StateID = AcceptFloat
			} else if (thisChar == 'd' || thisChar == 'D') && !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = AcceptDouble
			} else if (thisChar == 'd' || thisChar == 'D') && nextChar == '.' && nextNextChar == '.' {
				this.StateID = AcceptDouble
			} else if (thisChar == 'x' || thisChar == 'X') && !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = AcceptFixed
			} else if (thisChar == 'x' || thisChar == 'X') && nextChar == '.' && nextNextChar == '.' {
				this.StateID = AcceptFixed
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			}

		case ReadBinaryDigit:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptInteger
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptInteger
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			}

		case ReadHexDigit:
			if IsHexChar(thisChar) {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptInteger
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptInteger
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			}

		case ReadHexDigit1:
			if IsHexChar(thisChar) {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptByte
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptByte
				} else {
					this.StateID = ReadHexDigit2
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadHexDigit2
			}

		case ReadHexDigit2:
			if IsHexChar(thisChar) {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptByte
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptByte
				} else {
					this.StateID = ReadHexDigit3
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadHexDigit3
			}

		case ReadHexDigit3:
			if IsHexChar(thisChar) {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptWord
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptWord
				} else {
					this.StateID = ReadHexDigit4
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadHexDigit4
			}

		case ReadHexDigit4:
			if IsHexChar(thisChar) {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptWord
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptWord
				} else {
					this.StateID = ReadHexDigit5
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadHexDigit5
			}

		case ReadHexDigit5:
			if IsHexChar(thisChar) {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptLong
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptLong
				} else {
					this.StateID = ReadHexDigit6
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadHexDigit6
			}

		case ReadHexDigit6:
			if IsHexChar(thisChar) {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptLong
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptLong
				} else {
					this.StateID = ReadHexDigit7
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadHexDigit7
			}

		case ReadHexDigit7:
			if IsHexChar(thisChar) {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptDWord
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptDWord
				} else {
					this.StateID = ReadHexDigit8
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadHexDigit8
			}

		case ReadHexDigit8:
			if IsHexChar(thisChar) {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptDWord
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptDWord
				} else {
					this.StateID = ReadHexDigitQ
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadHexDigitQ
			}

		case ReadHexDigitQ:
			if IsHexChar(thisChar) {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptQWord
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptQWord
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			}

		case ReadBinaryDigit0:
			if thisChar == '%' && IsIdentChar(nextChar) {
				this.StateID = ReadBinaryDigit1
			} else if thisChar == '%' && nextChar == '.' {
				if nextNextChar == '.' {
					this.StateID = RejectNumber
				} else {
					this.StateID = ReadBinaryDigit1
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else {
				this.StateID = InvalidReadNumber
			}

		case ReadBinaryDigit1:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptByte
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptByte
				} else {
					this.StateID = ReadBinaryDigit2
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit2
			}

		case ReadBinaryDigit2:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptByte
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptByte
				} else {
					this.StateID = ReadBinaryDigit3
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit3
			}

		case ReadBinaryDigit3:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptByte
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptByte
				} else {
					this.StateID = ReadBinaryDigit4
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit4
			}

		case ReadBinaryDigit4:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptByte
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptByte
				} else {
					this.StateID = ReadBinaryDigit5
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit5
			}

		case ReadBinaryDigit5:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptByte
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptByte
				} else {
					this.StateID = ReadBinaryDigit6
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit6
			}

		case ReadBinaryDigit6:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptByte
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptByte
				} else {
					this.StateID = ReadBinaryDigit7
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit7
			}

		case ReadBinaryDigit7:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptByte
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptByte
				} else {
					this.StateID = ReadBinaryDigit8
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit8
			}

		case ReadBinaryDigit8:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptByte
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptByte
				} else {
					this.StateID = ReadBinaryDigit9
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit9
			}

		case ReadBinaryDigit9:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptWord
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptWord
				} else {
					this.StateID = ReadBinaryDigit10
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit10
			}

		case ReadBinaryDigit10:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptWord
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptWord
				} else {
					this.StateID = ReadBinaryDigit11
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit11
			}

		case ReadBinaryDigit11:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptWord
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptWord
				} else {
					this.StateID = ReadBinaryDigit12
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit12
			}

		case ReadBinaryDigit12:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptWord
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptWord
				} else {
					this.StateID = ReadBinaryDigit13
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit13
			}

		case ReadBinaryDigit13:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptWord
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptWord
				} else {
					this.StateID = ReadBinaryDigit14
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit14
			}

		case ReadBinaryDigit14:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptWord
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptWord
				} else {
					this.StateID = ReadBinaryDigit15
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit15
			}

		case ReadBinaryDigit15:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptWord
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptWord
				} else {
					this.StateID = ReadBinaryDigit16
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit16
			}

		case ReadBinaryDigit16:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptWord
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptWord
				} else {
					this.StateID = ReadBinaryDigit17
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit17
			}

		case ReadBinaryDigit17:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptLong
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptLong
				} else {
					this.StateID = ReadBinaryDigit18
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit18
			}

		case ReadBinaryDigit18:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptLong
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptLong
				} else {
					this.StateID = ReadBinaryDigit19
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit19
			}

		case ReadBinaryDigit19:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptLong
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptLong
				} else {
					this.StateID = ReadBinaryDigit20
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit20
			}

		case ReadBinaryDigit20:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptLong
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptLong
				} else {
					this.StateID = ReadBinaryDigit21
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit21
			}

		case ReadBinaryDigit21:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptLong
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptLong
				} else {
					this.StateID = ReadBinaryDigit22
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit22
			}

		case ReadBinaryDigit22:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptLong
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptLong
				} else {
					this.StateID = ReadBinaryDigit23
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit23
			}

		case ReadBinaryDigit23:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptLong
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptLong
				} else {
					this.StateID = ReadBinaryDigit24
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit24
			}

		case ReadBinaryDigit24:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptLong
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptLong
				} else {
					this.StateID = ReadBinaryDigit25
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit25
			}

		case ReadBinaryDigit25:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptDWord
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptDWord
				} else {
					this.StateID = ReadBinaryDigit26
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit26
			}

		case ReadBinaryDigit26:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptDWord
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptDWord
				} else {
					this.StateID = ReadBinaryDigit27
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit27
			}

		case ReadBinaryDigit27:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptDWord
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptDWord
				} else {
					this.StateID = ReadBinaryDigit28
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit28
			}

		case ReadBinaryDigit28:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptDWord
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptDWord
				} else {
					this.StateID = ReadBinaryDigit29
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit29
			}

		case ReadBinaryDigit29:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptDWord
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptDWord
				} else {
					this.StateID = ReadBinaryDigit30
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit30
			}

		case ReadBinaryDigit30:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptDWord
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptDWord
				} else {
					this.StateID = ReadBinaryDigit31
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit31
			}

		case ReadBinaryDigit31:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptDWord
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptDWord
				} else {
					this.StateID = ReadBinaryDigit32
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigit32
			}

		case ReadBinaryDigit32:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptDWord
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptDWord
				} else {
					this.StateID = ReadBinaryDigitQ
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			} else {
				this.StateID = ReadBinaryDigitQ
			}

		case ReadBinaryDigitQ:
			if thisChar == '0' || thisChar == '1' {
				if !IsIdentChar(nextChar) && nextChar != '.' {
					this.StateID = AcceptQWord
				} else if nextChar == '.' && nextNextChar == '.' {
					this.StateID = AcceptQWord
				}
			} else if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if IsIdentChar(thisChar) || thisChar == '.' {
				this.StateID = InvalidReadNumber
			}

		case InvalidReadNumber:
			if !IsIdentChar(nextChar) && nextChar != '.' {
				this.StateID = RejectNumber
			} else if nextChar == '.' && nextNextChar == '.' {
				this.StateID = RejectNumber
			}

		case ReadKeywordsA:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'l': this.StateID = ReadKeywordsAL
					case 'n': this.StateID = ReadKeywordsAN
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 's' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsAL:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'i': this.StateID = ReadKeywordsALI
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsALI:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'a': this.StateID = ReadKeywordsALIA
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsALIA:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 's' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsAN:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'd' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsB:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'r': this.StateID = ReadKeywordsBR
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsBR:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'e': this.StateID = ReadKeywordsBRE
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsBRE:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'a': this.StateID = ReadKeywordsBREA
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsBREA:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'k' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsC:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'a': this.StateID = ReadKeywordsCA
					case 'o': this.StateID = ReadKeywordsCO
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsCA:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'c': this.StateID = ReadKeywordsCAC
					case 's': this.StateID = ReadKeywordsCAS
					case 't': this.StateID = ReadKeywordsCAT
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsCAC:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'h': this.StateID = ReadKeywordsCACH
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsCACH:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'e': this.StateID = ReadKeywordsCACHE
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsCACHE:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'd' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsCAS:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'e' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsCAT:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'c': this.StateID = ReadKeywordsCATC
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsCATC:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'h' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsCO:
			if IsIdentChar(nextChar) {
				if thisChar == 'n' && (nextChar == 's' || nextChar == 't') {
					this.StateID = ReadKeywordsCON
				} else {
					this.StateID = ReadIdent
				}
			} else {
				if thisChar == 'n' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsCON:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 's': this.StateID = ReadKeywordsCONS
					case 't': this.StateID = ReadKeywordsCONT
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsCONS:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 't': this.StateID = ReadKeywordsCONST
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsCONST:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'r': this.StateID = ReadKeywordsCONSTR
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsCONSTR:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'a': this.StateID = ReadKeywordsCONSTRA
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsCONSTRA:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'i': this.StateID = ReadKeywordsCONSTRAI
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsCONSTRAI:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'n' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsCONT:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'i': this.StateID = ReadKeywordsCONTI
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsCONTI:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'n': this.StateID = ReadKeywordsCONTIN
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsCONTIN:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'u': this.StateID = ReadKeywordsCONTINU
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsCONTINU:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'e' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsD:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'e': this.StateID = ReadKeywordsDE
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'o' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsDE:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'f': this.StateID = ReadKeywordsDEF
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsDEF:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'a': this.StateID = ReadKeywordsDEFA
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsDEFA:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'u': this.StateID = ReadKeywordsDEFAU
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsDEFAU:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'l': this.StateID = ReadKeywordsDEFAUL
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsDEFAUL:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 't' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsE:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'a': this.StateID = ReadKeywordsEA
					case 'l': this.StateID = ReadKeywordsEL
					case 'n': this.StateID = ReadKeywordsEN
					case 'v': this.StateID = ReadKeywordsEV
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsEA:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'c': this.StateID = ReadKeywordsEAC
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsEAC:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'h' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsEL:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 's': this.StateID = ReadKeywordsELS
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsELS:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'e' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsEN:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'u': this.StateID = ReadKeywordsENU
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsENU:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'm' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsEV:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'a': this.StateID = ReadKeywordsEVA
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsEVA:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'l' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsF:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'a': this.StateID = ReadKeywordsFA
					case 'i': this.StateID = ReadKeywordsFI
					case 'o': this.StateID = ReadKeywordsFO
					case 'u': this.StateID = ReadKeywordsFU
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsFA:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'l': this.StateID = ReadKeywordsFAL
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsFAL:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 's': this.StateID = ReadKeywordsFALS
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsFALS:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'e' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsFI:
			if IsIdentChar(nextChar) {
				if thisChar == 'n' && nextChar == 'a' {
					this.StateID = ReadKeywordsFIN
				} else {
					this.StateID = ReadIdent
				}
			} else {
				if thisChar == 'n' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsFIN:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'a': this.StateID = ReadKeywordsFINA
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsFINA:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'l': this.StateID = ReadKeywordsFINAL
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsFINAL:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'l': this.StateID = ReadKeywordsFINALL
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsFINALL:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'y' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsFO:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'r' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsFU:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'n': this.StateID = ReadKeywordsFUN
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsFUN:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'c' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsI:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'm': this.StateID = ReadKeywordsIM
					case 'n': this.StateID = ReadKeywordsIN
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'f' || thisChar == 's' || thisChar == 'n' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsIM:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'p': this.StateID = ReadKeywordsIMP
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsIMP:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'l': this.StateID = ReadKeywordsIMPL
					case 'o': this.StateID = ReadKeywordsIMPO
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsIMPL:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'e': this.StateID = ReadKeywordsIMPLE
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsIMPLE:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'm': this.StateID = ReadKeywordsIMPLEM
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsIMPLEM:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'e': this.StateID = ReadKeywordsIMPLEME
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsIMPLEME:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'n': this.StateID = ReadKeywordsIMPLEMEN
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsIMPLEMEN:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 't' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsIMPO:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'r': this.StateID = ReadKeywordsIMPOR
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsIMPOR:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 't' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsIN:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'h': this.StateID = ReadKeywordsINH
					case 't': this.StateID = ReadKeywordsINT
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'f' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsINH:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'e': this.StateID = ReadKeywordsINHE
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsINHE:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'r': this.StateID = ReadKeywordsINHER
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsINHER:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'i': this.StateID = ReadKeywordsINHERI
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsINHERI:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 't' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsINT:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'e': this.StateID = ReadKeywordsINTE
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsINTE:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'r': this.StateID = ReadKeywordsINTER
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsINTER:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'f': this.StateID = ReadKeywordsINTERF
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsINTERF:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'a': this.StateID = ReadKeywordsINTERFA
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsINTERFA:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'c': this.StateID = ReadKeywordsINTERFAC
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsINTERFAC:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'e' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsL:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'e': this.StateID = ReadKeywordsLE
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsLE:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 't' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsM:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'a': this.StateID = ReadKeywordsMA
					case 'o': this.StateID = ReadKeywordsMO
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsMA:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'p' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsMO:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'd': this.StateID = ReadKeywordsMOD
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsMOD:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'u': this.StateID = ReadKeywordsMODU
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsMODU:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'l': this.StateID = ReadKeywordsMODUL
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsMODUL:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'e' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsN:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'a': this.StateID = ReadKeywordsNA
					case 'o': this.StateID = ReadKeywordsNO
					case 'u': this.StateID = ReadKeywordsNU
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsNA:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'n' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsNO:
			if IsIdentChar(nextChar) {
				if thisChar == 't' && nextChar == 'h' {
					this.StateID = ReadKeywordsNOT
				} else {
					this.StateID = ReadIdent
				}
			} else {
				if thisChar == 't' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsNOT:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'h': this.StateID = ReadKeywordsNOTH
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsNOTH:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'i': this.StateID = ReadKeywordsNOTHI
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsNOTHI:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'n': this.StateID = ReadKeywordsNOTHIN
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsNOTHIN:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'g' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsNU:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'l': this.StateID = ReadKeywordsNUL
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsNUL:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'l' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsO:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'p': this.StateID = ReadKeywordsOP
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'r' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsOP:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'e': this.StateID = ReadKeywordsOPE
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsOPE:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'r' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsP:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'r': this.StateID = ReadKeywordsPR
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsPR:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'o': this.StateID = ReadKeywordsPRO
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsPRO:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'c' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsR:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'e': this.StateID = ReadKeywordsRE
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsRE:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'p': this.StateID = ReadKeywordsREP
					case 't': this.StateID = ReadKeywordsRET
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsREP:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'e': this.StateID = ReadKeywordsREPE
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsREPE:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'a': this.StateID = ReadKeywordsREPEA
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsREPEA:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 't' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsRET:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'u': this.StateID = ReadKeywordsRETU
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsRETU:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'r': this.StateID = ReadKeywordsRETUR
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsRETUR:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'n' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsS:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'e': this.StateID = ReadKeywordsSE
					case 'i': this.StateID = ReadKeywordsSI
					case 't': this.StateID = ReadKeywordsST
					case 'w': this.StateID = ReadKeywordsSW
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsSE:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'l': this.StateID = ReadKeywordsSEL
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsSEL:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'e': this.StateID = ReadKeywordsSELE
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsSELE:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'c': this.StateID = ReadKeywordsSELEC
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsSELEC:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 't' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsSI:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'z': this.StateID = ReadKeywordsSIZ
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsSIZ:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'e': this.StateID = ReadKeywordsSIZE
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsSIZE:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'o': this.StateID = ReadKeywordsSIZEO
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsSIZEO:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'f' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsST:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'r': this.StateID = ReadKeywordsSTR
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsSTR:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'u': this.StateID = ReadKeywordsSTRU
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsSTRU:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'c': this.StateID = ReadKeywordsSTRUC
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsSTRUC:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 't' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsSW:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'i': this.StateID = ReadKeywordsSWI
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsSWI:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 't': this.StateID = ReadKeywordsSWIT
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsSWIT:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'c': this.StateID = ReadKeywordsSWITC
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsSWITC:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'h' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsT:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'e': this.StateID = ReadKeywordsTE
					case 'h': this.StateID = ReadKeywordsTH
					case 'r': this.StateID = ReadKeywordsTR
					case 'y': this.StateID = ReadKeywordsTY
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsTE:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 's': this.StateID = ReadKeywordsTES
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsTES:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 't': this.StateID = ReadKeywordsTEST
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsTEST:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 's' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsTH:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'i': this.StateID = ReadKeywordsTHI
					case 'r': this.StateID = ReadKeywordsTHR
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsTHI:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 's' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsTHR:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'o': this.StateID = ReadKeywordsTHRO
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsTHRO:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'w' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsTR:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'u': this.StateID = ReadKeywordsTRU
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'y' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsTRU:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'e' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsTY:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'p': this.StateID = ReadKeywordsTYP
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsTYP:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'e': this.StateID = ReadKeywordsTYPE
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsTYPE:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'o': this.StateID = ReadKeywordsTYPEO
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsTYPEO:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'f' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsU:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'n': this.StateID = ReadKeywordsUN
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsUN:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 's': this.StateID = ReadKeywordsUNS
					case 't': this.StateID = ReadKeywordsUNT
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsUNS:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'a': this.StateID = ReadKeywordsUNSA
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsUNSA:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'f': this.StateID = ReadKeywordsUNSAF
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsUNSAF:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'e' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsUNT:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'i': this.StateID = ReadKeywordsUNTI
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsUNTI:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'l' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsV:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'a': this.StateID = ReadKeywordsVA
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsVA:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'r' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsW:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'h': this.StateID = ReadKeywordsWH
					case 'i': this.StateID = ReadKeywordsWI
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsWH:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'e': this.StateID = ReadKeywordsWHE
					case 'i': this.StateID = ReadKeywordsWHI
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsWHE:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'n' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsWHI:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'l': this.StateID = ReadKeywordsWHIL
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsWHIL:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'e' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsWI:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 't': this.StateID = ReadKeywordsWIT
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsWIT:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'h' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		case ReadKeywordsX:
			if IsKeywordChar(nextChar) {
				switch thisChar {
					case 'o': this.StateID = ReadKeywordsXO
					default:  this.StateID = ReadIdent
				}
			} else if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				this.StateID = AcceptIdent
			}

		case ReadKeywordsXO:
			if IsIdentChar(nextChar) {
				this.StateID = ReadIdent
			} else {
				if thisChar == 'r' {
					this.StateID = AcceptKeyword
				} else {
					this.StateID = AcceptIdent
				}
			}

		default:
			return Token {
				ID:       "error",
				RawValue: []byte(fmt.Sprintf("Unexpected character: %s", string([] byte {thisChar}))),
				Value:    fmt.Sprintf("Unexpected character: %s", string([] byte {thisChar})),

				FilePath: this.FilePath,
				Line:     this.Line,
				Column:   this.Column,

				Error: true,
			}, true
	}

	// Build temp token raw value
	tempToken.RawValue = append(tempToken.RawValue, thisChar)

	return return_token, result
}

func IsIdentChar(c byte) bool {
	return c == '_' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9'
}

func IsKeywordChar(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func IsNumChar(c byte) bool {
	return c >= '0' && c <= '9'
}

func IsHexChar(c byte) bool {
	return c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f'
}

func IsWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r'
}