package auth

import "time"

type PublishACL struct {
	Pattern       string `msgpack:"pattern"`
	MaxQos        *uint8 `msgpack:"max_qos,omitempty"`
	AllowedRetain *bool  `msgpack:"allowed_retain,omitempty"`
}

type SubACL struct {
	Pattern string `msgpack:"pattern"`
	MaxQos  *uint8 `msgpack:"max_qos,omitempty"`
}

type User struct {
	Username   string       `msgpack:"username"`
	Password   string       `msgpack:"password"`
	ClientID   *string      `msgpack:"client_id,omitempty"`
	Mountpoint string       `msgpack:"mountpoint"`
	PublishACL []PublishACL `msgpack:"publish_acl,omitempty"`
	SubACL     []SubACL     `msgpack:"sub_acl,omitempty"`
	CreatedAt  time.Time    `msgpack:"created_at"`
	UpdatedAt  time.Time    `msgpack:"updated_at"`
}
