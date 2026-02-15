package envx

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func parseDotEnv(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	values := make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		key = strings.TrimPrefix(key, "export ")
		key = strings.TrimSpace(key)
		if key == "" {
			return nil, fmt.Errorf("invalid dotenv line: %q", line)
		}

		value := strings.TrimSpace(parts[1])
		value = stripInlineComment(value)

		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		values[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return values, nil
}

func stripInlineComment(value string) string {
	inSingle := false
	inDouble := false

	for i, r := range value {
		switch r {
		case '\'':
			if !inDouble {
				inSingle = !inSingle
			}
		case '"':
			if !inSingle {
				inDouble = !inDouble
			}
		case '#':
			if inSingle || inDouble {
				continue
			}

			if i == 0 {
				return ""
			}

			prev := []rune(value[:i])
			if len(prev) == 0 {
				return ""
			}
			if unicode.IsSpace(prev[len(prev)-1]) {
				return strings.TrimSpace(value[:i])
			}
		}
	}

	return strings.TrimSpace(value)
}
