package main

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	/* Note:
	 * Name string `json:"-"`
	 * placing a - omits it from the marshaled json
	 */
}

var users = map[string]User{
	"2d711642b726b04401627ca9fbac32": {
		ID:   "2d711642b726b04401627ca9fbac32",
		Name: "Alice",
	},
	"594e519ae499312b29433b7d": {
		ID:   "594e519ae499312b29433b7d",
		Name: "Bob",
	},
}
