package core

type Config struct {
	Label      string
	DamageMode bool
	Targets    []EnemyProfile
	Characters struct {
		Initial string
		Profile []CharacterProfile
	}
	Rotation []Action

	Hurt      HurtEvent
	FixedRand bool //if this is true then use the same seed
}

type RunOpt struct {
	LogDetails bool `json:"log_details"`
	Iteration  int  `json:"iter"`
	Workers    int  `json:"workers"`
	Duration   int  `json:"seconds"`
	Debug      bool `json:"debug"`
	DebugPaths []string
}

type CharacterProfile struct {
	Base    CharacterBase
	Weapon  WeaponProfile
	Talents TalentProfile
	Stats   []float64
	Sets    map[string]int
}

type CharacterBase struct {
	Name    string
	Element EleType
	Level   int
	HP      float64
	Atk     float64
	Def     float64
	Cons    int
	StartHP float64
}

type WeaponProfile struct {
	Name   string
	Class  WeaponClass
	Refine int
	Atk    float64
	Param  map[string]int
}

type TalentProfile struct {
	Attack int
	Skill  int
	Burst  int
}

type EnemyProfile struct {
	Level  int
	HP     float64
	Resist map[EleType]float64
}

type HurtEvent struct {
	WillHurt bool
	Once     bool //how often
	Start    int  //
	End      int
	Min      float64
	Max      float64
	Ele      EleType
}

func (e *EnemyProfile) Clone() EnemyProfile {
	r := EnemyProfile{
		Level:  e.Level,
		Resist: make(map[EleType]float64),
	}
	for k, v := range e.Resist {
		r.Resist[k] = v
	}
	return r
}

func (c *CharacterProfile) Clone() CharacterProfile {
	r := *c
	r.Weapon.Param = make(map[string]int)
	for k, v := range c.Weapon.Param {
		r.Weapon.Param[k] = v
	}
	r.Stats = make([]float64, len(c.Stats))
	copy(r.Stats, c.Stats)
	r.Sets = make(map[string]int)
	for k, v := range c.Sets {
		r.Sets[k] = v
	}

	return r
}

func (c *Config) Clone() Config {
	r := *c

	r.Targets = make([]EnemyProfile, len(c.Targets))

	for i, v := range c.Targets {
		r.Targets[i] = v.Clone()
	}

	return r
}
