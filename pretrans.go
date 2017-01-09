package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
  "regexp"
	"flag"
)

func main() {
	var fp *os.File
	var fpe *os.File
	var err error

  headPtr := flag.Bool("h", false, "extract headers")
  jpPtr := flag.Bool("j", false, "Japanese google txt to markdown")
	mixPtr := flag.Bool("m", false, "Mix English and Japanese")
	flag.Parse()

	nextArg := 0
	if *mixPtr == true {
		if len(flag.Args()) < 2 {
			panic("2 args required for -m option")
		}
		fpe, err = os.Open(flag.Args()[nextArg])
		nextArg++
		if err != nil {
			panic(err)
		}
		defer fpe.Close()
	}

	if len(flag.Args()) < 1 {
		fp = os.Stdin
	} else {
		fp, err = os.Open(flag.Args()[nextArg])
		nextArg++
		if err != nil {
			panic(err)
		}
		defer fp.Close()
	}
	if *mixPtr == true {
		mixEngilishAndJapanease(fpe, fp)
	}else	if *jpPtr == true {
		trimJapanease(fp)
	}else if *headPtr==true {
		extractHeaders(fp)
	}else{
		trimEnglish(fp)
	}
}

func mixEngilishAndJapanease(fpE *os.File, fp *os.File)  () {
	scannerE := bufio.NewScanner(fpE)
	scanner := bufio.NewScanner(fp)
	succeededE := true
	succeeded := true
	for ;; {
		for ;; {
			succeededE = scannerE.Scan()
			if succeededE == false {
				fmt.Println("")
				break
			}
			scanTextE := scannerE.Text()
			fmt.Println(scanTextE)
			if len(strings.TrimSpace(scanTextE)) == 0 {
				break;
			}
		}
		nonPrint := false
		for ;; {
			succeeded = scanner.Scan()
			if succeeded == false {
				fmt.Println("")
				break
			}
			scanText := scanner.Text()
			if nonPrint == false {
				nonPrint = isOneLine(scanText,"..")
			}
			if nonPrint == false {
				fmt.Println(scanText)
			}
			if len(strings.TrimSpace(scanText)) == 0 {
				break;
			}
		}
		if succeededE == false && succeeded == false {
			break;
		}
	}
}

func extractHeaders(fp *os.File)  () {
	scanner := bufio.NewScanner(fp)
  rep := regexp.MustCompile(`^\#+`)
	for scanner.Scan() {
		scanText := scanner.Text()
		matched := rep.MatchString(scanText)
		if matched==true {
			fmt.Println(scanText)
		}
	}
}

func trimJapanease(fp *os.File)  () {
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		scanText := scanner.Text()
    rep := regexp.MustCompile(`(^\*)([^ ])`)
    str := rep.ReplaceAllString(scanText, "$1 $2")
    rep = regexp.MustCompile(`(^\#+)([^\# ])`)
		str = rep.ReplaceAllString(str, "$1 $2")
		fmt.Println(str)
	}
}

func trimEnglish(fp *os.File)  () {
	scanner := bufio.NewScanner(fp)
  isLeft := true
  lastHyphenDeleted := false

	for scanner.Scan() {
    nlPrefix := false
    nlSuffix := false
    hyphenDeleted := false
		scanText := scanner.Text()
		hyphenDeleted, str := trimHyphen(scanText)
    nonHyphenText := str
    matched := false;

		/*
		matched, _ = regexp.MatchString("^[^0-9A-Za-z_ ]", nonHyphenText)
		if matched {
			nlPrefix = true
		}

		matched, _ = regexp.MatchString("[^0-9A-Za-z_ ]$", nonHyphenText)
		if matched {
      nlSuffix = true
		}
		*/

    if len(nonHyphenText) == 0 {
      nlPrefix = true
      nlSuffix = true
		}
		if isOneLine(nonHyphenText,"* ") {
//			nonHyphenText = nonHyphenText + "\n"
      nlPrefix = true
//      nlSuffix = true
    }

    if isOneLine(nonHyphenText,"..") {
//			nonHyphenText = strings.TrimPrefix(nonHyphenText, "..")
//			nonHyphenText = "<!--" + nonHyphenText + "-->"
      nlPrefix = true
      nlSuffix = true
    }

    if isOneLine(nonHyphenText,"#") {
      nlPrefix = true
      nlSuffix = true
    }

		// conversion
		if matched==false {
			matched, nonHyphenText = convertQuote(nonHyphenText, "| ")
		}
		if matched==false {
			matched, nonHyphenText = convertList(nonHyphenText, "> ")
		}
		if matched==false {
			matched, nonHyphenText = convertList(nonHyphenText, "• ")
		}
		if matched==false {
    	matched, nonHyphenText = convertNumlist(nonHyphenText)
		}
    if( matched==true ){
      nlPrefix = true
    }

    if isLeft == false && nlPrefix == false && lastHyphenDeleted == false {
			nonHyphenText = " " + nonHyphenText
		}
    // print
    if isLeft == false && nlPrefix {
      fmt.Print("\n")
    }
		fmt.Print(nonHyphenText)
    if nlSuffix == true {
      fmt.Print("\n")
      isLeft = true
    }else{
      isLeft = false
    }
    lastHyphenDeleted = hyphenDeleted
	}


	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func trimHyphen(x string)  (bool, string) {
	if strings.HasSuffix(x, "-") == true {
        return true, strings.TrimSuffix(x, "-")
	}
	return false, x
}

func convertNumlist(x string)  (bool, string) {
  	matched, _ := regexp.MatchString("^[0-9]+[.] ", x)
  if matched == true {
    return true, x
	}
	return false, x
}

func isOneLine(x string, prefix string) bool {
	if len(x) == 0 {
		return false
	}
	if strings.HasPrefix(x, prefix) == true {
		return true
	}
	return false
}

func convertList(x string, prefix string) (bool, string) {
	if strings.HasPrefix(x, prefix) == true {
		str := strings.TrimPrefix(x, prefix)
		str = "* " + str
		return true, str
	}
	return false, x
}

func convertQuote(x string, prefix string) (bool, string) {
	if strings.HasPrefix(x, prefix) == true {
		str := strings.TrimPrefix(x, prefix)
		str = "> " + str
		return true, str
	}
	return false, x
}
