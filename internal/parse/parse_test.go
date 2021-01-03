package parse

import (
	"bytes"
	"errors"
	"testing"
)

func Test_parse(t *testing.T) {
	_, err := parser([]byte("+[->+<blahotherstuff]"))
	if err != nil {
		t.Fatalf("unexpected error (%s)", err.Error())
	}
}

func Test_parse_unbalancedLoop(t *testing.T) {
	_, err := parser([]byte("[->+<blahotherstuff[[]][[][][][][]]]]]]][[[]]][]["))
	if err == nil {
		t.Fatal("error expected on unbalanced loop set (none provided)")
	} else if !errors.Is(err, ErrorUnbalancedLoop) {
		t.Fatalf("error is of wrong type (got %s)", err.Error())
	}
}

func Test_parse_backwardsLoop(t *testing.T) {
	res, err := parser([]byte("+]-["))
	if err == nil {
		t.Fatal("error expected on backwards loop set (none provided)", string(res))
	} else if !errors.Is(err, ErrorUnbalancedLoop) {
		t.Fatalf("error is of wrong type (got %s)", err.Error())
	}
}

func Test_parse_removeInitialCommentLoop(t *testing.T) {
	inp := "[hello world +++++]]++"
	expected := "++"
	res, err := parser([]byte(inp))
	if err != nil {
		t.Fatalf("unexpected error (%s)", err.Error())
	} else if !bytes.Equal(res, []byte(expected)) {
		t.Fatalf("initial comment loop improperly removed (got %s, expected %s)", string(res), expected)
	}
}
