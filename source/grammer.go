package main

import "strings"

func trimQuotes(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func splitArgs(line string) []string {
	out := []string{}
	cur := strings.Builder{}
	inQuotes := false
	var q byte
	for i := 0; i < len(line); i++ {
		c := line[i]
		if inQuotes {
			if c == q {
				inQuotes = false
				out = append(out, cur.String())
				cur.Reset()
			} else {
				cur.WriteByte(c)
			}
			continue
		}
		if c == '"' || c == '\'' {
			inQuotes = true
			q = c
			continue
		}
		if c == ' ' || c == '\t' {
			if cur.Len() > 0 {
				out = append(out, cur.String())
				cur.Reset()
			}
			continue
		}
		cur.WriteByte(c)
	}
	if cur.Len() > 0 {
		out = append(out, cur.String())
	}
	return out
}
