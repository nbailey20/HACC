package vault

type state struct {
	Name   string
	Value  string
	Path   string
	loaded bool // is the local copy up-to-date with the current value
	saved  bool // is the current value saved to the backend
	client *ssmClient
}

// If Value is empty, use existing value from backend
// Otherwise immediately save provided value and use it
func NewState(name string, value string, path string) (*state, error) {
	state := state{
		Name:   name,
		Value:  value,
		Path:   path,
		loaded: false,
		saved:  false,
		client: NewSsmClient(),
	}
	if value != "" {
		state.loaded = true
		err := state.Save()
		if err == nil {
			state.saved = true
		}
		return &state, err
	}
	return &state, nil
}

func (s *state) GetName() string {
	return s.Name
}

// Lazy backend fetching, don't pull the value before its needed
func (s *state) GetValue() (string, error) {
	if !s.loaded {
		err := s.Load()
		if err != nil {
			return "", err
		}
	}
	return s.Value, nil
}

func (s *state) SetValue(value string) {
	s.Value = value
	s.loaded = true
	s.saved = false
}

func (s *state) Save() error {
	if s.saved {
		return nil
	}
	err := s.client.PutParameter(s.Path+s.Name, s.Value)
	if err == nil {
		s.saved = true
	}
	return err
}

func (s *state) Load() error {
	if s.loaded {
		return nil
	}
	value, err := s.client.GetParameter(s.Path + s.Name)
	if err == nil {
		s.Value = value
		s.loaded = true
	}
	return err
}

func (s *state) Delete() error {
	err := s.client.DeleteParameter(s.Path + s.Name)
	s.Name = ""
	s.Value = ""
	s.loaded = false
	s.saved = false
	return err
}
