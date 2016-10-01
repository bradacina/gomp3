package main

func main() {
	c := getConfig()

	p := NewPlayer(c.Mp3Location)

	StartWeb(p, c.WebappPath, c.Port)
}
