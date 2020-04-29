package caching

import (
	"reflect"
	"testing"
	"time"

	"github.com/hayrullahcansu/fastmta-core/dns"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
)

//New Cache instance al, içine değer set et geliyor mu test et, birde timerkoyup test et.
func TestDomainCache(t *testing.T) {
	instance := newInstanceDomain(time.Second * 2)
	if instance == nil {
		t.Fatalf("instance cannot be nil")
	}
	_host := "gmail.com"
	expected, _ := dns.NewDomain(_host)
	instance.AddOrUpdate(_host, expected)
	o := instance.Get(_host)
	if o == nil {
		t.Fatalf("The cache removed the data earlier")
	}

	if !assert.EqualValues(t, expected, o) {
		t.Errorf("the cache returns different value")
	}
	if !reflect.DeepEqual(expected, o) {
		t.Errorf("the cache returns different value")
	}
	if diff := deep.Equal(expected, o); diff != nil {
		t.Errorf("the cache returns different value")
	}
	time.Sleep(time.Second * 3)
	o = instance.Get(_host)
	if o != nil {
		t.Errorf("the cache didn't clean in expected time")
	}
}
