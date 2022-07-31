package jsonchain

import (
	"testing"
)

func TestNewJsonChain(t *testing.T) {
	chain := NewJsonChain[string, string]()
	if chain == nil {
		t.Error("NewJsonChain returned nil")
	}
}

func TestJsonChain_Set(t *testing.T) {
	chain := NewJsonChain[string, string]()
	if chain.Set("name", "Mattéo").Set("belovedLanguage", "Go") == nil {
		t.Error("JsonChain.Set returned nil")
		return
	}
	if chain.Get("name") != "Mattéo" && chain.Get("belovedLanguage") != "Go" {
		t.Error("JsonChain.Set failed")
	}
}

func TestJsonChain_SetWithError(t *testing.T) {
	chain := NewJsonChain[string, string]()
	if err := chain.SetWithError("name", "Mattéo"); err != nil {
		t.Error("JsonChain.SetWithError failed")
	}
}

func TestJsonChain_Get(t *testing.T) {
	chain := NewJsonChain[string, string]()
	chain.Set("name", "Mattéo")
	if chain.Get("name") != "Mattéo" {
		t.Error("JsonChain.Get failed")
	}
}

func TestJsonChain_AppendFromBytes(t *testing.T) {
	chain := NewJsonChain[string, string]()
	if err := chain.AppendFromBytes([]byte(`{"name":"Mattéo"}`)); err != nil {
		t.Error("JsonChain.Append failed")
	}
	value, err := chain.GetWithError("name")
	if err != nil {
		t.Error("JsonChain.GetWithError failed")
		return
	}
	if value != "Mattéo" {
		t.Error("JsonChain.GetWithError failed")
	}
}

func TestJsonChain_Clone(t *testing.T) {
	chain := NewJsonChain[string, string]()
	chain.Set("name", "Mattéo")
	clone := chain.Clone()
	if !chain.Equal(clone) {
		t.Error("JsonChain.Clone failed")
		return
	}
	clone.Set("belovedLanguage", "Go")
	if chain.Equal(clone) {
		t.Error("JsonChain.Clone failed")
	}
}
