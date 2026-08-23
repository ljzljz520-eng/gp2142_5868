package store

var bucketNames = map[string][]byte{
	"records":     []byte("records"),
	"workflows":   []byte("workflows"),
	"audits":      []byte("audits"),
	"attachments": []byte("attachments"),
}

func bucket(name string) []byte {
	return bucketNames[name]
}
