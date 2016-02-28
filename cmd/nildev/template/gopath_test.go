package template

import (
	"github.com/nildev/lib/utils"
	. "gopkg.in/check.v1"
	"io/ioutil"
)

type GopathSuite struct{}

var _ = Suite(&GopathSuite{})

func (s *GopathSuite) TestIfTemplatesAreFound(c *C) {
	gpl := NewGoPathLoader()
	data, err := gpl.Load("nildev", "test-good", "v0.1.9")
	c.Assert(err, IsNil)
	expected, err := ioutil.ReadFile("./fixtures/good.tpl")
	c.Assert(err, IsNil)
	remaining, err := utils.PopLine(expected)
	c.Assert(err, IsNil)
	c.Assert(string(data), Equals, string(remaining))
}

func (s *GopathSuite) TestIfTemplateIsFoundWhenVersionIsIgnored(c *C) {
	gpl := NewGoPathLoader()
	data, err := gpl.Load("nildev", "test-good", "")
	c.Assert(err, IsNil)
	expected, err := ioutil.ReadFile("./fixtures/good.tpl")
	c.Assert(err, IsNil)
	remaining, err := utils.PopLine(expected)
	c.Assert(err, IsNil)
	c.Assert(string(data), Equals, string(remaining))
}

func (s *GopathSuite) TestIfTemplateIsNotFoundWhenVersionIsBad(c *C) {
	gpl := NewGoPathLoader()
	data, err := gpl.Load("nildev", "test-good", "v0.1.0")
	c.Assert(err, IsNil)
	c.Assert(err, IsNil)
	c.Assert(string(data), Equals, string([]byte{}))
}
