// phpx-compiler: компилятор + режим --out для генерации бинарника

package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: phpxc файл.php [--out имя]")
		return
	}
	outFile := ""
	if len(os.Args) > 3 && os.Args[2] == "--out" {
		outFile = os.Args[3]
	}
	
	data, _ := os.ReadFile(os.Args[1])
	source := string(data)

	if strings.HasPrefix(source, "<?php") {
		source = strings.TrimPrefix(source, "<?php")
		source = strings.TrimLeft(source, " \n\r\t")
	}

	lexer := NewLexer(source)
	var tokens []Token
	for {
		tok := lexer.NextToken()
		if tok.Type == T_EOF {
			break
		}
		tokens = append(tokens, tok)
	}

	parser := NewParser(tokens)
	stmts := parser.Parse()

	for _, stmt := range stmts {
		CompileStmt(stmt)
	}

	bytecode = append(bytecode, OP_HALT)

	tmpFile := "vm_exec_temp.c"
	f, err := os.Create(tmpFile)
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpFile)
	defer f.Close()
	f.WriteString("#include <stdint.h>\n")
	f.WriteString("#include <stddef.h>\n")
	f.WriteString("#include \"vm.h\"\n\n")

	f.WriteString("uint8_t bytecode[] = {\n")
	for i, b := range bytecode {
		f.WriteString(fmt.Sprintf("0x%02X", b))
		if i != len(bytecode)-1 {
			f.WriteString(", ")
		}
	}
	f.WriteString("\n};\n")

	f.WriteString(fmt.Sprintf("\nsize_t bytecode_len = %d;\n\n", len(bytecode)))

	f.WriteString("Value constants[] = {\n")
	for _, c := range constants {
		if c.Type == "string" {
			f.WriteString(fmt.Sprintf("    {.type = TYPE_STRING, .value.str_val = \"%s\"},\n", c.Value))
		} else {
			num, _ := strconv.Atoi(c.Value)
			f.WriteString(fmt.Sprintf("    {.type = TYPE_INT, .value.int_val = %d},\n", num))
		}
	}
	f.WriteString("};\n")

	target := "vm_exec"
	if outFile != "" {
		target = outFile
	}
	cmd := exec.Command("gcc", "-DCOMPILE_AS_EXECUTABLE", "-I../phpc-vm", "-o", target, tmpFile, "../phpc-vm/vm.c")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Println("Ошибка компиляции VM")
		return
	}

	if outFile == "" {
		execCmd := exec.Command("./" + target)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		execCmd.Run()
		os.Remove(target)
	}
}