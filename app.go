package main

func run() error {
	g := NewGame()
	g.Init()
	if err := g.Play(); err != nil {
		return err
	}
	return nil
}
