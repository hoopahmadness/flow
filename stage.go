package flowchart

type Stage struct {
	Name        string   `json:"name"`
	Transitions []string `json:"transitions"`
}

func NewStage(name string) Stage {
	return Stage{
		Name:        name,
		Transitions: []string{},
	}
}

func (s *Stage) addTransition(t string) {
	s.Transitions = append(s.Transitions, t)
}
