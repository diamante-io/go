package cmp

import "bufio"

type Scanner struct {
	*bufio.Scanner
	linesRead int
}

func (s *Scanner) Scan() bool {
	ret := s.Scanner.Scan()
	if ret {
		s.linesRead++
	}
	return ret
}

func (s *Scanner) LinesRead() int {
	return s.linesRead
}
