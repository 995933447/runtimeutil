package runtimeutil

import jsoniter "github.com/json-iterator/go"

var ForceNumberJson = jsoniter.Config{
	UseNumber: true,
}.Froze()
