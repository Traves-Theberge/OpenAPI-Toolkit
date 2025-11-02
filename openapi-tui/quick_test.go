package main_test

import (
"testing"
"github.com/getkin/kin-openapi/openapi3"
)

func TestQuick(t *testing.T) {
responses := openapi3.NewResponses()
responses.Set("200", &openapi3.ResponseRef{
Value: &openapi3.Response{},
})

def := responses.Default()
t.Logf("Default: %v", def)

respMap := responses.Map()
for k, v := range respMap {
t.Logf("Response key=%s, val=%v", k, v)
}
}
