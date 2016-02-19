package setup

// NewConfigLoader returns instance of config Loader based on extension
func NewConfigLoader(pathToConfig string) ConfigReader {

	// todo: implement switch
	return NewJSONLoader()
}
