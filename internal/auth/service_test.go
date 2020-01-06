package auth

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestPrivatePublishAuthorize(t *testing.T) {
	srv, err := New(1, nil, logrus.New())
	if err != nil {
		t.Fatal(err)
	}

	var qos1 uint8 = 1
	var qos2 uint8 = 2

	var allowedRetain bool = true
	var notAllowedRetain bool = false

	testCases := []struct {
		msg    string
		u      string
		i      string
		t      string
		q      uint8
		r      bool
		p      []PublishACL
		result bool
	}{
		{msg: "should allow all publishes", u: "username", i: "my-client-id", p: nil, t: "foo/bar/baz", q: 1, r: true, result: true},
		{msg: "topic should match with pattern", u: "username", i: "my-client-id", p: []PublishACL{{Pattern: "foo/+/baz"}}, t: "foo/bar/baz", q: 1, r: true, result: true},
		{msg: "topic shouldn't match with pattern", u: "username", i: "my-client-id", p: []PublishACL{{Pattern: "foo/BAR/baz"}}, t: "foo/bar/baz", q: 1, r: true, result: false},
		{msg: "publish qos is equal or lesser than max q", u: "username", i: "my-client-id", p: []PublishACL{{Pattern: "foo/+/baz", MaxQos: &qos2}}, t: "foo/bar/baz", q: 2, r: true, result: true},
		{msg: "publish qos is greater than max q", u: "username", i: "my-client-id", p: []PublishACL{{Pattern: "foo/+/baz", MaxQos: &qos1}}, t: "foo/bar/baz", q: 2, r: true, result: false},
		{msg: "should allow publishing retained messages", u: "username", i: "my-client-id", p: []PublishACL{{Pattern: "foo/+/baz", AllowedRetain: &allowedRetain}}, t: "foo/bar/baz", q: 2, r: true, result: true},
		{msg: "shouldn't allow publishing retained messages", u: "username", i: "my-client-id", p: []PublishACL{{Pattern: "foo/+/baz", AllowedRetain: &notAllowedRetain}}, t: "foo/bar/baz", q: 2, r: true, result: false},
		{msg: "should allow topics with clientID", u: "username", i: "my-client-id", p: []PublishACL{{Pattern: "foo/{{client_id}}/baz"}}, t: "foo/my-client-id/baz", q: 2, r: true, result: true},
		{msg: "shouldn't allow topics with invalid clientID in it", u: "username", i: "my-client-id", p: []PublishACL{{Pattern: "foo/{{client_id}}/baz"}}, t: "foo/invlaid-client-id/baz", q: 2, r: true, result: false},
		{msg: "should allow topics with username", u: "username", i: "my-client-id", p: []PublishACL{{Pattern: "foo/{{username}}/baz"}}, t: "foo/username/baz", q: 2, r: true, result: true},
		{msg: "shouldn't allow topics with invalid username in it", u: "username", i: "my-client-id", p: []PublishACL{{Pattern: "foo/{{username}}/baz"}}, t: "foo/invalid-username/baz", q: 2, r: true, result: false},
		{msg: "should return false on invalid templates", u: "username", i: "my-client-id", p: []PublishACL{{Pattern: "foo/{{username/baz"}}, t: "foo/invalid-username/baz", q: 2, r: true, result: false},
	}

	for _, tt := range testCases {
		isOk := srv.authorizePublish(tt.i, tt.u, tt.t, tt.q, tt.r, tt.p)
		assert.Equalf(t, tt.result, isOk, tt.msg)
	}
}

func TestPrivateSubscribeAuthorize(t *testing.T) {
	srv, err := New(1, nil, logrus.New())
	if err != nil {
		t.Fatal(err)
	}

	var qos1 uint8 = 1

	testCases := []struct {
		msg    string
		u      string
		i      string
		t      string
		q      uint8
		r      bool
		subs   []Subscription
		acl    []SubACL
		result []Subscription
	}{
		{msg: "should allow all subscriptions", u: "username", i: "my-client-id", subs: []Subscription{{Topic: "foo/bar/baz", Qos: 1}}, acl: nil, result: []Subscription{{Topic: "foo/bar/baz", Qos: 1}}},
		{msg: "should allow matching patterns", u: "username", i: "my-client-id", subs: []Subscription{{Topic: "foo/bar/baz", Qos: 1}}, acl: []SubACL{{Pattern: "foo/+/baz"}}, result: []Subscription{{Topic: "foo/bar/baz", Qos: 1}}},
		{msg: "shouldn't allow invalid patterns", u: "username", i: "my-client-id", subs: []Subscription{{Topic: "foo/bar/baz", Qos: 1}}, acl: []SubACL{{Pattern: "foo/BAZ/baz"}}, result: []Subscription{{Topic: "foo/bar/baz", Qos: InvalidSubscription}}},
		{msg: "should allow subscriptions with valid qos", u: "username", i: "my-client-id", subs: []Subscription{{Topic: "foo/bar/baz", Qos: 1}}, acl: []SubACL{{Pattern: "foo/bar/baz", MaxQos: &qos1}}, result: []Subscription{{Topic: "foo/bar/baz", Qos: 1}}},
		{msg: "shouldn't allow subscriptions with invalid qos", u: "username", i: "my-client-id", subs: []Subscription{{Topic: "foo/bar/baz", Qos: 2}}, acl: []SubACL{{Pattern: "foo/bar/baz", MaxQos: &qos1}}, result: []Subscription{{Topic: "foo/bar/baz", Qos: InvalidSubscription}}},
		{msg: "shouldn't allow topics with valid clientID", u: "username", i: "my-client-id", subs: []Subscription{{Topic: "foo/my-client-id/baz", Qos: 1}}, acl: []SubACL{{Pattern: "foo/{{client_id}}/baz"}}, result: []Subscription{{Topic: "foo/my-client-id/baz", Qos: 1}}},
		{msg: "shouldn't allow invalid parameterized patterns", u: "username", i: "my-client-id", subs: []Subscription{{Topic: "foo/my-client-id/baz", Qos: 2}}, acl: []SubACL{{Pattern: "foo/{{client_id/baz"}}, result: []Subscription{{Topic: "foo/my-client-id/baz", Qos: InvalidSubscription}}},
		{msg: "should handle multiple subscriptions", u: "username", i: "my-client-id", subs: []Subscription{{Topic: "foo/bar/baz", Qos: 1}, {Topic: "baz/bar/foo", Qos: 1}}, acl: []SubACL{{Pattern: "foo/bar/baz", MaxQos: &qos1}}, result: []Subscription{{Topic: "foo/bar/baz", Qos: 1}, {Topic: "baz/bar/foo", Qos: InvalidSubscription}}},
	}

	for _, tt := range testCases {
		isOk := srv.authorizeSubscribe(tt.i, tt.u, tt.subs, tt.acl)
		assert.Equalf(t, tt.result, isOk, tt.msg)
	}
}
