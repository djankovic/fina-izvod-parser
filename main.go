package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type _900 struct {
	IZ900DVBDIPOS string `decoder:"7"`
	IZ900NAZBAN   string `decoder:"50"`
	IZ9000OIBBNK  string `decoder:"11"`
	IZ900VRIZ     string `decoder:"4"`
	IZ900DATUM    string `decoder:"8"`
	IZ900REZ2     string `decoder:"917" json:"-"`
	IZ900TIPSL    string `decoder:"3" json:"-"`
}

type _903 struct {
	IZ903VBDI    string `decoder:"7"`
	IZ903BIC     string `decoder:"11"`
	IZ903RACUN   string `decoder:"21"`
	IZ903VLRN    string `decoder:"3"`
	IZ903NAZKLI  string `decoder:"70"`
	IZ903SJEDKLI string `decoder:"35"`
	IZ903MB      string `decoder:"8"`
	IZ903OIBKLI  string `decoder:"11"`
	IZ903RBIZV   string `decoder:"3"`
	IZ903PODBR   string `decoder:"3"`
	IZ903DATUM   string `decoder:"8"`
	IZ903BRGRU   string `decoder:"4"`
	IZ903VRIZ    string `decoder:"4"`
	IZ903REZ     string `decoder:"809" json:"-"`
	IZ903TIPSL   string `decoder:"3" json:"-"`
}

type _905 struct {
	IZ905OZTRA         string `decoder:"2"`
	IZ905RNPRPL        string `decoder:"34"`
	IZ905NAZPRPL       string `decoder:"70"`
	IZ905ADRPRPL       string `decoder:"35"`
	IZ905SJPRPL        string `decoder:"35"`
	IZ905DATVAL        string `decoder:"8"`
	IZ905DATIZVR       string `decoder:"8"`
	IZ905VLPL          string `decoder:"3"`
	IZ905TECAJ         string `decoder:"15"`
	IZ905PREDZNVL      string `decoder:"1"`
	IZ905IZNOSPPVALUTE string `decoder:"15"`
	IZ905PREDZN        string `decoder:"1"`
	IZ905IZNOS         string `decoder:"15"`
	IZ905PNBPL         string `decoder:"26"`
	IZ905PNBPR         string `decoder:"26"`
	IZ905SIFNAM        string `decoder:"4"`
	IZ905OPISPL        string `decoder:"140"`
	IZ905IDTRFINA      string `decoder:"42"`
	IZ905IDTRBAN       string `decoder:"35"`
	IZ905REZ2          string `decoder:"482" json:"-"`
	IZ905TIPSL         string `decoder:"3" json:"-"`
}

type _907 struct {
	IZ907RAČUN    string `decoder:"21"`
	IZ907VLRN     string `decoder:"3"`
	IZ907NAZKLI   string `decoder:"70"`
	IZ907RBIZV    string `decoder:"3"`
	IZ907PRRBIZV  string `decoder:"3"`
	IZ907DATUM    string `decoder:"8"`
	IZ907DATPRSAL string `decoder:"8"`
	IZ907PPPOS    string `decoder:"1"`
	IZ907PRSAL    string `decoder:"15"`
	IZ907PREREZ   string `decoder:"1"`
	IZ907IZNREZ   string `decoder:"15"`
	IZ907DATOKV   string `decoder:"8"`
	IZ907IZNOKV   string `decoder:"15"`
	IZ907IZNZAPSR string `decoder:"15"`
	IZ907PRASPSTA string `decoder:"1"`
	IZ907IZNRASP  string `decoder:"15"`
	IZ907PDUGU    string `decoder:"1"`
	IZ907KDUGU    string `decoder:"15"`
	IZ907PPOTR    string `decoder:"1"`
	IZ907KPOTR    string `decoder:"15"`
	IZ07PRNOS     string `decoder:"1"`
	IZ907KOSAL    string `decoder:"15"`
	IZ907BRGRU    string `decoder:"4"`
	IZ907BRSTA    string `decoder:"6"`
	IZ907TEKST    string `decoder:"420"`
	IZ907REZ2     string `decoder:"317" json:"-"`
	IZ907TIPSL    string `decoder:"3" json:"-"`
}

type _909 struct {
	IZ909DATUM string `decoder:"8"`
	IZ909UKGRU string `decoder:"5"`
	IZ909UKSLG string `decoder:"6"`
	IZ909REZ3  string `decoder:"978" json:"-"`
	IZ909TIPSL string `decoder:"3" json:"-"`
}

type _999 struct {
	IZ999REZ1  string `decoder:"997" json:"-"`
	IZ999TIPSL string `decoder:"3" json:"-"`
}

// Group repesents a group of transactions in a single Statement
type Group struct {
	*_903
	*_907
	Transactions []_905
}

// Statement represents a single bank statement
type Statement struct {
	*_900
	*_909
	Groups []Group
}

func decode(str interface{}, runes []rune) {
	reflected := reflect.Indirect(reflect.ValueOf(str))

	for i := 0; i < reflected.NumField(); i++ {
		tag := reflected.Type().Field(i).Tag.Get("decoder")
		fieldLen, _ := strconv.ParseInt(tag, 10, 64)
		reflected.Field(i).SetString(strings.TrimSpace(string(runes[:fieldLen])))
		runes = runes[fieldLen:]
	}
}

func main() {
	decoder := charmap.Windows1250.NewDecoder()
	stat, _ := os.Stdin.Stat()
	hasStdin := (stat.Mode() & os.ModeCharDevice) == 0
	hasFile := len(os.Args) >= 2

	if hasFile == hasStdin {
		fmt.Println("Provide either the filename or stdin input")
		os.Exit(1)
	}

	var decoderReader io.Reader
	if hasFile {
		f, _ := os.Open(os.Args[1])
		defer f.Close()
		decoderReader = decoder.Reader(f)
	} else {
		decoderReader = decoder.Reader(os.Stdin)
	}

	reader := bufio.NewReader(decoderReader)
	linesScanned := 0
	statement := Statement{_900: &_900{}, _909: &_909{}, Groups: []Group{}}

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		runes := []rune(line)
		runes = runes[:1000]
		strType := string(runes[997:])

		if linesScanned == 0 && strType != "900" {
			fmt.Println(fmt.Sprintf("Invalid line %d with type %s", linesScanned, strType))
			os.Exit(1)
		}

		switch strType {
		case "900":
			decode(statement._900, runes)
			break
		case "903":
			str := _903{}
			decode(&str, runes)
			statement.Groups = append(statement.Groups, Group{_903: &str, Transactions: []_905{}})
			break
		case "905":
			str := _905{}
			decode(&str, runes)
			iz905 := &statement.Groups[len(statement.Groups)-1].Transactions
			*iz905 = append(*iz905, str)
			break
		case "907":
			str := _907{}
			decode(&str, runes)
			statement.Groups[len(statement.Groups)-1]._907 = &str
			break
		case "909":
			decode(statement._909, runes)
			break
		case "999":
			break
		default:
			fmt.Println(fmt.Sprintf("Unsupported type %s", strType))
			os.Exit(1)
		}

		linesScanned++
	}

	json, _ := json.Marshal(statement)
	fmt.Println(string(json))
}
