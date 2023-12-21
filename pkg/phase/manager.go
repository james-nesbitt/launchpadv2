package phase

import "context"

type Manager struct {
	phases []Phase
}

// Run execute all of the loaded phase.
func (m *Manager) Run(ctx context.Context) error {
	cctx, cancel := context.WithCancel(ctx)

	for _, p := range m.phases {
		if err := p.Run(cctx); err != nil {
			cancel()
			return err
		}
	}

	cancel()

	return nil
}
