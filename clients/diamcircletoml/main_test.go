package diamcircletoml

import "log"

// ExampleGetTOML gets the diamcircle.toml file for coins.asia
func ExampleClient_GetdiamcircleToml() {
	_, err := DefaultClient.GetdiamcircleToml("coins.asia")
	if err != nil {
		log.Fatal(err)
	}
}
