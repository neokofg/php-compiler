// PHP Compiler - compiles php code to IR and then running it on PHPC VM
// Copyright (C) 2025  Andrey Vasilev (neokofg)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
package main

import (
	"fmt"
	"github.com/neokofg/php-compiler/internal/compiler"
	"github.com/neokofg/php-compiler/internal/lexer"
	"github.com/neokofg/php-compiler/internal/parser"
	"github.com/neokofg/php-compiler/internal/token"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	source, outFile, err := processArgs()
	if err != nil {
		fmt.Println(err)
		return
	}

	tokens, err := lexicalAnalysis(source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	phpCompiler, err := parseAndCompile(tokens)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	tmpFile := "vm_exec_temp.c"
	if err := generateVMCode(phpCompiler, tmpFile); err != nil {
		panic(err)
	}
	defer os.Remove(tmpFile)

	if err := compileAndRunVM(tmpFile, outFile); err != nil {
		fmt.Println(err)
	}
}

func processArgs() (string, string, error) {
	if len(os.Args) < 2 {
		return "", "", fmt.Errorf("Usage: phpc file.php [--out name]")
	}

	outFile := ""
	if len(os.Args) > 3 && os.Args[2] == "--out" {
		outFile = os.Args[3]
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		return "", "", fmt.Errorf("File reading error %s: %v", os.Args[1], err)
	}

	source := string(data)
	if strings.HasPrefix(source, "<?php") {
		source = strings.TrimPrefix(source, "<?php")
		source = strings.TrimLeft(source, " \n\r\t")
	}

	return source, outFile, nil
}

func lexicalAnalysis(source string) ([]token.Token, error) {
	lexerInstance := lexer.NewLexer(source)
	var tokens []token.Token

	for {
		tok := lexerInstance.NextToken()
		if tok.Type == token.T_ILLEGAL {
			return nil, fmt.Errorf("Lexer analyze error: %s", tok.Value)
		}
		if tok.Type == token.T_EOF {
			break
		}
		tokens = append(tokens, tok)
	}

	return tokens, nil
}

func parseAndCompile(tokens []token.Token) (*compiler.Compiler, error) {
	parserInstance := parser.NewParser(tokens)
	stmts, err := parserInstance.Parse()
	if err != nil {
		return nil, fmt.Errorf("Syntax analyze error: %v", err)
	}

	phpCompiler := compiler.New()
	err = phpCompiler.CompileProgram(stmts)
	if err != nil {
		return nil, fmt.Errorf("Compilation error: %v", err)
	}

	return phpCompiler, nil
}

func generateVMCode(phpCompiler *compiler.Compiler, tmpFile string) error {
	bytecode := phpCompiler.GetBytecode()
	constants := phpCompiler.GetConstants()

	f, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("Cannot create temp c file %s: %v", tmpFile, err)
	}
	defer f.Close()

	headers := []string{
		"#include <stdint.h>\n",
		"#include <stddef.h>\n",
		"#include <stdio.h>\n",
		"#include \"vm.h\"\n\n",
	}

	for _, header := range headers {
		if _, err := f.WriteString(header); err != nil {
			return fmt.Errorf("Error writing header in %s: %v", tmpFile, err)
		}
	}

	if _, err := f.WriteString("uint8_t bytecode[] = {\n    "); err != nil {
		return fmt.Errorf("Error writing bytecode start in %s: %v", tmpFile, err)
	}

	for i, b := range bytecode {
		if _, err := f.WriteString(fmt.Sprintf("0x%02X", b)); err != nil {
			return fmt.Errorf("Error writing bytecode in %s: %v", tmpFile, err)
		}

		if i != len(bytecode)-1 {
			if _, err := f.WriteString(", "); err != nil {
				return fmt.Errorf("Error writing bytecode separator in %s: %v", tmpFile, err)
			}

			if (i+1)%16 == 0 {
				if _, err := f.WriteString("\n    "); err != nil {
					return fmt.Errorf("Error writing bytecode newline in %s: %v", tmpFile, err)
				}
			}
		}
	}

	if _, err := f.WriteString("\n};\n"); err != nil {
		return fmt.Errorf("Error writing bytecode end in %s: %v", tmpFile, err)
	}

	if _, err := f.WriteString(fmt.Sprintf("\nsize_t bytecode_len = %d;\n\n", len(bytecode))); err != nil {
		return fmt.Errorf("Error writing bytecode_len in %s: %v", tmpFile, err)
	}

	if _, err := f.WriteString("Value constants[] = {\n"); err != nil {
		return fmt.Errorf("Error writing constants start in %s: %v", tmpFile, err)
	}

	for _, c := range constants {
		if c.Type == "string" {
			if _, err := f.WriteString(fmt.Sprintf("    {.type = TYPE_STRING, .value.str_val = %s},\n",
				strconv.Quote(c.Value))); err != nil {
				return fmt.Errorf("Error writing string constant in %s: %v", tmpFile, err)
			}
		} else {
			num, errAtoi := strconv.Atoi(c.Value)
			if errAtoi != nil {
				return fmt.Errorf("Critical error: Cannot convert const '%s' to integer: %v", c.Value, errAtoi)
			}

			if _, err := f.WriteString(fmt.Sprintf("    {.type = TYPE_INT, .value.int_val = %d},\n", num)); err != nil {
				return fmt.Errorf("Error writing int constant in %s: %v", tmpFile, err)
			}
		}
	}

	if _, err := f.WriteString("};\n"); err != nil {
		return fmt.Errorf("Error writing constants end in %s: %v", tmpFile, err)
	}

	if _, err := f.WriteString(fmt.Sprintf("\nsize_t constants_len = %d;\n\n", len(constants))); err != nil {
		return fmt.Errorf("Error writing constants_len in %s: %v", tmpFile, err)
	}

	mainFunctionCode := []string{
		"int main(int argc, char** argv) {\n",
		"    VM* vm = vm_new();\n",
		"    if (!vm) {\n",
		"        fprintf(stderr, \"Failed to create VM\\n\");\n",
		"        return 1;\n",
		"    }\n\n",
		"    if (argc > 1 && strcmp(argv[1], \"--debug\") == 0) {\n",
		"        vm_set_debug_mode(vm, true);\n",
		"    }\n\n",
		"    status_t status = vm_execute(vm, bytecode, bytecode_len, constants, constants_len);\n",
		"    vm_free(vm);\n\n",
		"    return (status == STATUS_SUCCESS) ? 0 : 1;\n",
		"}\n",
	}

	for _, line := range mainFunctionCode {
		if _, err := f.WriteString(line); err != nil {
			return fmt.Errorf("Error writing main function in %s: %v", tmpFile, err)
		}
	}

	return nil
}

func compileAndRunVM(tmpFile string, outFile string) error {
	target := "vm_exec"
	if outFile != "" {
		target = outFile
	}

	vmDir := "vm"
	cFlags := "-DCOMPILE_AS_EXECUTABLE"
	includeFlag := "-I" + vmDir + "/includes"

	cmd := exec.Command("gcc", cFlags, includeFlag, "-o", target,
		tmpFile,
		vmDir+"/src/core/vm.c",
		vmDir+"/src/core/context.c",
		vmDir+"/src/core/dispatcher.c",
		vmDir+"/src/handlers/arithmetic.c",
		vmDir+"/src/handlers/core.c",
		vmDir+"/src/handlers/flow.c",
		vmDir+"/src/handlers/logic.c",
		vmDir+"/src/handlers/string.c",
		vmDir+"/src/handlers/function.c",
		vmDir+"/src/components/value.c",
		vmDir+"/src/components/memory.c",
		vmDir+"/src/components/stack.c",
		vmDir+"/src/components/error.c")

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	//fmt.Printf("Compiling C with command: %s\n", cmd.String()) // cmd string debug

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("VM compilation error, please report this error to repository owners: %v", err)
	}

	if outFile == "" {
		execCmd := exec.Command("./" + target)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr

		if runErr := execCmd.Run(); runErr != nil {
			return fmt.Errorf("VM runtime error, please report this error to repository owners: %v", runErr)
		}

		os.Remove(target)
	}

	return nil
}
