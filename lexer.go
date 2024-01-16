package main

import (
	"fmt"
	"strconv"
	"unicode"
)

type enumType string

const (
	x          = "X"
	digit      = "Digit"
	operator   = "Operator"
	whitespace = "Whitespace"
)

type Token struct {
	valueString string
	valueNumber float64
	tokenType   enumType
}

type invalidDuplicateOp struct{}

func (m *invalidDuplicateOp) Error() string {
	return "invalid duplicate operator"
}

func lexerFunction(argument string) ([]Token, error) {
	var tokensTab []Token
	var tokenStart int = 0
	var currentType enumType = whitespace

	for i, letter := range argument {
		letterType := getTokenType(letter)
		if currentType != letterType || currentType != digit {
			tmpToken, err := newToken(argument[tokenStart:i], currentType)
			if err != nil {
				return tokensTab, err
			}
			tokensTab = append(tokensTab, tmpToken)
			tokenStart = i
			currentType = letterType
		}
	}
	tmpToken, err := newToken(argument[tokenStart:], currentType)
	if err != nil {
		return tokensTab, err
	}
	tokensTab = append(tokensTab, tmpToken)
	tokensTab = removeWhitespaces(tokensTab)
	// fmt.Println(tokensTab)
	tokensTab, err = cleanToken(tokensTab)
	if err != nil {
		return tokensTab, err
	}
	// fmt.Println(tokensTab)
	return tokensTab, nil
}

func removeWhitespaces(tab []Token) []Token {
	var newTab []Token

	for _, token := range tab {
		if token.tokenType == whitespace {
			continue
		}
		newTab = append(newTab, token)
	}
	return newTab
}

func getTokenType(c rune) enumType {

	if unicode.IsDigit(c) || c == '.' {
		return digit
	} else if c == 'x' || c == 'X' {
		return x
	} else if unicode.IsSpace(c) {
		return whitespace
	}
	return operator
}

func newToken(str string, currentType enumType) (Token, error) {

	var token Token
	var err error
	token.valueString = str
	if currentType == digit {
		token.valueString = ""
		token.valueNumber, err = strconv.ParseFloat(str, 64)
		if err != nil {
			return token, err
		}
	}
	token.tokenType = currentType
	return token, nil
}

func concatOperator(tab []Token, oper int, nb int) (Token, int) {
	if tab[oper].tokenType == operator && tab[nb].tokenType == digit {
		if tab[oper].valueString == "-" {
			newToken := Token{
				valueNumber: tab[nb].valueNumber * -1,
				tokenType:   digit,
			}
			return newToken, 2
		}
		if tab[oper].valueString == "+" {
			newToken := Token{
				valueNumber: tab[nb].valueNumber,
				tokenType:   digit,
			}
			return newToken, 2
		}
	}
	return tab[oper], 1
}

func changeOperator(oper string) Token {

	newOperator := Token{
		valueString: oper,
		tokenType:   operator,
	}
	return newOperator
}

func cleanToken(tab []Token) ([]Token, error) {

	var newTok Token
	var newTab []Token
	var i int

	newTok, i = concatOperator(tab, 0, 1)
	newTab = append(newTab, newTok)
	for i < len(tab) {
		var token Token
		token = tab[i]
		if i != len(tab)-1 && token.tokenType == operator && tab[i+1].tokenType == operator {
			firstOperator := token.valueString
			secondOperator := tab[i+1].valueString
			if (firstOperator == "+" || firstOperator == "-") &&
				(secondOperator == "+" || secondOperator == "-") {
				if (firstOperator == "-" && firstOperator == secondOperator) ||
					(firstOperator == "+" && firstOperator == secondOperator) {
					newOperator := changeOperator("+")
					newTab = append(newTab, newOperator)
				} else if firstOperator != secondOperator {
					newOperator := changeOperator("-")
					newTab = append(newTab, newOperator)
				}
				i = i + 1
			} else if (firstOperator == "*" || firstOperator == "/") &&
				(secondOperator == "+" || secondOperator == "-") {
				var reduceOp Token
				reduceOp, _ = concatOperator(tab, i+1, i+2)
				newTab = append(newTab, token)
				newTab = append(newTab, reduceOp)
				i = i + 2
			} else {
				return newTab, &invalidDuplicateOp{}
			}
			fmt.Println(token)
			i = i + 1
		} else if token.tokenType == x {
			if i != 0 && tab[i-1].tokenType == digit {
				newOperator := changeOperator("*")
				newTab = append(newTab, newOperator)
				newTab = append(newTab, token)
				i = i + 1
			}
		} else {
			newTab = append(newTab, token)
			i = i + 1
		}
	}
	return newTab, nil
}
