package vault

type credential struct {
	name           string
	value          string
	path           string
	kmsId          string
	obfuscationKey string
	loaded         bool // is the local copy up-to-date with the current value
	saved          bool // is the current value saved to the backend
	client         *ssmClient
}

// If Value is empty, use existing value from backend
// Otherwise immediately save provided value and use it
func NewCredential(
	name string,
	value string,
	path string,
	kmsId string,
	obfuscationKey string,
	client *ssmClient,
) (*credential, error) {
	credential := credential{
		name:           name,
		value:          value,
		path:           path,
		kmsId:          kmsId,
		obfuscationKey: obfuscationKey,
		loaded:         false,
		saved:          false,
		client:         client,
	}
	if value != "" {
		credential.loaded = true
		err := credential.Save()
		if err == nil {
			credential.saved = true
		}
		return &credential, err
	}
	// Lazy backend fetching, don't pull the value before its needed
	return &credential, nil
}

func (c *credential) GetName() string {
	return c.name
}

func (c *credential) GetValue() (string, error) {
	if !c.loaded {
		err := c.Load()
		if err != nil {
			return "", err
		}
	}
	return c.value, nil
}

func (c *credential) SetValue(value string) {
	c.value = value
	c.loaded = true
	c.saved = false
}

func (c *credential) Save() error {
	if c.saved {
		return nil
	}
	obf_path, err := obfuscatePath(c.path, c.name, c.obfuscationKey)
	if err != nil {
		return err
	}
	err = c.client.PutParameter(obf_path, c.value, c.kmsId)
	if err == nil {
		c.saved = true
	}
	return err
}

func (c *credential) Load() error {
	if c.loaded {
		return nil
	}
	obf_path, err := obfuscatePath(c.path, c.name, c.obfuscationKey)
	if err != nil {
		return err
	}
	value, err := c.client.GetParameter(obf_path)
	if err == nil {
		c.value = value
		c.loaded = true
	}
	return err
}

func (c *credential) Delete() error {
	obfPath, err := obfuscatePath(c.path, c.name, c.obfuscationKey)
	if err != nil {
		return err
	}
	err = c.client.DeleteParameter(obfPath)
	c.name = ""
	c.value = ""
	c.loaded = false
	c.saved = false
	return err
}
