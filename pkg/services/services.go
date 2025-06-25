package services

type Config struct {
	AuthService      string `json:"auth_service" yaml:"auth_service" toml:"auth_service" mapstructure:"auth_service"`                     // nolint
	MeService        string `json:"me_service" yaml:"me_service" toml:"me_service" mapstructure:"me_service"`                             // nolint
	BazarService     string `json:"bazar_service" yaml:"bazar_service" toml:"bazar_service" mapstructure:"bazar_service"`                 // nolint
	BankService      string `json:"bank_service" yaml:"bank_service" toml:"bank_service" mapstructure:"bank_service"`                     // nolint
	PortfolioService string `json:"portfolio_service" yaml:"portfolio_service" toml:"portfolio_service" mapstructure:"portfolio_service"` // nolint
	IBService        string `json:"ib_service" yaml:"ib_service" toml:"ib_service" mapstructure:"ib_service"`                             // nolint
	EmailPusher      string `json:"email_pusher" yaml:"email_pusher" toml:"email_pusher" mapstructure:"email_pusher"`
	BankInfo         string `json:"bank_info" yaml:"bank_info" toml:"bank_info" mapstructure:"bank_info"`
	AppSettings      string `json:"app_settings" yaml:"app_settings" toml:"app_settings" mapstructure:"app_settings"`
	Tyrion           string `json:"tyrion" yaml:"tyrion" toml:"tyrion" mapstructure:"tyrion"`
}
