package vault

type credential struct {
	name   string
	value  string
	path   string
	kms_id string
	loaded bool // is the local copy up-to-date with the current value
	saved  bool // is the current value saved to the backend
	client *ssmClient
}

// If Value is empty, use existing value from backend
// Otherwise immediately save provided value and use it
func NewCredential(
	name string,
	value string,
	path string,
	kms_id string,
	client *ssmClient,
) (*credential, error) {
	credential := credential{
		name:   name,
		value:  value,
		path:   path,
		kms_id: kms_id,
		loaded: false,
		saved:  false,
		client: client,
	}
	if value != "" {
		credential.loaded = true
		err := credential.Save()
		if err == nil {
			credential.saved = true
		}
		return &credential, err
	}
	return &credential, nil
}

func (c *credential) GetName() string {
	return c.name
}

// Lazy backend fetching, don't pull the value before its needed
func (s *credential) GetValue() (string, error) {
	if !s.loaded {
		err := s.Load()
		if err != nil {
			return "", err
		}
	}
	return s.value, nil
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
	err := c.client.PutParameter(c.path+c.name, c.value, c.kms_id)
	if err == nil {
		c.saved = true
	}
	return err
}

func (c *credential) Load() error {
	if c.loaded {
		return nil
	}
	value, err := c.client.GetParameter(c.path + c.name)
	if err == nil {
		c.value = value
		c.loaded = true
	}
	return err
}

func (c *credential) Delete() error {
	err := c.client.DeleteParameter(c.path + c.name)
	c.name = ""
	c.value = ""
	c.loaded = false
	c.saved = false
	return err
}
