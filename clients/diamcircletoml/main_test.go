package diamcircletoml

import "log"

// ExampleGetTOML gets the diamcircle.toml file for coins.asia
func ExampleClient_GetDiamcircleToml() {
	_, err := DefaultClient.GetDiamcircleToml("coins.asia")
	if err != nil {
		log.Fatal(err)
	}
}
