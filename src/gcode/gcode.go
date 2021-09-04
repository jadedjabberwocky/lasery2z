package gcode

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/jadedjabberwocky/lasery2z/coordmap"
)

type param struct {
	name      string
	value     float64
	decPoints int
}

type instruction struct {
	raw     string
	cmd     string
	params  []param
	comment string
}

//var instructionRegex = regexp.MustCompile("^([GM][0-9]+ *)?([A-Z][-+\\.0-9]* *)(;.*)?$")
var instructionRegex = regexp.MustCompile(`^([GM][0-9]+)? *([^;]*)?(;.*)?$`)
var paramRegex = regexp.MustCompile(`([A-Z]-?[0-9]*\.?[0-9]*) *`)

// Errors
var (
	ErrBadCommand = errors.New("bad command")
)

func Process(r io.Reader, w io.Writer, cm coordmap.CoordMap) error {
	br := bufio.NewReader(r)

	mustBreak := false
	for {
		if mustBreak {
			return nil
		}

		line, err := br.ReadString('\n')
		if err == io.EOF {
			// Must still continue through the for loop
			// to handle the line that was read
			mustBreak = true
		} else if err != nil {
			return err
		}

		line = strings.Trim(line, "\n")
		instr, err := parse(line)
		if err != nil {
			return fmt.Errorf("%w (%v)", err, instr.raw)
		}

		switch instr.cmd {
		case "G0", "G1":
			if err := instr.addZifY(cm); err != nil {
				return err
			}
		}

		_, err = w.Write([]byte(instr.String()))
		if err != nil {
			return err
		}

		_, err = w.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}
}

func (instr *instruction) addZifY(cm coordmap.CoordMap) error {
	y, yindex := instr.getParam("Y")
	if yindex == -1 {
		return nil
	}

	z, zindex := instr.getParam("Z")
	dz, err := cm.Map(y.value)
	if err != nil {
		return err
	}

	z.value += dz
	if zindex == -1 {
		instr.params = append(instr.params, param{"Z", z.value, y.decPoints})
		return nil
	}

	instr.params[zindex].value = z.value
	return nil
}

func (instr *instruction) getParam(name string) (param, int) {
	for i, p := range instr.params {
		if p.name == name {
			return p, i
		}
	}
	return param{}, -1
}

func parse(line string) (*instruction, error) {
	m := instructionRegex.FindStringSubmatch(line)
	if m == nil {
		return &instruction{raw: line}, ErrBadCommand
	}

	params := m[2]
	pms := paramRegex.FindAllStringSubmatch(params, -1)
	paramList := []param(nil)
	for i := 0; i < len(pms); i++ {
		pm := pms[i][1]
		if pm == "" {
			continue
		}

		p, err := parseParam(pm)
		if err != nil {
			return &instruction{raw: line}, err
		}
		paramList = append(paramList, p)
	}

	return &instruction{
		raw:     line,
		cmd:     m[1],
		params:  paramList,
		comment: m[3],
	}, nil
}

func parseParam(p string) (param, error) {
	key := p[:1]
	value, err := strconv.ParseFloat(p[1:], 64)
	if err != nil {
		return param{}, ErrBadCommand
	}
	decPoint := strings.IndexRune(p, '.')
	if decPoint < 0 {
		return param{
			name:      key,
			value:     value,
			decPoints: 0,
		}, nil
	}

	return param{
		name:      key,
		value:     value,
		decPoints: len(p) - decPoint - 1,
	}, nil
}

func (i *instruction) String() string {
	sb := strings.Builder{}

	if i.cmd != "" {
		if sb.Len() > 0 {
			sb.WriteRune(' ')
		}
		sb.WriteString(i.cmd)
	}

	for _, param := range i.params {
		if sb.Len() > 0 {
			sb.WriteRune(' ')
		}
		format := fmt.Sprintf("%%s%%.%df", param.decPoints)
		sb.WriteString(fmt.Sprintf(format, param.name, param.value))
	}

	if i.comment != "" {
		if sb.Len() > 0 {
			sb.WriteRune(' ')
		}
		sb.WriteString(i.comment)
	}

	return sb.String()
}
