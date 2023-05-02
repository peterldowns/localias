package hostctl

var _ Controller = MultiController{}

type MultiController []Controller

func NewMultiController(controllers ...Controller) MultiController {
	return controllers
}

func (mc MultiController) Set(ip string, alias string) error {
	for _, c := range mc {
		if err := c.Set(ip, alias); err != nil {
			return err
		}
	}
	return nil
}

func (mc MultiController) SetLocal(alias string) error {
	for _, c := range mc {
		if err := c.SetLocal(alias); err != nil {
			return err
		}
	}
	return nil
}

func (mc MultiController) Remove(alias string) error {
	for _, c := range mc {
		if err := c.Remove(alias); err != nil {
			return err
		}
	}
	return nil
}

func (mc MultiController) Clear() error {
	for _, c := range mc {
		if err := c.Clear(); err != nil {
			return err
		}
	}
	return nil
}

func (mc MultiController) Apply() (bool, error) {
	anyChanges := false
	for _, c := range mc {
		changed, err := c.Apply()
		if err != nil {
			return false, err
		}
		anyChanges = anyChanges || changed
	}
	return anyChanges, nil
}
