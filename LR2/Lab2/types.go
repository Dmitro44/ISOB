package main

import (
	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
)

type (
	Method    int
	Mode      int
	InputMode int
)

const (
	ManualInput InputMode = iota
	File
)

const (
	Caesar Method = iota
	Vigenere
)

const (
	Encrypt Mode = iota
	Decrypt
)

const (
	FocusInputMode int = iota
	FocusFilepicker
	FocusMethod
	FocusMode
	FocusKey
	FocusShift
	FocusInput
	FocusBtn
	FocusSaveBtn
	FocusCount
)

type model struct {
	Filepicker    filepicker.Model
	SelectedFile  string
	FileConfirmed bool
	InputMode     InputMode
	Method        Method
	Mode          Mode
	Status        string
	Focus         int
	Width         int
	Height        int
	Key           textinput.Model
	Shift         textinput.Model
	Input         textarea.Model
	Output        textarea.Model
	SaveFilename  textinput.Model
	ShowSavePopup bool
	keymap        keymap
}

type keymap = struct {
	next, prev, quit, copy, clearInput key.Binding
}

type statusMsg string
