package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/0xMoonrise/asql/internal/kernel"
	"github.com/chzyer/readline"
)

func run() error {
	k := kernel.NewKernel("db.json")

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "asql> ",
		HistoryFile:     os.TempDir() + "/asql_history",
		InterruptPrompt: "^C",
		EOFPrompt:       ".exit",
	})

	if err != nil {
		return err
	}

	defer rl.Close()

	fmt.Println("Mini SQL shell. Asql 0.1v")
	var sql strings.Builder

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		if strings.TrimSpace(line) == ".exit" {
			break
		}

		if sql.Len() == 0 && strings.TrimSpace(line) == "" {
			continue
		}

		sql.WriteString(line)
		sql.WriteString("\n")

		if strings.Contains(line, ";") {
			stmt := strings.TrimSpace(sql.String())
			if err := k.ProcessSQL(stmt); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			}
			sql.Reset()
			rl.SetPrompt("asql> ")
		} else {
			rl.SetPrompt("   -> ")
		}
	}

	fmt.Println("bye")
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
