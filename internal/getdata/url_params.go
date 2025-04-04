package getdata

type BaseUrlParams struct {
	Host     string
	User     string
	Password string
}

type PmodeUrlParams struct {
	BaseUrlParams
	Machine string
	Point   string
	Pmode   string
}

func NewPmodeUrlParams(host, machine, point, pmode, user, password string) PmodeUrlParams {
	return PmodeUrlParams{
		BaseUrlParams: BaseUrlParams{
			Host:     host,
			User:     user,
			Password: password,
		},
		Machine: machine,
		Point:   point,
		Pmode:   pmode,
	}
}

type PmodeUrlTimeParams struct {
	PmodeUrlParams
	DateTime string
}

func NewPmodeUrlTimeParams(
	host, machine, point, pmode, time, user, password string,
) PmodeUrlTimeParams {
	return PmodeUrlTimeParams{
		PmodeUrlParams: NewPmodeUrlParams(host, machine, point, pmode, user, password),
		DateTime:       time,
	}
}
