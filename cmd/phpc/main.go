// phpx-compiler: компилятор + режим --out для генерации бинарника

package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"os/exec"
	"github.com/neokofg/php-compiler/internal/lexer"
	"github.com/neokofg/php-compiler/internal/parser"
	"github.com/neokofg/php-compiler/internal/token"
	"github.com/neokofg/php-compiler/internal/compiler"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: phpc file.php [--out name]")
		return
	}
	outFile := ""
	if len(os.Args) > 3 && os.Args[2] == "--out" {
		outFile = os.Args[3]
	}
	
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("File reading error %s: %v\n", os.Args[1], err)
		return
	}
	source := string(data)

	if strings.HasPrefix(source, "<?php") {
		source = strings.TrimPrefix(source, "<?php")
		source = strings.TrimLeft(source, " \n\r\t")
	}

	lexerInstance := lexer.NewLexer(source)
	var tokens []token.Token
	for {
		tok := lexerInstance.NextToken()
		if tok.Type == token.T_ILLEGAL {
			fmt.Fprintf(os.Stderr, "Lexer analyze error: %s\n", tok.Value)
			os.Exit(1)
		}
		if tok.Type == token.T_EOF {
			break
		}
		tokens = append(tokens, tok)
	}	

	parserInstance := parser.NewParser(tokens)
	stmts, err := parserInstance.Parse()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Syntax analyze error: %v\n", err)
        os.Exit(1)
    }

	for _, stmt := range stmts {
		compiler.CompileStmt(stmt)
	}

	compiler.Bytecode = append(compiler.Bytecode, compiler.OP_HALT)

	tmpFile := "vm_exec_temp.c"
	f, err := os.Create(tmpFile)
	if err != nil {
		panic(fmt.Sprintf("Cannot create temp c file %s: %v", tmpFile, err))
	}
	defer f.Close()
	defer os.Remove(tmpFile)
	
	_, err = f.WriteString("#include <stdint.h>\n")
	if err != nil {
		panic(fmt.Sprintf("Error writing in %s: %v", tmpFile, err))
	}
	_, err = f.WriteString("#include <stddef.h>\n")
	if err != nil {
		panic(fmt.Sprintf("Error writing in %s: %v", tmpFile, err))
	}
	_, err = f.WriteString("#include \"vm.h\"\n\n")
	if err != nil {
		panic(fmt.Sprintf("Error writing in %s: %v", tmpFile, err))
	}

	_, err = f.WriteString("uint8_t bytecode[] = {\n    ")
	if err != nil {
		panic(fmt.Sprintf("Error writing in %s: %v", tmpFile, err))
	}
	for i, b := range compiler.Bytecode {
		_, err = f.WriteString(fmt.Sprintf("0x%02X", b))
		if err != nil {
			panic(fmt.Sprintf("Error writing in %s: %v", tmpFile, err))
		}
		if i != len(compiler.Bytecode)-1 {
			_, err = f.WriteString(", ")
			if err != nil {
				panic(fmt.Sprintf("Error writing in %s: %v", tmpFile, err))
			}
			if (i+1)%16 == 0 {
				_, err = f.WriteString("\n    ")
				if err != nil {
					panic(fmt.Sprintf("Error writing in %s: %v", tmpFile, err))
				}
			}
		}
	}
	_, err = f.WriteString("\n};\n")
	if err != nil {
		panic(fmt.Sprintf("Error writing in %s: %v", tmpFile, err))
	}

	_, err = f.WriteString(fmt.Sprintf("\nsize_t bytecode_len = %d;\n\n", len(compiler.Bytecode)))
	if err != nil {
		panic(fmt.Sprintf("Error writing in %s: %v", tmpFile, err))
	}

	_, err = f.WriteString("Value constants[] = {\n")
	if err != nil {
		panic(fmt.Sprintf("Error writing in %s: %v", tmpFile, err))
	}
	for _, c := range compiler.Constants {
		if c.Type == "string" {
			_, err = f.WriteString(fmt.Sprintf("    {.type = TYPE_STRING, .value.str_val = %s},\n", strconv.Quote(c.Value)))
			if err != nil {
				panic(fmt.Sprintf("Error writing string in %s: %v", tmpFile, err))
			}
		} else {
			num, errAtoi := strconv.Atoi(c.Value)
			if errAtoi != nil {
				panic(fmt.Sprintf("Critical error: Cannot convert const '%s' to integer: %v", c.Value, errAtoi))
			}
			_, err = f.WriteString(fmt.Sprintf("    {.type = TYPE_INT, .value.int_val = %d},\n", num))
			if err != nil {
				panic(fmt.Sprintf("Error writing int in %s: %v", tmpFile, err))
			}
		}
	}
	_, err = f.WriteString("};\n")
	if err != nil {
		panic(fmt.Sprintf("Error writing in %s: %v", tmpFile, err))
	}

	_, err = f.WriteString(fmt.Sprintf("\nsize_t constants_len = %d;\n\n", len(compiler.Constants)))
    if err != nil {
        panic(fmt.Sprintf("Error writing constants_len in %s: %v", tmpFile, err))
    }

	err = f.Close()
	if err != nil {
		panic(fmt.Sprintf("Error closing file %s: %v", tmpFile, err))
	}

	target := "vm_exec"
	if outFile != "" {
		target = outFile
	}
	cmd := exec.Command("gcc", "-DCOMPILE_AS_EXECUTABLE", "-Ivm", "-o", target, tmpFile, "vm/vm.c")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	fmt.Printf("Compiling C with command: %s\n", cmd.String())
	err = cmd.Run()
	if err != nil {
		fmt.Println("VM compilation error, please report this error to repository owners")
		return
	}
	if outFile == "" {
		fmt.Println("--- Starting VM ---")
		execCmd := exec.Command("./" + target)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		runErr := execCmd.Run()
		if runErr != nil {
			fmt.Println("VM runtime error, please report this error to repository owners:", runErr)
		}
		fmt.Println("--- Ending VM ---")
		os.Remove(target)
	}
}