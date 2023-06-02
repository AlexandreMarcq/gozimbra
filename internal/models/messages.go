package models

import (
	"fmt"
	"io"
	"strings"

	"github.com/AlexandreMarcq/gozimbra/internal/utils"
)

type message interface {
	toResult() result
	write(w io.StringWriter) error
}

type getMsg struct {
	item  string
	attrs utils.AttrsMap
	err   error
}

func (m getMsg) toResult() result {
	return result{
		item: m.item,
		err:  m.err,
	}
}

func (m getMsg) write(w io.StringWriter) error {
	var sb strings.Builder
	_, err := sb.WriteString(m.item)
	if err != nil {
		return err
	}

	for _, attr := range m.attrs.Keys() {
		_, err = sb.WriteString(fmt.Sprintf(";%s", m.attrs[attr]))
		if err != nil {
			return err
		}
	}

	_, err = w.WriteString(fmt.Sprintf("%s\n", sb.String()))
	if err != nil {
		return err
	}

	return nil
}

func NewGetMsg(item string, attrs utils.AttrsMap, err error) getMsg {
	return getMsg{
		item,
		attrs,
		err,
	}
}

type modifyMsg struct {
	item     string
	oldAttrs utils.AttrsMap
	newAttrs utils.AttrsMap
	err      error
}

func (m modifyMsg) toResult() result {
	return result{
		item: m.item,
		err:  m.err,
	}
}

func (m modifyMsg) write(w io.StringWriter) error {
	var sb strings.Builder
	_, err := sb.WriteString(m.item)
	if err != nil {
		return err
	}

	for _, k := range m.oldAttrs.Keys() {
		var oldString, newString string
		if m.err != nil {
			oldString = "ERR"
			newString = "ERR"
		} else {
			oldString = m.oldAttrs[k]
			newString = m.newAttrs[k]
		}

		_, err := sb.WriteString(fmt.Sprintf(";%s;%s", oldString, newString))
		if err != nil {
			return err
		}
	}

	_, err = w.WriteString(fmt.Sprintf("%s\n", sb.String()))
	if err != nil {
		return err
	}

	return nil
}

func NewModifyMsg(item string, oldAttrs, newAttrs utils.AttrsMap, err error) modifyMsg {
	return modifyMsg{
		item,
		oldAttrs,
		newAttrs,
		err,
	}
}
