package applogger

import "context"

const (
	fieldTypeString = "string"
	fieldTypeInt    = "int"
	fieldTypeFloat  = "float"
	fieldTypeByte   = "byte"
	fieldTypeError  = "error"
	fieldTypeBool   = "bool"
)

type Field struct {
	Key       string
	String    string
	Int       int
	Float     float32
	Byte      []byte
	Bool      bool
	Error     error
	fieldType string
}

func String(key, value string) Field {
	return Field{
		Key:       key,
		String:    value,
		fieldType: fieldTypeString,
	}
}

func Int(key string, value int) Field {
	return Field{
		Key:       key,
		Int:       value,
		fieldType: fieldTypeInt,
	}
}

func Float(key string, value float32) Field {
	return Field{
		Key:       key,
		Float:     value,
		fieldType: fieldTypeFloat,
	}
}

func Byte(key string, value []byte) Field {
	return Field{
		Key:       key,
		Byte:      value,
		fieldType: fieldTypeFloat,
	}
}

func Bool(key string, value bool) Field {
	return Field{
		Key:       key,
		Bool:      value,
		fieldType: fieldTypeBool,
	}
}

func Error(err error) Field {
	return Field{
		Key:       "error",
		Error:     err,
		fieldType: fieldTypeError,
	}

}

type AppLogger interface {
	Info(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, fields ...Field)
}
